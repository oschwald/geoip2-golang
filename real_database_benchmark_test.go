package geoip2

import (
	"encoding/binary"
	"errors"
	"math/rand/v2"
	"net/netip"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2"
)

// BenchmarkRealDatabase uses GEOIP2_REAL_DB_DIR, or /var/lib/GeoIP when unset.
// Traffic is a fixed-seed sample of addresses with records, not a traffic trace.
// IPv4 addresses are uniform across the address space. IPv6 addresses are
// sampled from 2001::/16 and the 2400, 2600, 2800, 2a00, and 2c00 /12
// allocations. Each family has 16,384 addresses.
func BenchmarkRealDatabase(b *testing.B) {
	benchmarkRealDatabase(b, "City", "GeoIP2-City.mmdb", (*Reader).City)
	benchmarkRealDatabase(b, "GeoLiteCity", "GeoLite2-City.mmdb", (*Reader).City)
	benchmarkRealDatabase(b, "Enterprise", "GeoIP2-Enterprise.mmdb", (*Reader).Enterprise)
	benchmarkRealDatabase(b, "Country", "GeoIP2-Country.mmdb", (*Reader).Country)
	benchmarkRealDatabase(b, "ASN", "GeoLite2-ASN.mmdb", (*Reader).ASN)
	benchmarkRealDatabase(b, "ISP", "GeoIP2-ISP.mmdb", (*Reader).ISP)
}

func openRealDatabase(tb testing.TB, name string) *Reader {
	tb.Helper()
	dir := os.Getenv("GEOIP2_REAL_DB_DIR")
	explicit := dir != ""
	if !explicit {
		dir = "/var/lib/GeoIP"
	}
	reader, err := Open(filepath.Join(dir, name))
	if !explicit && errors.Is(err, os.ErrNotExist) {
		tb.Skip("real database is absent; set GEOIP2_REAL_DB_DIR")
	}
	require.NoError(tb, err)
	tb.Cleanup(func() { require.NoError(tb, reader.Close()) })
	return reader
}

func realDatabaseAddresses(tb testing.TB, reader *Reader, ipv6 bool) []netip.Addr {
	tb.Helper()
	const count = 16_384
	addresses := make([]netip.Addr, 0, count)
	//nolint:gosec // A fixed seed makes benchmark traffic reproducible.
	rng := rand.New(rand.NewPCG(1, 2))
	for attempts := 0; len(addresses) < count && attempts < count*1024; attempts++ {
		var address netip.Addr
		if ipv6 {
			var bytes [16]byte
			binary.BigEndian.PutUint64(bytes[:8], rng.Uint64())
			binary.BigEndian.PutUint64(bytes[8:], rng.Uint64())
			prefix := [...]byte{0x20, 0x24, 0x26, 0x28, 0x2a, 0x2c}[rng.IntN(6)]
			bytes[0] = prefix
			bytes[1] &= 0x0f
			if prefix == 0x20 {
				bytes[1] = 1
			}
			address = netip.AddrFrom16(bytes)
		} else {
			var bytes [4]byte
			binary.BigEndian.PutUint32(bytes[:], rng.Uint32())
			address = netip.AddrFrom4(bytes)
		}
		result := reader.mmdbReader.Lookup(address)
		require.NoError(tb, result.Err())
		if result.Found() {
			addresses = append(addresses, address)
		}
	}
	if len(addresses) != count {
		tb.Fatalf("database supplied only %d of %d matching addresses", len(addresses), count)
	}
	return addresses
}

func benchmarkRealDatabase[T any](
	b *testing.B,
	name, file string,
	lookup func(*Reader, netip.Addr) (*T, error),
) {
	b.Helper()
	b.Run(name, func(b *testing.B) {
		for _, family := range []string{"IPv4", "IPv6"} {
			b.Run(family, func(b *testing.B) {
				reader := openRealDatabase(b, file)
				addresses := realDatabaseAddresses(b, reader, family == "IPv6")
				results := make([]maxminddb.Result, len(addresses))
				for i, address := range addresses {
					results[i] = reader.mmdbReader.Lookup(address)
				}
				// Warm the string cache and mapped pages before timing.
				for range 2 {
					for _, address := range addresses {
						_, err := lookup(reader, address)
						require.NoError(b, err)
					}
				}
				b.Run("Lookup", func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; b.Loop(); i++ {
						if _, err := lookup(reader, addresses[i%len(addresses)]); err != nil {
							b.Fatal(err)
						}
					}
				})
				b.Run("Parallel", func(b *testing.B) {
					var starts atomic.Uint64
					b.ReportAllocs()
					b.ResetTimer()
					b.RunParallel(func(pb *testing.PB) {
						i := int(starts.Add(1)*104729) % len(addresses)
						for pb.Next() {
							if _, err := lookup(reader, addresses[i]); err != nil {
								b.Error(err)
								return
							}
							i++
							if i == len(addresses) {
								i = 0
							}
						}
					})
				})
				b.Run("Decode", func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; b.Loop(); i++ {
						if err := results[i%len(results)].Decode(new(T)); err != nil {
							b.Fatal(err)
						}
					}
				})
				b.Run("Traverse", func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; b.Loop(); i++ {
						result := reader.mmdbReader.Lookup(addresses[i%len(addresses)])
						if err := result.Err(); err != nil {
							b.Fatal(err)
						}
					}
				})
			})
		}
	})
}

// TestRealDatabaseParity is opt-in because it needs commercial database files.
func TestRealDatabaseParity(t *testing.T) {
	if os.Getenv("GEOIP2_REAL_DB_DIR") == "" {
		t.Skip("set GEOIP2_REAL_DB_DIR to test real databases")
	}
	testRealDatabaseParity(t, "GeoIP2-City.mmdb", cityFromReflection)
	testRealDatabaseParity(t, "GeoLite2-City.mmdb", cityFromReflection)
	testRealDatabaseParity(t, "GeoIP2-Enterprise.mmdb", enterpriseFromReflection)
}

func testRealDatabaseParity[T, R any](t *testing.T, file string, fromReflection func(R) T) {
	t.Helper()
	t.Run(file, func(t *testing.T) {
		reader := openRealDatabase(t, file)
		for _, ipv6 := range []bool{false, true} {
			var generated T
			var reflection R
			for _, address := range realDatabaseAddresses(t, reader, ipv6) {
				result := reader.mmdbReader.Lookup(address)
				// Reuse destinations to check merge behavior as well as values.
				require.NoError(t, result.Decode(&generated))
				require.NoError(t, result.Decode(&reflection))
				require.Equal(t, fromReflection(reflection), generated,
					"reused, address %s", address)

				var freshGenerated T
				var freshReflection R
				require.NoError(t, result.Decode(&freshGenerated))
				require.NoError(t, result.Decode(&freshReflection))
				require.Equal(t, fromReflection(freshReflection), freshGenerated,
					"fresh, address %s", address)
			}
		}
	})
}
