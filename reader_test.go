package geoip2

import (
	"math/rand"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	assert.Nil(t, err)

	defer reader.Close()

	record, err := reader.City(net.ParseIP("81.2.69.160"))
	assert.Nil(t, err)

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
	assert.True(t, record.Country.IsInEuropeanUnion)
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
	assert.Nil(t, err)
	defer reader.Close()

	record, err := reader.City(net.ParseIP("216.160.83.56"))
	assert.Nil(t, err)

	assert.Equal(t, uint(819), record.Location.MetroCode)
}

func TestAnonymousIP(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Anonymous-IP-Test.mmdb")
	assert.Nil(t, err)
	defer reader.Close()

	record, err := reader.AnonymousIP(net.ParseIP("1.2.0.0"))
	assert.Nil(t, err)

	assert.Equal(t, true, record.IsAnonymous)

	assert.Equal(t, true, record.IsAnonymousVPN)
	assert.Equal(t, false, record.IsHostingProvider)
	assert.Equal(t, false, record.IsPublicProxy)
	assert.Equal(t, false, record.IsTorExitNode)
}

func TestASN(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoLite2-ASN-Test.mmdb")
	assert.Nil(t, err)
	defer reader.Close()

	record, err := reader.ASN(net.ParseIP("1.128.0.0"))
	assert.Nil(t, err)

	assert.Equal(t, uint(1221), record.AutonomousSystemNumber)

	assert.Equal(t, "Telstra Pty Ltd", record.AutonomousSystemOrganization)
}

func TestConnectionType(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Connection-Type-Test.mmdb")
	assert.Nil(t, err)

	defer reader.Close()

	record, err := reader.ConnectionType(net.ParseIP("1.0.1.0"))
	assert.Nil(t, err)

	assert.Equal(t, "Cable/DSL", record.ConnectionType)
}

func TestCountry(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Country-Test.mmdb")
	assert.Nil(t, err)

	defer reader.Close()

	record, err := reader.Country(net.ParseIP("81.2.69.160"))
	assert.Nil(t, err)

	assert.True(t, record.Country.IsInEuropeanUnion)
	assert.False(t, record.RegisteredCountry.IsInEuropeanUnion)
	assert.False(t, record.RepresentedCountry.IsInEuropeanUnion)
}

func TestDomain(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-Domain-Test.mmdb")
	assert.Nil(t, err)
	defer reader.Close()

	record, err := reader.Domain(net.ParseIP("1.2.0.0"))
	assert.Nil(t, err)
	assert.Equal(t, "maxmind.com", record.Domain)
}

func TestISP(t *testing.T) {
	reader, err := Open("test-data/test-data/GeoIP2-ISP-Test.mmdb")
	assert.Nil(t, err)
	defer reader.Close()

	record, err := reader.ISP(net.ParseIP("1.128.0.0"))
	assert.Nil(t, err)

	assert.Equal(t, uint(1221), record.AutonomousSystemNumber)

	assert.Equal(t, "Telstra Pty Ltd", record.AutonomousSystemOrganization)
	assert.Equal(t, "Telstra Internet", record.ISP)
	assert.Equal(t, "Telstra Internet", record.Organization)
}

// This ensures the compiler does not optimize away the function call
var cityResult *City

func BenchmarkMaxMindDB(b *testing.B) {
	db, err := Open("GeoLite2-City.mmdb")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	r := rand.New(rand.NewSource(0))

	var city *City

	for i := 0; i < b.N; i++ {
		ip := randomIPv4Address(b, r)
		city, err = db.City(ip)
		if err != nil {
			b.Fatal(err)
		}
	}
	cityResult = city
}

func randomIPv4Address(b *testing.B, r *rand.Rand) net.IP {
	num := r.Uint32()
	return []byte{byte(num >> 24), byte(num >> 16), byte(num >> 8),
		byte(num)}
}
