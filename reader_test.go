package geoip2

import (
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	require.NoError(t, err)

	defer reader.Close()

	record, err := reader.City(net.ParseIP("81.2.69.160"))
	require.NoError(t, err)

	m := reader.Metadata()
	assert.Equal(t, uint(2), m.BinaryFormatMajorVersion)
	assert.Equal(t, uint(0), m.BinaryFormatMinorVersion)
	assert.NotZero(t, m.BuildEpoch)
	assert.Equal(t, "GeoIP2-City", m.DatabaseType)
	assert.Equal(t,
		map[string]string{
			"en": "GeoIP2 City Test Database (fake GeoIP2 data, for example purposes only)",
			"zh": "小型数据库",
		},
		m.Description,
	)
	assert.Equal(t, uint(6), m.IPVersion)
	assert.Equal(t, []string{"en", "zh"}, m.Languages)
	assert.NotZero(t, m.NodeCount)
	assert.Equal(t, uint(28), m.RecordSize)

	assert.Equal(t, uint(2643743), record.City.GeoNameID)
	assert.Equal(t,
		map[string]string{
			"de":    "London",
			"en":    "London",
			"es":    "Londres",
			"fr":    "Londres",
			"ja":    "ロンドン",
			"pt-BR": "Londres",
			"ru":    "Лондон",
		},
		record.City.Names,
	)
	assert.Equal(t, uint(6255148), record.Continent.GeoNameID)
	assert.Equal(t, "EU", record.Continent.Code)
	assert.Equal(t,
		map[string]string{
			"de":    "Europa",
			"en":    "Europe",
			"es":    "Europa",
			"fr":    "Europe",
			"ja":    "ヨーロッパ",
			"pt-BR": "Europa",
			"ru":    "Европа",
			"zh-CN": "欧洲",
		},
		record.Continent.Names,
	)

	assert.Equal(t, uint(2635167), record.Country.GeoNameID)
	assert.False(t, record.Country.IsInEuropeanUnion)
	assert.Equal(t, "GB", record.Country.IsoCode)
	assert.Equal(t,
		map[string]string{
			"de":    "Vereinigtes Königreich",
			"en":    "United Kingdom",
			"es":    "Reino Unido",
			"fr":    "Royaume-Uni",
			"ja":    "イギリス",
			"pt-BR": "Reino Unido",
			"ru":    "Великобритания",
			"zh-CN": "英国",
		},
		record.Country.Names,
	)

	assert.Equal(t, uint16(100), record.Location.AccuracyRadius)
	assert.Equal(t, 51.5142, record.Location.Latitude)
	assert.Equal(t, -0.0931, record.Location.Longitude)
	assert.Equal(t, "Europe/London", record.Location.TimeZone)

	assert.Equal(t, uint(6269131), record.Subdivisions[0].GeoNameID)
	assert.Equal(t, "ENG", record.Subdivisions[0].IsoCode)
	assert.Equal(t,
		map[string]string{
			"en":    "England",
			"pt-BR": "Inglaterra",
			"fr":    "Angleterre",
			"es":    "Inglaterra",
		},
		record.Subdivisions[0].Names,
	)

	assert.Equal(t, uint(6252001), record.RegisteredCountry.GeoNameID)
	assert.False(t, record.RegisteredCountry.IsInEuropeanUnion)
	assert.Equal(t, "US", record.RegisteredCountry.IsoCode)
	assert.Equal(t,
		map[string]string{
			"de":    "USA",
			"en":    "United States",
			"es":    "Estados Unidos",
			"fr":    "États-Unis",
			"ja":    "アメリカ合衆国",
			"pt-BR": "Estados Unidos",
			"ru":    "США",
			"zh-CN": "美国",
		},
		record.RegisteredCountry.Names,
	)

	assert.False(t, record.RepresentedCountry.IsInEuropeanUnion)
}

func TestMetroCode(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	require.NoError(t, err)
	defer reader.Close()

	record, err := reader.City(net.ParseIP("216.160.83.56"))
	require.NoError(t, err)

	assert.Equal(t, uint(819), record.Location.MetroCode)
}

func TestAnonymousIP(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Anonymous-IP-Test.mmdb")
	require.NoError(t, err)
	defer reader.Close()

	record, err := reader.AnonymousIP(net.ParseIP("1.2.0.0"))
	require.NoError(t, err)

	assert.True(t, record.IsAnonymous)

	assert.True(t, record.IsAnonymousVPN)
	assert.False(t, record.IsHostingProvider)
	assert.False(t, record.IsPublicProxy)
	assert.False(t, record.IsTorExitNode)
	assert.False(t, record.IsResidentialProxy)
}

