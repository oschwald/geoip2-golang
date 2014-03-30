package geoip2

import (
	. "launchpad.net/gocheck"
	"net"
	"testing"
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

	record, err := reader.City(net.ParseIP("81.2.69.160"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}

	c.Assert(record.City.GeoNameID, Equals, 2643743)
	c.Assert(record.City.Names, DeepEquals, map[string]string{
		"de":    "London",
		"en":    "London",
		"es":    "Londres",
		"fr":    "Londres",
		"ja":    "ロンドン",
		"pt-BR": "Londres",
		"ru":    "Лондон",
	})
	c.Assert(record.Continent.GeoNameID, Equals, 6255148)
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

	c.Assert(record.Country.GeoNameID, Equals, 2635167)
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

	c.Assert(record.Subdivisions[0].GeoNameID, Equals, 6269131)
	c.Assert(record.Subdivisions[0].IsoCode, Equals, "ENG")
	c.Assert(record.Subdivisions[0].Names, DeepEquals, map[string]string{
		"en":    "England",
		"pt-BR": "Inglaterra",
		"fr":    "Angleterre",
		"es":    "Inglaterra",
	})

	c.Assert(record.RegisteredCountry.GeoNameID, Equals, 6252001)
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

	reader.Close()
}

func (s *MySuite) TestMetroCode(c *C) {
	reader, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		c.Log(err)
		c.Fail()
	}

	record, err := reader.City(net.ParseIP("216.160.83.56"))
	if err != nil {
		c.Log(err)
		c.Fail()
	}

	c.Assert(record.Location.MetroCode, Equals, 819)

	reader.Close()
}
