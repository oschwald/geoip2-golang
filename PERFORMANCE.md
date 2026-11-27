# Real-MMDB performance investigation

## Current implementation

[maxminddb-golang PR #267](https://github.com/oschwald/maxminddb-golang/pull/267)
now optimizes `decodeKeyAt`, not `Cursor.ReadMapKey`. The shared key decoder
handles one-byte size extensions (29–284-byte strings) without entering the
general decoder. Inline keys and all four pointer widths use this path.
Generated, reflection, and iterator decoders can benefit.

Review identified regressions in the original cursor-only implementation:
inline keys and wider pointers paid for an extra fallback call. Moving the
optimization restores the original `ReadMapKey` call chain and its `validate()`
check. The change also removes the extra copy of pointer-following logic.
Compact keys keep their existing shared-decoder path.

A separate commit limits borrowed key slices to their payload capacity. An
append now allocates instead of overwriting subsequent database bytes. The
returned bytes still alias the input and must not be modified directly.

Neither change adds a cache, retained reader state, or decoding allocations.
Cache capacity and policy remain unchanged. This library's `go.mod`, public
methods, and generated models are unchanged.

## Method

The review follow-up compares three dependency revisions:

- Baseline: `origin/main` at `3359e15`, the v2.9.0 release.
- Original PR: `ace4f3f`, the cursor-only optimization.
- Revised PR: `18d984a`, including shared decoding (`19c3273`) and key-capacity
  hardening.

The benchmark harness is geoip2-golang commit `09fe119`. Environment: Go 1.27.3,
Linux/amd64, Intel Core Ultra 7 265K. Databases came from `/var/lib/GeoIP`:
GeoIP2 City, Enterprise, Country, ISP, GeoLite2 City, and GeoLite2 ASN.

The real-database benchmarks use:

- 16,384 deterministic successful addresses per database and family.
- Uniform IPv4 candidates, filtered to addresses with records.
- IPv6 candidates sampled from `2001::/16` and the `2400`, `2600`, `2800`,
  `2a00`, and `2c00` /12 allocations, filtered to addresses with records.
- Two warm-up passes before timing.
- Fresh generated results for full public lookups and decode-only measurements.
- Separate fresh reflection decoding with no search-tree traversal. City and
  Enterprise use deep mirror types so nested generated methods cannot run.

These are address-space samples, not production traffic traces. They do not
measure misses, cold starts, or actual traffic's distribution of fields. IPv6
records in the sample are often less detailed. The corpus repeats during timing.

The main comparison used eight independent 250 ms samples per case. Execution
order rotated across all three versions, pinned to CPU 0 with `GOMAXPROCS=1`.
The key microbenchmark used eight 150 ms samples per case. Tables show medians,
analyzed with `benchstat`. Negative changes mean less time per operation.
“No clear change” means p >= 0.05, not proof of equivalent performance.

## Key microbenchmarks

Unlike the original benchmark, `BenchmarkMapKeyEncodings` covers inline keys,
all pointer widths, and both `Cursor.ReadMapKey` and the shared decoder.

Selected cursor results, in ns/key:

| Encoding                 | Key bytes | Baseline | Original PR | Revised PR |
| ------------------------ | --------: | -------: | ----------: | ---------: |
| Inline                   |         2 |    1.882 |       2.438 |      1.872 |
| Inline                   |        30 |    3.761 |       4.482 |      2.062 |
| One-byte pointer payload |         2 |    2.237 |       1.889 |      2.083 |
| One-byte pointer payload |        30 |    5.764 |       1.959 |      2.434 |
| Two-byte pointer payload |         2 |    2.523 |       3.087 |      2.493 |
| Two-byte pointer payload |        30 |    6.401 |       6.962 |      3.094 |
| One-byte pointer payload |       285 |    5.948 |       6.731 |      6.205 |

The revised code removes the roughly 29% inline-short-key and 22%
two-byte-pointer regressions. Thirty-byte keys improve 45–58% through the
cursor API and 52–60% through the shared decoder across the tested encodings.

There is still a trade-off. Some 285-byte-key cases are 1–5% slower than the
baseline, about 0.05–0.3 ns/key. They incur an extra size-extension check before
falling back. GeoIP2 model field names do not reach this size. All key
microbenchmarks allocate zero bytes.

## Real-database results

Full public lookups, revised PR versus baseline:

| Database     | Family | Baseline ns/op | Revised ns/op |          Change |
| ------------ | ------ | -------------: | ------------: | --------------: |
| City         | IPv4   |          1,114 |         1,133 | No clear change |
| City         | IPv6   |          737.2 |         742.0 | No clear change |
| GeoLite City | IPv4   |          1,065 |         1,044 | No clear change |
| GeoLite City | IPv6   |          735.3 |         744.2 | No clear change |
| Enterprise   | IPv4   |          1,476 |         1,526 |         +3.39%* |
| Enterprise   | IPv6   |          876.6 |         882.7 | No clear change |
| Country      | IPv4   |          565.1 |         569.9 | No clear change |
| Country      | IPv6   |          581.5 |         578.9 | No clear change |
| ASN          | IPv4   |          147.0 |         141.6 |          -3.64% |
| ASN          | IPv6   |          162.7 |         157.1 |          -3.44% |
| ISP          | IPv4   |          205.6 |         199.2 |          -3.09% |
| ISP          | IPv6   |          210.5 |         203.5 |          -3.30% |

*The Enterprise IPv4 regression did not repeat in a longer follow-up: ten
alternating one-second samples measured 1,526 versus 1,543 ns/op, p=0.101.
That follow-up also found no significant City or GeoLite City IPv4 difference.
Their revised lookup medians remained slightly higher, so this is not a claim
that regressions are impossible.

Generated decode-only results improve 7.6–8.4% for ASN and 5.2–5.4% for ISP.
Fresh reflection decoding improves 3.1–4.2% for ASN and 2.9–3.1% for ISP.
Most larger-record decode comparisons are within noise. Small 0.5–0.8%
regressions appeared in some IPv6 City, GeoLite City, and Country decode cases.
These results do not establish a general City/Enterprise speed-up.

Bytes and allocations per operation had no significant change.

The original PR retains a larger benefit on pointer-heavy generated lookups,
including City and GeoLite City. The revised implementation gives up most of
that benefit to avoid penalizing inline-key workloads and to share the
extended-size improvement with reflection and iterator callers.

The parallel benchmark remains available. Earlier parallel and single-address
results measured `ace4f3f`, not the revised implementation, and should not be
used as revised-PR performance claims.

## Review coverage and validation

- Restored the original cursor validation and removed its extra fallback call.
- Shared extended-size handling across inline keys and all pointer widths.
- Tested `extendedKeyBytes` directly: supported sizes must return a non-nil key,
  proving they are accepted without the general decoder.
- Compared cursor results and cache offsets with `decodeKeyAtSlow`.
- Tested capacity limits and append behavior for compact and extended keys,
  all pointer widths, and both optimized and general decoding.
- Added benchmarks for the inline and wider-pointer cases missing originally.
- Added unreleased changelog entries in maxminddb-golang.
- Passed both libraries' tests and race tests, plus decoder and `mmdbdata` tests
  on `GOARCH=386`.
- Passed fresh/reused generated-versus-reflection parity on 98,304 real City,
  GeoLite City, and Enterprise addresses, including a race-enabled run.
- Passed 60 seconds of differential fuzzing: 7,664,178 executions, no failures.
- Passed `golangci-lint` and formatting checks.

The review found no correctness bug in the original fast path. The follow-up
changes address performance coverage, shared implementation, and defensive
slice capacity rather than changing valid decoded values.

## Reproduce

Run from this checkout. Benchmarks use `GEOIP2_REAL_DB_DIR`, defaulting to
`/var/lib/GeoIP`. Missing default files are skipped. An explicit directory makes
missing files fail. The real-database parity test requires the environment
variable and is skipped by ordinary CI.

```sh
GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -run '^$' \
  -bench '^BenchmarkRealDatabase$/.*/.*/(Lookup|Decode)$' \
  -benchmem -benchtime=250ms -count=8 -cpu=1

GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -run '^$' \
  -bench '^BenchmarkRealDatabaseReflection$' \
  -benchmem -benchtime=250ms -count=8 -cpu=1

GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -run TestRealDatabaseParity
```

Check out the candidate dependency in `../maxminddb-golang-geoip2-perf` and use
an alternate module file without changing `go.mod`:

```sh
benchdir=$(mktemp -d)
cp go.mod "$benchdir/local.mod"
cp go.sum "$benchdir/local.sum"
go mod edit -modfile="$benchdir/local.mod" \
  -replace="github.com/oschwald/maxminddb-golang/v2=$(realpath ../maxminddb-golang-geoip2-perf)"
GEOIP2_REAL_DB_DIR=/var/lib/GeoIP go test -modfile="$benchdir/local.mod" \
  -run '^$' -bench '^BenchmarkRealDatabase$/.*/.*/(Lookup|Decode)$' \
  -benchmem -benchtime=250ms -count=8 -cpu=1
```

For close comparisons, compile each version with `go test -c`, pin the binaries
with `taskset`, and rotate execution order. Save each version's output and
compare it with `benchstat`. Run `BenchmarkMapKeyEncodings` in the dependency's
`internal/decoder` package for the key-encoding microbenchmarks.
