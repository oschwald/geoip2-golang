package geoip2

import "net/netip"

// These mirrors deliberately use distinct nested types. A shallow defined type
// such as `type cityReflection City` would now dispatch to the generated
// methods on City's nested model types and would no longer measure reflection.
type namesReflectionBenchmark struct {
	German              string `maxminddb:"de"`
	English             string `maxminddb:"en"`
	Spanish             string `maxminddb:"es"`
	French              string `maxminddb:"fr"`
	Japanese            string `maxminddb:"ja"`
	BrazilianPortuguese string `maxminddb:"pt-BR"`
	Russian             string `maxminddb:"ru"`
	SimplifiedChinese   string `maxminddb:"zh-CN"`
}

type continentReflectionBenchmark struct {
	Names     namesReflectionBenchmark `maxminddb:"names"`
	Code      string                   `maxminddb:"code"`
	GeoNameID uint                     `maxminddb:"geoname_id"`
}

type locationReflectionBenchmark struct {
	Latitude       *float64 `maxminddb:"latitude"`
	Longitude      *float64 `maxminddb:"longitude"`
	TimeZone       string   `maxminddb:"time_zone"`
	MetroCode      uint     `maxminddb:"metro_code"`
	AccuracyRadius uint16   `maxminddb:"accuracy_radius"`
}

type representedCountryReflectionBenchmark struct {
	Names             namesReflectionBenchmark `maxminddb:"names"`
	ISOCode           string                   `maxminddb:"iso_code"`
	Type              string                   `maxminddb:"type"`
	GeoNameID         uint                     `maxminddb:"geoname_id"`
	IsInEuropeanUnion bool                     `maxminddb:"is_in_european_union"`
}

type countryRecordReflectionBenchmark struct {
	Names             namesReflectionBenchmark `maxminddb:"names"`
	ISOCode           string                   `maxminddb:"iso_code"`
	GeoNameID         uint                     `maxminddb:"geoname_id"`
	IsInEuropeanUnion bool                     `maxminddb:"is_in_european_union"`
}

type cityRecordReflectionBenchmark struct {
	Names     namesReflectionBenchmark `maxminddb:"names"`
	GeoNameID uint                     `maxminddb:"geoname_id"`
}

type cityPostalReflectionBenchmark struct {
	Code string `maxminddb:"code"`
}

type citySubdivisionReflectionBenchmark struct {
	Names     namesReflectionBenchmark `maxminddb:"names"`
	ISOCode   string                   `maxminddb:"iso_code"`
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
	Subdivisions       []citySubdivisionReflectionBenchmark  `maxminddb:"subdivisions"`
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
	Code       string `maxminddb:"code"`
	Confidence uint8  `maxminddb:"confidence"`
}

type enterpriseSubdivisionReflectionBenchmark struct {
	Names      namesReflectionBenchmark `maxminddb:"names"`
	ISOCode    string                   `maxminddb:"iso_code"`
	GeoNameID  uint                     `maxminddb:"geoname_id"`
	Confidence uint8                    `maxminddb:"confidence"`
}

type enterpriseCountryRecordReflectionBenchmark struct {
	Names             namesReflectionBenchmark `maxminddb:"names"`
	ISOCode           string                   `maxminddb:"iso_code"`
	GeoNameID         uint                     `maxminddb:"geoname_id"`
	Confidence        uint8                    `maxminddb:"confidence"`
	IsInEuropeanUnion bool                     `maxminddb:"is_in_european_union"`
}

type enterpriseTraitsReflectionBenchmark struct {
	Network                      netip.Prefix `maxminddb:"-"`
	IPAddress                    netip.Addr   `maxminddb:"-"`
	AutonomousSystemOrganization string       `maxminddb:"autonomous_system_organization"`
	ConnectionType               string       `maxminddb:"connection_type"`
	Domain                       string       `maxminddb:"domain"`
	ISP                          string       `maxminddb:"isp"`
	MobileCountryCode            string       `maxminddb:"mobile_country_code"`
	MobileNetworkCode            string       `maxminddb:"mobile_network_code"`
	Organization                 string       `maxminddb:"organization"`
	UserType                     string       `maxminddb:"user_type"`
	StaticIPScore                float64      `maxminddb:"static_ip_score"`
	AutonomousSystemNumber       uint         `maxminddb:"autonomous_system_number"`
	IsAnycast                    bool         `maxminddb:"is_anycast"`
	IsLegitimateProxy            bool         `maxminddb:"is_legitimate_proxy"`
}

type enterpriseReflectionBenchmark struct {
	Continent          continentReflectionBenchmark               `maxminddb:"continent"`
	Subdivisions       []enterpriseSubdivisionReflectionBenchmark `maxminddb:"subdivisions"`
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
