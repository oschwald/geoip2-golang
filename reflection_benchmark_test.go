package geoip2

import "net/netip"

// These mirrors deliberately use distinct nested types. A shallow defined type
// such as `type cityReflection City` would now dispatch to the generated
// methods on City's nested model types and would no longer measure reflection.
type namesReflectionBenchmark struct {
	German              string `maxminddb:"de,maxsize:1024"`
	English             string `maxminddb:"en,maxsize:1024"`
	Spanish             string `maxminddb:"es,maxsize:1024"`
	French              string `maxminddb:"fr,maxsize:1024"`
	Japanese            string `maxminddb:"ja,maxsize:1024"`
	BrazilianPortuguese string `maxminddb:"pt-BR,maxsize:1024"`
	Russian             string `maxminddb:"ru,maxsize:1024"`
	SimplifiedChinese   string `maxminddb:"zh-CN,maxsize:1024"`
}

type continentReflectionBenchmark struct {
	Names     namesReflectionBenchmark `maxminddb:"names"`
	Code      string                   `maxminddb:"code,maxsize:16"`
	GeoNameID uint                     `maxminddb:"geoname_id"`
}

type locationReflectionBenchmark struct {
	Latitude       *float64 `maxminddb:"latitude"`
	Longitude      *float64 `maxminddb:"longitude"`
	TimeZone       string   `maxminddb:"time_zone,maxsize:256"`
	MetroCode      uint     `maxminddb:"metro_code"`
	AccuracyRadius uint16   `maxminddb:"accuracy_radius"`
}

type representedCountryReflectionBenchmark struct {
	Names             namesReflectionBenchmark `maxminddb:"names"`
	ISOCode           string                   `maxminddb:"iso_code,maxsize:16"`
	Type              string                   `maxminddb:"type,maxsize:128"`
	GeoNameID         uint                     `maxminddb:"geoname_id"`
	IsInEuropeanUnion bool                     `maxminddb:"is_in_european_union"`
}

type countryRecordReflectionBenchmark struct {
	Names             namesReflectionBenchmark `maxminddb:"names"`
	ISOCode           string                   `maxminddb:"iso_code,maxsize:16"`
	GeoNameID         uint                     `maxminddb:"geoname_id"`
	IsInEuropeanUnion bool                     `maxminddb:"is_in_european_union"`
}

type cityRecordReflectionBenchmark struct {
	Names     namesReflectionBenchmark `maxminddb:"names"`
	GeoNameID uint                     `maxminddb:"geoname_id"`
}

type cityPostalReflectionBenchmark struct {
	Code string `maxminddb:"code,maxsize:128"`
}

type citySubdivisionReflectionBenchmark struct {
	Names     namesReflectionBenchmark `maxminddb:"names"`
	ISOCode   string                   `maxminddb:"iso_code,maxsize:16"`
	GeoNameID uint                     `maxminddb:"geoname_id"`
}

type cityTraitsReflectionBenchmark struct {
	IPAddress netip.Addr   `maxminddb:"-"`
	Network   netip.Prefix `maxminddb:"-"`
	IsAnycast bool         `maxminddb:"is_anycast"`
}

type cityReflectionBenchmark struct {
	Traits             cityTraitsReflectionBenchmark         `maxminddb:"traits"`
	Postal             cityPostalReflectionBenchmark         `maxminddb:"postal"`
	Continent          continentReflectionBenchmark          `maxminddb:"continent"`
	City               cityRecordReflectionBenchmark         `maxminddb:"city"`
	Subdivisions       []citySubdivisionReflectionBenchmark  `maxminddb:"subdivisions,maxsize:32"`
	RepresentedCountry representedCountryReflectionBenchmark `maxminddb:"represented_country"`
	Country            countryRecordReflectionBenchmark      `maxminddb:"country"`
	RegisteredCountry  countryRecordReflectionBenchmark      `maxminddb:"registered_country"`
	Location           locationReflectionBenchmark           `maxminddb:"location"`
}

type enterpriseCityRecordReflectionBenchmark struct {
	Names      namesReflectionBenchmark `maxminddb:"names"`
	GeoNameID  uint                     `maxminddb:"geoname_id"`
	Confidence uint8                    `maxminddb:"confidence"`
}

type enterprisePostalReflectionBenchmark struct {
	Code       string `maxminddb:"code,maxsize:128"`
	Confidence uint8  `maxminddb:"confidence"`
}

type enterpriseSubdivisionReflectionBenchmark struct {
	Names      namesReflectionBenchmark `maxminddb:"names"`
	ISOCode    string                   `maxminddb:"iso_code,maxsize:16"`
	GeoNameID  uint                     `maxminddb:"geoname_id"`
	Confidence uint8                    `maxminddb:"confidence"`
}

type enterpriseCountryRecordReflectionBenchmark struct {
	Names             namesReflectionBenchmark `maxminddb:"names"`
	ISOCode           string                   `maxminddb:"iso_code,maxsize:16"`
	GeoNameID         uint                     `maxminddb:"geoname_id"`
	Confidence        uint8                    `maxminddb:"confidence"`
	IsInEuropeanUnion bool                     `maxminddb:"is_in_european_union"`
}

