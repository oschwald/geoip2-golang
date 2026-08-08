package geoip2

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2/mmdbdata"
)

var (
	_ mmdbdata.CursorUnmarshaler = (*City)(nil)
	_ mmdbdata.CursorUnmarshaler = (*Enterprise)(nil)
)

var (
	generatedCitySink   *City
	reusedCitySink      City
	benchmarkCitySink   *cityReflectionBenchmark
	reusedBenchCitySink cityReflectionBenchmark
	generatedEntSink    *Enterprise
	reusedEntSink       Enterprise
	benchmarkEntSink    *enterpriseReflectionBenchmark
	reusedBenchEntSink  enterpriseReflectionBenchmark
)

func TestGeneratedCityParity(t *testing.T) {
	testGeneratedParity(t, "test-data/test-data/GeoIP2-City-Test.mmdb",
		cityBenchmarkAddresses(), cityFromReflection)
}

func TestGeneratedEnterpriseParity(t *testing.T) {
	testGeneratedParity(t, "test-data/test-data/GeoIP2-Enterprise-Test.mmdb",
		enterpriseBenchmarkAddresses(), enterpriseFromReflection)
}

func testGeneratedParity[T, R any](
	t *testing.T,
	database string,
	addresses []netip.Addr,
	fromReflection func(R) T,
) {
	t.Helper()
	reader, err := Open(database)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reader.Close()) })

	var generated T
	var reflection R
	seed := reader.mmdbReader.Lookup(addresses[len(addresses)-1])
	require.True(t, seed.Found())
	require.NoError(t, seed.Decode(&generated))
	require.NoError(t, seed.Decode(&reflection))
	require.Equal(t, fromReflection(reflection), generated, "dirty-destination seed")
	for iteration := range 3 {
		for _, address := range addresses {
			result := reader.mmdbReader.Lookup(address)
			require.True(t, result.Found(), "address %s", address)
			require.NoError(t, result.Decode(&generated))
			require.NoError(t, result.Decode(&reflection))
			require.Equal(t, fromReflection(reflection), generated,
				"iteration %d, address %s", iteration, address)
		}
	}
}

func BenchmarkCityDecodeGeneratedReused(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := cityBenchmarkAddresses()
	results := make([]resultDecode, len(addresses))
	for i, address := range addresses {
		result := reader.mmdbReader.Lookup(address)
		results[i] = resultDecode{decode: result.Decode}
	}
	var city City
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		if err := results[i%len(results)].decode(&city); err != nil {
			b.Fatal(err)
		}
	}
	reusedCitySink = city
}

func BenchmarkCityDecodeReflectionReused(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := cityBenchmarkAddresses()
	results := make([]resultDecode, len(addresses))
	for i, address := range addresses {
		result := reader.mmdbReader.Lookup(address)
		results[i] = resultDecode{decode: result.Decode}
	}
	var city cityReflectionBenchmark
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		if err := results[i%len(results)].decode(&city); err != nil {
			b.Fatal(err)
		}
	}
	reusedBenchCitySink = city
}

func BenchmarkCityLookupGeneratedFresh(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := cityBenchmarkAddresses()
	var city *City
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		city, err = reader.City(addresses[i%len(addresses)])
		if err != nil {
			b.Fatal(err)
		}
	}
	generatedCitySink = city
}

func BenchmarkCityLookupReflectionFresh(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := cityBenchmarkAddresses()
	var city *cityReflectionBenchmark
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		address := addresses[i%len(addresses)]
		result := reader.mmdbReader.Lookup(address)
		city = new(cityReflectionBenchmark)
		if err = result.Decode(city); err != nil {
			b.Fatal(err)
		}
		city.Traits.IPAddress = address
		city.Traits.Network = result.Prefix()
	}
	benchmarkCitySink = city
}

func BenchmarkCityCommercialGeneratedFresh(b *testing.B) {
	reader, err := Open("/var/lib/GeoIP/GeoIP2-City.mmdb")
	if err != nil {
		b.Skip(err)
	}
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	address := netip.MustParseAddr("128.101.101.101")
	var city *City
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		city, err = reader.City(address)
		if err != nil {
			b.Fatal(err)
		}
	}
	generatedCitySink = city
}

func BenchmarkEnterpriseDecodeGeneratedReused(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-Enterprise-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := enterpriseBenchmarkAddresses()
	results := make([]resultDecode, len(addresses))
	for i, address := range addresses {
		result := reader.mmdbReader.Lookup(address)
		results[i] = resultDecode{decode: result.Decode}
	}
	var enterprise Enterprise
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		if err := results[i%len(results)].decode(&enterprise); err != nil {
			b.Fatal(err)
		}
	}
	reusedEntSink = enterprise
}

func BenchmarkEnterpriseDecodeReflectionReused(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-Enterprise-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := enterpriseBenchmarkAddresses()
	results := make([]resultDecode, len(addresses))
	for i, address := range addresses {
		result := reader.mmdbReader.Lookup(address)
		results[i] = resultDecode{decode: result.Decode}
	}
	var enterprise enterpriseReflectionBenchmark
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		if err := results[i%len(results)].decode(&enterprise); err != nil {
			b.Fatal(err)
		}
	}
	reusedBenchEntSink = enterprise
}

func BenchmarkEnterpriseLookupGeneratedFresh(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-Enterprise-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := enterpriseBenchmarkAddresses()
	var enterprise *Enterprise
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		enterprise, err = reader.Enterprise(addresses[i%len(addresses)])
		if err != nil {
			b.Fatal(err)
		}
	}
	generatedEntSink = enterprise
}

func BenchmarkEnterpriseLookupReflectionFresh(b *testing.B) {
	reader, err := Open("test-data/test-data/GeoIP2-Enterprise-Test.mmdb")
	require.NoError(b, err)
	b.Cleanup(func() { require.NoError(b, reader.Close()) })
	addresses := enterpriseBenchmarkAddresses()
	var enterprise *enterpriseReflectionBenchmark
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		address := addresses[i%len(addresses)]
		result := reader.mmdbReader.Lookup(address)
		enterprise = new(enterpriseReflectionBenchmark)
		if err = result.Decode(enterprise); err != nil {
			b.Fatal(err)
		}
		enterprise.Traits.IPAddress = address
		enterprise.Traits.Network = result.Prefix()
	}
	benchmarkEntSink = enterprise
}

type resultDecode struct {
	decode func(any) error
}

func cityBenchmarkAddresses() []netip.Addr {
	return []netip.Addr{
		netip.MustParseAddr("2.125.160.216"),
		netip.MustParseAddr("67.43.156.0"),
		netip.MustParseAddr("81.2.69.160"),
		netip.MustParseAddr("89.160.20.128"),
		netip.MustParseAddr("175.16.199.0"),
		netip.MustParseAddr("202.196.224.0"),
		netip.MustParseAddr("216.160.83.56"),
	}
}

func enterpriseBenchmarkAddresses() []netip.Addr {
	return []netip.Addr{
		netip.MustParseAddr("74.209.24.0"),
		netip.MustParseAddr("81.2.69.160"),
		netip.MustParseAddr("149.101.100.0"),
		netip.MustParseAddr("175.16.199.0"),
		netip.MustParseAddr("202.196.224.0"),
		netip.MustParseAddr("214.1.1.0"),
		netip.MustParseAddr("216.160.83.56"),
	}
}