func TestASN(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoLite2-ASN-Test.mmdb")
	require.NoError(t, err)
	defer reader.Close()

	record, err := reader.ASN(net.ParseIP("1.128.0.0"))
	require.NoError(t, err)

	assert.Equal(t, uint(1221), record.AutonomousSystemNumber)

	assert.Equal(t, "Telstra Pty Ltd", record.AutonomousSystemOrganization)
}

func TestConnectionType(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Connection-Type-Test.mmdb")
	require.NoError(t, err)

	defer reader.Close()

	record, err := reader.ConnectionType(net.ParseIP("1.0.1.0"))
	require.NoError(t, err)

	assert.Equal(t, "Cellular", record.ConnectionType)
}

func TestCountry(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Country-Test.mmdb")
	require.NoError(t, err)

	defer reader.Close()

	record, err := reader.Country(net.ParseIP("81.2.69.160"))
	require.NoError(t, err)

	assert.False(t, record.Country.IsInEuropeanUnion)
	assert.False(t, record.RegisteredCountry.IsInEuropeanUnion)
	assert.False(t, record.RepresentedCountry.IsInEuropeanUnion)
}

func TestDomain(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Domain-Test.mmdb")
	require.NoError(t, err)
	defer reader.Close()

	record, err := reader.Domain(net.ParseIP("1.2.0.0"))
	require.NoError(t, err)
	assert.Equal(t, "maxmind.com", record.Domain)
}

func TestEnterprise(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Enterprise-Test.mmdb")
	require.Nil(t, err)

	defer reader.Close()

	record, err := reader.Enterprise(net.ParseIP("74.209.24.0"))
	require.Nil(t, err)

	assert.Equal(t, uint8(11), record.City.Confidence)

	assert.Equal(t, uint(14671), record.Traits.AutonomousSystemNumber)
	assert.Equal(t, "FairPoint Communications", record.Traits.AutonomousSystemOrganization)
	assert.Equal(t, "Cable/DSL", record.Traits.ConnectionType)
	assert.Equal(t, "frpt.net", record.Traits.Domain)
	assert.Equal(t, float64(0.34), record.Traits.StaticIPScore)

	record, err = reader.Enterprise(net.ParseIP("149.101.100.0"))
	require.NoError(t, err)

	assert.Equal(t, uint(6167), record.Traits.AutonomousSystemNumber)

	assert.Equal(t, "CELLCO-PART", record.Traits.AutonomousSystemOrganization)
	assert.Equal(t, "Verizon Wireless", record.Traits.ISP)
	assert.Equal(t, "310", record.Traits.MobileCountryCode)
	assert.Equal(t, "004", record.Traits.MobileNetworkCode)
}

func TestISP(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-ISP-Test.mmdb")
	require.NoError(t, err)
	defer reader.Close()

	record, err := reader.ISP(net.ParseIP("149.101.100.0"))
	require.NoError(t, err)

	assert.Equal(t, uint(6167), record.AutonomousSystemNumber)

	assert.Equal(t, "CELLCO-PART", record.AutonomousSystemOrganization)
	assert.Equal(t, "Verizon Wireless", record.ISP)
	assert.Equal(t, "310", record.MobileCountryCode)
	assert.Equal(t, "004", record.MobileNetworkCode)
	assert.Equal(t, "Verizon Wireless", record.Organization)
}

// This ensures the compiler does not optimize away the function call.
var cityResult *City

func BenchmarkCity(b *testing.B) {
	db, err := Open("GeoLite2-City.mmdb")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	//nolint:gosec // this is just a benchmark
	r := rand.New(rand.NewSource(0))

	var city *City

	ip := make(net.IP, 4)
	for i := 0; i < b.N; i++ {
		randomIPv4Address(r, ip)
		city, err = db.City(ip)
		if err != nil {
			b.Fatal(err)
		}
	}
	cityResult = city
}

// This ensures the compiler does not optimize away the function call.
var asnResult *ASN

func BenchmarkASN(b *testing.B) {
	db, err := Open("GeoLite2-ASN.mmdb")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	//nolint:gosec // this is just a benchmark
	r := rand.New(rand.NewSource(0))

	var asn *ASN

	ip := make(net.IP, 4)
	for i := 0; i < b.N; i++ {
		randomIPv4Address(r, ip)
		asn, err = db.ASN(ip)
		if err != nil {
			b.Fatal(err)
		}
	}
	asnResult = asn
}

func randomIPv4Address(r *rand.Rand, ip net.IP) {
	num := r.Uint32()
	ip[0] = byte(num >> 24)
	ip[1] = byte(num >> 16)
	ip[2] = byte(num >> 8)
	ip[3] = byte(num)
}
