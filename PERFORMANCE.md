# Real-MMDB performance investigation

## Recommended change

Add a fast path to `Cursor.ReadMapKey` in `maxminddb-golang/v2`. Real databases
reuse map keys near the start of the data section. Read pointers with a one-byte
payload directly, without the `decodeKey` / `decodeKeyAt` call chain. Handle
both compact string lengths and the one-byte size extension used by keys such
as `autonomous_system_organization`.

The implementation is in
[maxminddb-golang PR #267](https://github.com/oschwald/maxminddb-golang/pull/267):

- `internal/decoder/cursor.go`: fast path and unchanged general-decoder fallback.
- `internal/decoder/cursor_map_key_test.go`: encoding, malformed-input,
  differential fuzz, and microbenchmark coverage.

This change needs no generated-code or public API changes. It adds no cache,
retained reader state, or allocations. The string-cache capacity and policy
remain unchanged. All other prototypes were removed. This checkout's `go.mod`
and generated models are unchanged.

## Method

Baseline: geoip2-golang commit `876fc55`, using maxminddb-golang v2.9.0.
Candidate: maxminddb-golang commit `ace4f3f`, based on `3359e15`. The baseline
`internal` directory matched the v2.9.0 module contents.

Environment: Go 1.27.3, Linux/amd64, Intel Core Ultra 7 265K. Databases came from
`/var/lib/GeoIP`: GeoIP2 City, Enterprise, Country, ISP, and GeoLite2 ASN.

`real_database_benchmark_test.go` supplies:

- 16,384 deterministic successful lookups per database and address family.
- Uniform IPv4 candidates, filtered to addresses with records.
- IPv6 candidates sampled from `2001::/16` and the `2400`, `2600`, `2800`,
  `2a00`, and `2c00` /12 allocations, filtered to addresses with records.
- Two warm-up passes before timing.
- Full public lookups with fresh results, decode-only with fresh results,
  tree traversal, and concurrent public lookups.

These are address-space samples, not a production traffic trace. They do not
measure misses, cold starts, or the distribution of fields in actual traffic.
IPv6 geolocation records in this sample are often less detailed than IPv4
records. The corpus repeats during timing.

Final comparisons used ten independent 300 ms samples per case, alternating
baseline/candidate order. Serial runs used CPU 0 and `GOMAXPROCS=1`. Parallel
runs used CPUs 0–7 and `GOMAXPROCS=8`. Results below are medians, analyzed with
`benchstat`. Earlier non-interleaved runs overstated some gains.

## Full-lookup results

Negative changes mean less time per operation. “No clear change” means the
comparison did not reach p < 0.05.

| Database   | Family | Baseline ns/op | Candidate ns/op |          Change |
| ---------- | ------ | -------------: | --------------: | --------------: |
| City       | IPv4   |          1,157 |           1,130 |          -2.29% |
| City       | IPv6   |          743.6 |           735.5 | No clear change |
| Enterprise | IPv4   |          1,555 |           1,548 | No clear change |
| Enterprise | IPv6   |          904.9 |           867.5 |          -4.14% |
| Country    | IPv4   |          569.0 |           550.8 |          -3.20% |
| Country    | IPv6   |          580.3 |           564.5 |          -2.71% |
| ASN        | IPv4   |          148.0 |           141.3 |          -4.56% |
| ASN        | IPv6   |          164.0 |           156.3 |          -4.70% |
| ISP        | IPv4   |          208.4 |           197.9 |          -5.01% |
| ISP        | IPv6   |          213.4 |           203.3 |          -4.73% |

Decode-only improvements were 1.8–3.2% for City, Enterprise, and Country,
10.6–11.1% for ASN, and 8.4–8.6% for ISP. Bytes and allocations per operation
had no significant change. Tree traversal was effectively unchanged.

The existing single-address commercial City benchmark improved from 802.9 to
754.4 ns/op, or 6.04%. This larger gain is not representative of varied traffic.

At eight goroutines, aggregate time per operation improved 3.0–4.6% for ASN,
2.0–2.5% for ISP, and 1.6% for Enterprise IPv4. Other parallel comparisons were
within measurement noise. Parallel ns/op measures throughput, not request
latency.

The key microbenchmark shows why the extended-size case matters: a 29-byte
key takes about 2.3 ns through the fast path versus 6.0 ns through the general
path. The fallback has a cost: a 285-byte key takes about 7.1 ns versus 6.2 ns.
Such long keys are outside this fast path and are not GeoIP2 model field names.

## Other experiments

- A pointer fast path in `MapReader` did not produce a consistent overall gain.
- Probing the string cache's home slot before scanning did not help overall.
- Consolidating compact string decoding did not improve overall performance.
- Always scanning string-cache buckets produced only small, inconsistent gains.
- Increasing the cache size was excluded because it increases retained memory.

The mixed City/Enterprise CPU profile put about 20% of samples under string
interning and about 7% under map-key reading. Decoding, allocation, and garbage
collection still dominate large result types. This patch removes some key
handling overhead without changing those memory costs.

## Validation

Passed with the candidate dependency:

- Both libraries' normal and race-enabled test suites.
- Decoder and `mmdbdata` tests on `GOARCH=386`.
- Generated-versus-reflection parity on 65,536 real City/Enterprise addresses,
  checking both fresh and reused destinations.
- All pointer widths, high pointer bits, compact and extended string sizes,
  borrowed key bytes, and the cursor's position after the pointer.
- Malformed keys, truncated headers and payloads, pointer chains, invalid
  offsets, and zero cursors.
- 45 seconds of differential fuzzing: 5,741,212 executions, no failures.
- `golangci-lint run ./...` in both checkouts.

## Reproduce

Run from this checkout. An unset `GEOIP2_REAL_DB_DIR` defaults to
`/var/lib/GeoIP` for benchmarks. Missing default files are skipped. An explicit
path makes missing files fail. The real-database parity test is opt-in.

```sh
GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -run '^$' \
  -bench '^BenchmarkRealDatabase/.*/.*/(Lookup|Decode|Traverse)$' \
  -benchmem -benchtime=300ms -count=10 -cpu=1

GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -run TestRealDatabaseParity
```

For the candidate run, check out maxminddb-golang commit `ace4f3f` in
`../maxminddb-golang-geoip2-perf`. Use an alternate module file to test that
checkout without changing `go.mod`:

```sh
benchdir=$(mktemp -d)
cp go.mod "$benchdir/local.mod"
cp go.sum "$benchdir/local.sum"
go mod edit -modfile="$benchdir/local.mod" \
  -replace="github.com/oschwald/maxminddb-golang/v2=$(realpath ../maxminddb-golang-geoip2-perf)"
GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -modfile="$benchdir/local.mod" \
  -run '^$' -bench '^BenchmarkRealDatabase/.*/.*/(Lookup|Decode|Traverse)$' \
  -benchmem -benchtime=300ms -count=10 -cpu=1
```

For close comparisons, compile each version with `go test -c`, pin the test
binaries with `taskset`, and alternate their execution order. Save both versions'
output and compare the results with `benchstat`.
