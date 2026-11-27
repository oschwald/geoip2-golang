package geoip2

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2"
)

// These scalar-only models do not inherit the generated unmarshal methods.
// City and Enterprise use the deep mirrors in reflection_benchmark_test.go.
type asnReflectionBenchmark ASN

type ispReflectionBenchmark ISP

// BenchmarkRealDatabaseReflection measures fresh reflection decoding separately
// from search-tree traversal and the public Reader methods' generated decoders.
func BenchmarkRealDatabaseReflection(b *testing.B) {
	benchmarkRealReflection[cityReflectionBenchmark](b, "City", "GeoIP2-City.mmdb")
	benchmarkRealReflection[cityReflectionBenchmark](b, "GeoLiteCity", "GeoLite2-City.mmdb")
	benchmarkRealReflection[enterpriseReflectionBenchmark](
		b, "Enterprise", "GeoIP2-Enterprise.mmdb",
	)
	benchmarkRealReflection[asnReflectionBenchmark](b, "ASN", "GeoLite2-ASN.mmdb")
	benchmarkRealReflection[ispReflectionBenchmark](b, "ISP", "GeoIP2-ISP.mmdb")
}

func benchmarkRealReflection[T any](b *testing.B, name, file string) {
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
				for range 2 {
					for _, result := range results {
						require.NoError(b, result.Decode(new(T)))
					}
				}
				b.ReportAllocs()
				for i := 0; b.Loop(); i++ {
					if err := results[i%len(results)].Decode(new(T)); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	})
}
