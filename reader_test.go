package geoip2

import (
	"fmt"
	. "launchpad.net/gocheck"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestGeoIP2(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestReader(c *C) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	defer reader.Close()

	record, err := reader.City(net.ParseIP("81.2.69.160"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}

	c.Assert(reader.Metadata.BinaryFormatMajorVersion, Equals, uint(2))
	c.Assert(reader.Metadata.BinaryFormatMinorVersion, Equals, uint(0))
	c.Assert(reader.Metadata.BuildEpoch, Equals, uint(1403110838))
	c.Assert(reader.Metadata.DatabaseType, Equals, "GeoIP2 City")
	c.Assert(reader.Metadata.Description, DeepEquals, map[string]string{
		"en": "GeoIP2 City Test Database (a small sample of real GeoIP2 data)",
		"zh": "小型数据库",
	})
	c.Assert(reader.Metadata.IPVersion, Equals, uint(6))
	c.Assert(reader.Metadata.Languages, DeepEquals, []string{"en", "zh"})
	c.Assert(reader.Metadata.NodeCount, Equals, uint(1218))
	c.Assert(reader.Metadata.RecordSize, Equals, uint(28))

	c.Assert(record.City.GeoNameID, Equals, uint(2643743))
	c.Assert(record.City.Names, DeepEquals, map[string]string{
		"de":    "London",
		"en":    "London",
		"es":    "Londres",
		"fr":    "Londres",
		"ja":    "ロンドン",
		"pt-BR": "Londres",
		"ru":    "Лондон",
	})
	c.Assert(record.Continent.GeoNameID, Equals, uint(6255148))
	c.Assert(record.Continent.Code, Equals, "EU")
	c.Assert(record.Continent.Names, DeepEquals, map[string]string{
		"de":    "Europa",
		"en":    "Europe",
		"es":    "Europa",
		"fr":    "Europe",
		"ja":    "ヨーロッパ",
		"pt-BR": "Europa",
		"ru":    "Европа",
		"zh-CN": "欧洲",
	})

	c.Assert(record.Country.GeoNameID, Equals, uint(2635167))
	c.Assert(record.Country.IsoCode, Equals, "GB")
	c.Assert(record.Country.Names, DeepEquals, map[string]string{
		"de":    "Vereinigtes Königreich",
		"en":    "United Kingdom",
		"es":    "Reino Unido",
		"fr":    "Royaume-Uni",
		"ja":    "イギリス",
		"pt-BR": "Reino Unido",
		"ru":    "Великобритания",
		"zh-CN": "英国",
	})

	c.Assert(record.Location.Latitude, Equals, 51.5142)
	c.Assert(record.Location.Longitude, Equals, -0.0931)
	c.Assert(record.Location.TimeZone, Equals, "Europe/London")

	c.Assert(record.Subdivisions[0].GeoNameID, Equals, uint(6269131))
	c.Assert(record.Subdivisions[0].IsoCode, Equals, "ENG")
	c.Assert(record.Subdivisions[0].Names, DeepEquals, map[string]string{
		"en":    "England",
		"pt-BR": "Inglaterra",
		"fr":    "Angleterre",
		"es":    "Inglaterra",
	})

	c.Assert(record.RegisteredCountry.GeoNameID, Equals, uint(6252001))
	c.Assert(record.RegisteredCountry.IsoCode, Equals, "US")
	c.Assert(record.RegisteredCountry.Names, DeepEquals, map[string]string{
		"de":    "USA",
		"en":    "United States",
		"es":    "Estados Unidos",
		"fr":    "États-Unis",
		"ja":    "アメリカ合衆国",
		"pt-BR": "Estados Unidos",
		"ru":    "США",
		"zh-CN": "美国",
	})
}

func (s *MySuite) TestMetroCode(c *C) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	defer reader.Close()

	record, err := reader.City(net.ParseIP("216.160.83.56"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}

	c.Assert(record.Location.MetroCode, Equals, uint(819))
}

func (s *MySuite) TestConnectionType(c *C) {
	reader, err := Open("test-data/test-data/GeoIP2-Connection-Type-Test.mmdb")
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	defer reader.Close()

	record, err := reader.ConnectionType(net.ParseIP("1.0.1.0"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	c.Assert(record.ConnectionType, Equals, "Cable/DSL")

}

func (s *MySuite) TestDomain(c *C) {
	reader, err := Open("test-data/test-data/GeoIP2-Domain-Test.mmdb")
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	defer reader.Close()

	record, err := reader.Domain(net.ParseIP("1.2.0.0"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	c.Assert(record.Domain, Equals, "maxmind.com")

}

func (s *MySuite) TestISP(c *C) {
	reader, err := Open("test-data/test-data/GeoIP2-ISP-Test.mmdb")
	if err != nil {
		c.Log(err)
		c.Fail()
	}
	defer reader.Close()

	record, err := reader.ISP(net.ParseIP("1.128.0.0"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}

	c.Assert(record.AutonomousSystemNumber, Equals, uint(1221))

	c.Assert(record.AutonomousSystemOrganization, Equals, "Telstra Pty Ltd")
	c.Assert(record.ISP, Equals, "Telstra Internet")
	c.Assert(record.Organization, Equals, "Telstra Internet")

}

func BenchmarkMaxMindDB(b *testing.B) {
	db, err := Open("GeoLite2-City.mmdb")
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		num := r.Uint32()
		ip := net.ParseIP(fmt.Sprintf("%d.%d.%d.%d", (0xFF000000&num)>>24,
			(0x00FF0000&num)>>16, (0x0000FF00&num)>>8, 0x000000F&num))
		_, err := db.City(ip)
		if err != nil {
			b.Fatal(err)
		}
	}
}