type enterpriseTraitsReflectionBenchmark struct {
	Network                      netip.Prefix `maxminddb:"-"`
	IPAddress                    netip.Addr   `maxminddb:"-"`
	AutonomousSystemOrganization string       `maxminddb:"autonomous_system_organization,maxsize:4096"`
	ConnectionType               string       `maxminddb:"connection_type,maxsize:128"`
	Domain                       string       `maxminddb:"domain,maxsize:512"`
	ISP                          string       `maxminddb:"isp,maxsize:4096"`
	MobileCountryCode            string       `maxminddb:"mobile_country_code,maxsize:16"`
	MobileNetworkCode            string       `maxminddb:"mobile_network_code,maxsize:16"`
	Organization                 string       `maxminddb:"organization,maxsize:4096"`
	UserType                     string       `maxminddb:"user_type,maxsize:128"`
	StaticIPScore                float64      `maxminddb:"static_ip_score"`
	AutonomousSystemNumber       uint         `maxminddb:"autonomous_system_number"`
	IsAnycast                    bool         `maxminddb:"is_anycast"`
	IsLegitimateProxy            bool         `maxminddb:"is_legitimate_proxy"`
}

type enterpriseReflectionBenchmark struct {
	Continent          continentReflectionBenchmark               `maxminddb:"continent"`
	Subdivisions       []enterpriseSubdivisionReflectionBenchmark `maxminddb:"subdivisions,maxsize:32"`
	Postal             enterprisePostalReflectionBenchmark        `maxminddb:"postal"`
	RepresentedCountry representedCountryReflectionBenchmark      `maxminddb:"represented_country"`
	Country            enterpriseCountryRecordReflectionBenchmark `maxminddb:"country"`
	RegisteredCountry  countryRecordReflectionBenchmark           `maxminddb:"registered_country"`
	City               enterpriseCityRecordReflectionBenchmark    `maxminddb:"city"`
	Location           locationReflectionBenchmark                `maxminddb:"location"`
	Traits             enterpriseTraitsReflectionBenchmark        `maxminddb:"traits"`
}

func cityFromReflection(value cityReflectionBenchmark) City {
	var subdivisions []CitySubdivision
	if value.Subdivisions != nil {
		subdivisions = make([]CitySubdivision, len(value.Subdivisions))
		for i, subdivision := range value.Subdivisions {
			subdivisions[i] = CitySubdivision{
				Names:     Names(subdivision.Names),
				ISOCode:   subdivision.ISOCode,
				GeoNameID: subdivision.GeoNameID,
			}
		}
	}
	return City{
		Traits:    CityTraits(value.Traits),
		Postal:    CityPostal(value.Postal),
		Continent: continentFromReflection(value.Continent),
		City: CityRecord{
			Names:     Names(value.City.Names),
			GeoNameID: value.City.GeoNameID,
		},
		Subdivisions:       subdivisions,
		RepresentedCountry: representedCountryFromReflection(value.RepresentedCountry),
		Country:            countryRecordFromReflection(value.Country),
		RegisteredCountry:  countryRecordFromReflection(value.RegisteredCountry),
		Location:           Location(value.Location),
	}
}

func enterpriseFromReflection(value enterpriseReflectionBenchmark) Enterprise {
	var subdivisions []EnterpriseSubdivision
	if value.Subdivisions != nil {
		subdivisions = make([]EnterpriseSubdivision, len(value.Subdivisions))
		for i, subdivision := range value.Subdivisions {
			subdivisions[i] = EnterpriseSubdivision{
				Names:      Names(subdivision.Names),
				ISOCode:    subdivision.ISOCode,
				GeoNameID:  subdivision.GeoNameID,
				Confidence: subdivision.Confidence,
			}
		}
	}
	return Enterprise{
		Continent:    continentFromReflection(value.Continent),
		Subdivisions: subdivisions,
		Postal: EnterprisePostal{
			Code:       value.Postal.Code,
			Confidence: value.Postal.Confidence,
		},
		RepresentedCountry: representedCountryFromReflection(value.RepresentedCountry),
		Country: EnterpriseCountryRecord{
			Names:             Names(value.Country.Names),
			ISOCode:           value.Country.ISOCode,
			GeoNameID:         value.Country.GeoNameID,
			Confidence:        value.Country.Confidence,
			IsInEuropeanUnion: value.Country.IsInEuropeanUnion,
		},
		RegisteredCountry: countryRecordFromReflection(value.RegisteredCountry),
		City: EnterpriseCityRecord{
			Names:      Names(value.City.Names),
			GeoNameID:  value.City.GeoNameID,
			Confidence: value.City.Confidence,
		},
		Location: Location(value.Location),
		Traits:   EnterpriseTraits(value.Traits),
	}
}

func continentFromReflection(value continentReflectionBenchmark) Continent {
	return Continent{
		Names:     Names(value.Names),
		Code:      value.Code,
		GeoNameID: value.GeoNameID,
	}
}

func representedCountryFromReflection(
	value representedCountryReflectionBenchmark,
) RepresentedCountry {
	return RepresentedCountry{
		Names:             Names(value.Names),
		ISOCode:           value.ISOCode,
		Type:              value.Type,
		GeoNameID:         value.GeoNameID,
		IsInEuropeanUnion: value.IsInEuropeanUnion,
	}
}

func countryRecordFromReflection(value countryRecordReflectionBenchmark) CountryRecord {
	return CountryRecord{
		Names:             Names(value.Names),
		ISOCode:           value.ISOCode,
		GeoNameID:         value.GeoNameID,
		IsInEuropeanUnion: value.IsInEuropeanUnion,
	}
}
