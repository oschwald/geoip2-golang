// Package geoip2 provides an easy-to-use API for the MaxMind GeoIP2 and
// GeoLite2 databases; this package does not support GeoIP Legacy databases.
//
// The structs provided by this package match the internal structure of
// the data in the MaxMind databases.
//
// See github.com/oschwald/maxminddb-golang for more advanced used cases.
package geoip2

import (
	"fmt"
	"net/netip"

	"github.com/oschwald/maxminddb-golang/v2"
)

// The Enterprise struct corresponds to the data in the GeoIP2 Enterprise
// database.
type Enterprise struct {
	Continent struct {
		Names     map[string]string `maxminddb:"names"`
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"continent"`
	City struct {
		Names      map[string]string `maxminddb:"names"`
		GeoNameID  uint              `maxminddb:"geoname_id"`
		Confidence uint8             `maxminddb:"confidence"`
	} `maxminddb:"city"`
	Postal struct {
		Code       string `maxminddb:"code"`
		Confidence uint8  `maxminddb:"confidence"`
	} `maxminddb:"postal"`
	Subdivisions []struct {
		Names      map[string]string `maxminddb:"names"`
		ISOCode    string            `maxminddb:"iso_code"`
		GeoNameID  uint              `maxminddb:"geoname_id"`
		Confidence uint8             `maxminddb:"confidence"`
	} `maxminddb:"subdivisions"`
	RepresentedCountry struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		Type              string            `maxminddb:"type"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"represented_country"`
	Country struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		Confidence        uint8             `maxminddb:"confidence"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		Confidence        uint8             `maxminddb:"confidence"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	Traits struct {
		AutonomousSystemOrganization string  `maxminddb:"autonomous_system_organization"`
		ConnectionType               string  `maxminddb:"connection_type"`
		Domain                       string  `maxminddb:"domain"`
		ISP                          string  `maxminddb:"isp"`
		MobileCountryCode            string  `maxminddb:"mobile_country_code"`
		MobileNetworkCode            string  `maxminddb:"mobile_network_code"`
		Organization                 string  `maxminddb:"organization"`
		UserType                     string  `maxminddb:"user_type"`
		AutonomousSystemNumber       uint    `maxminddb:"autonomous_system_number"`
		StaticIPScore                float64 `maxminddb:"static_ip_score"`
		IsAnonymousProxy             bool    `maxminddb:"is_anonymous_proxy"`
		IsAnycast                    bool    `maxminddb:"is_anycast"`
		IsLegitimateProxy            bool    `maxminddb:"is_legitimate_proxy"`
		IsSatelliteProvider          bool    `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
	Location struct {
		TimeZone       string  `maxminddb:"time_zone"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for Enterprise
func (e *Enterprise) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "continent":
			for contKey, contErr := range d.DecodeMap() {
				if contErr != nil {
					return contErr
				}
				switch string(contKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					e.Continent.Names = names
				case "code":
					code, codeErr := d.DecodeString()
					if codeErr != nil {
						return codeErr
					}
					e.Continent.Code = code
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					e.Continent.GeoNameID = uint(geoID)
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "city":
			for cityKey, cityErr := range d.DecodeMap() {
				if cityErr != nil {
					return cityErr
				}
				switch string(cityKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					e.City.Names = names
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					e.City.GeoNameID = uint(geoID)
				case "confidence":
					conf, confErr := d.DecodeUInt16()
					if confErr != nil {
						return confErr
					}
					e.City.Confidence = uint8(conf)
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "postal":
			for postalKey, postalErr := range d.DecodeMap() {
				if postalErr != nil {
					return postalErr
				}
				switch string(postalKey) {
				case "code":
					code, codeErr := d.DecodeString()
					if codeErr != nil {
						return codeErr
					}
					e.Postal.Code = code
				case "confidence":
					conf, confErr := d.DecodeUInt16()
					if confErr != nil {
						return confErr
					}
					e.Postal.Confidence = uint8(conf)
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "subdivisions":
			var subdivisions []struct {
				Names      map[string]string `maxminddb:"names"`
				ISOCode    string            `maxminddb:"iso_code"`
				GeoNameID  uint              `maxminddb:"geoname_id"`
				Confidence uint8             `maxminddb:"confidence"`
			}

			for sliceErr := range d.DecodeSlice() {
				if sliceErr != nil {
					return sliceErr
				}

				var subdivision struct {
					Names      map[string]string `maxminddb:"names"`
					ISOCode    string            `maxminddb:"iso_code"`
					GeoNameID  uint              `maxminddb:"geoname_id"`
					Confidence uint8             `maxminddb:"confidence"`
				}

				for subKey, subErr := range d.DecodeMap() {
					if subErr != nil {
						return subErr
					}
					switch string(subKey) {
					case "names":
						names := make(map[string]string)
						for nameKey, nameErr := range d.DecodeMap() {
							if nameErr != nil {
								return nameErr
							}
							value, valueErr := d.DecodeString()
							if valueErr != nil {
								return valueErr
							}
							names[string(nameKey)] = value
						}
						subdivision.Names = names
					case "iso_code":
						isoCode, isoErr := d.DecodeString()
						if isoErr != nil {
							return isoErr
						}
						subdivision.ISOCode = isoCode
					case "geoname_id":
						geoID, geoErr := d.DecodeUInt32()
						if geoErr != nil {
							return geoErr
						}
						subdivision.GeoNameID = uint(geoID)
					case "confidence":
						conf, confErr := d.DecodeUInt16()
						if confErr != nil {
							return confErr
						}
						subdivision.Confidence = uint8(conf)
					default:
						if err := d.SkipValue(); err != nil {
							return err
						}
					}
				}
				subdivisions = append(subdivisions, subdivision)
			}
			e.Subdivisions = subdivisions
		case "represented_country":
			for repKey, repErr := range d.DecodeMap() {
				if repErr != nil {
					return repErr
				}
				switch string(repKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					e.RepresentedCountry.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					e.RepresentedCountry.ISOCode = isoCode
				case "type":
					typeStr, typeErr := d.DecodeString()
					if typeErr != nil {
						return typeErr
					}
					e.RepresentedCountry.Type = typeStr
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					e.RepresentedCountry.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.RepresentedCountry.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "country":
			for countryKey, countryErr := range d.DecodeMap() {
				if countryErr != nil {
					return countryErr
				}
				switch string(countryKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					e.Country.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					e.Country.ISOCode = isoCode
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					e.Country.GeoNameID = uint(geoID)
				case "confidence":
					conf, confErr := d.DecodeUInt16()
					if confErr != nil {
						return confErr
					}
					e.Country.Confidence = uint8(conf)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.Country.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "registered_country":
			for regKey, regErr := range d.DecodeMap() {
				if regErr != nil {
					return regErr
				}
				switch string(regKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					e.RegisteredCountry.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					e.RegisteredCountry.ISOCode = isoCode
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					e.RegisteredCountry.GeoNameID = uint(geoID)
				case "confidence":
					conf, confErr := d.DecodeUInt16()
					if confErr != nil {
						return confErr
					}
					e.RegisteredCountry.Confidence = uint8(conf)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.RegisteredCountry.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "traits":
			for traitsKey, traitsErr := range d.DecodeMap() {
				if traitsErr != nil {
					return traitsErr
				}
				switch string(traitsKey) {
				case "autonomous_system_organization":
					org, orgErr := d.DecodeString()
					if orgErr != nil {
						return orgErr
					}
					e.Traits.AutonomousSystemOrganization = org
				case "connection_type":
					connType, connErr := d.DecodeString()
					if connErr != nil {
						return connErr
					}
					e.Traits.ConnectionType = connType
				case "domain":
					domain, domainErr := d.DecodeString()
					if domainErr != nil {
						return domainErr
					}
					e.Traits.Domain = domain
				case "isp":
					isp, ispErr := d.DecodeString()
					if ispErr != nil {
						return ispErr
					}
					e.Traits.ISP = isp
				case "mobile_country_code":
					mcc, mccErr := d.DecodeString()
					if mccErr != nil {
						return mccErr
					}
					e.Traits.MobileCountryCode = mcc
				case "mobile_network_code":
					mnc, mncErr := d.DecodeString()
					if mncErr != nil {
						return mncErr
					}
					e.Traits.MobileNetworkCode = mnc
				case "organization":
					org, orgErr := d.DecodeString()
					if orgErr != nil {
						return orgErr
					}
					e.Traits.Organization = org
				case "user_type":
					userType, userErr := d.DecodeString()
					if userErr != nil {
						return userErr
					}
					e.Traits.UserType = userType
				case "autonomous_system_number":
					asn, asnErr := d.DecodeUInt32()
					if asnErr != nil {
						return asnErr
					}
					e.Traits.AutonomousSystemNumber = uint(asn)
				case "static_ip_score":
					score, scoreErr := d.DecodeFloat64()
					if scoreErr != nil {
						return scoreErr
					}
					e.Traits.StaticIPScore = score
				case "is_anonymous_proxy":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.Traits.IsAnonymousProxy = val
				case "is_anycast":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.Traits.IsAnycast = val
				case "is_legitimate_proxy":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.Traits.IsLegitimateProxy = val
				case "is_satellite_provider":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					e.Traits.IsSatelliteProvider = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "location":
			for locKey, locErr := range d.DecodeMap() {
				if locErr != nil {
					return locErr
				}
				switch string(locKey) {
				case "time_zone":
					tz, tzErr := d.DecodeString()
					if tzErr != nil {
						return tzErr
					}
					e.Location.TimeZone = tz
				case "latitude":
					lat, latErr := d.DecodeFloat64()
					if latErr != nil {
						return latErr
					}
					e.Location.Latitude = lat
				case "longitude":
					lon, lonErr := d.DecodeFloat64()
					if lonErr != nil {
						return lonErr
					}
					e.Location.Longitude = lon
				case "metro_code":
					metro, metroErr := d.DecodeUInt16()
					if metroErr != nil {
						return metroErr
					}
					e.Location.MetroCode = uint(metro)
				case "accuracy_radius":
					radius, radiusErr := d.DecodeUInt16()
					if radiusErr != nil {
						return radiusErr
					}
					e.Location.AccuracyRadius = radius
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The City struct corresponds to the data in the GeoIP2/GeoLite2 City
// databases.
type City struct {
	City struct {
		Names     map[string]string `maxminddb:"names"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"city"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	Continent struct {
		Names     map[string]string `maxminddb:"names"`
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"continent"`
	Subdivisions []struct {
		Names     map[string]string `maxminddb:"names"`
		ISOCode   string            `maxminddb:"iso_code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"subdivisions"`
	RepresentedCountry struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		Type              string            `maxminddb:"type"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"represented_country"`
	Country struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	Location struct {
		TimeZone       string  `maxminddb:"time_zone"`
		Latitude       float64 `maxminddb:"latitude"`
		Longitude      float64 `maxminddb:"longitude"`
		MetroCode      uint    `maxminddb:"metro_code"`
		AccuracyRadius uint16  `maxminddb:"accuracy_radius"`
	} `maxminddb:"location"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsAnycast           bool `maxminddb:"is_anycast"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for City.
func (c *City) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "city":
			for cityKey, cityErr := range d.DecodeMap() {
				if cityErr != nil {
					return cityErr
				}
				switch string(cityKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.City.Names = names
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.City.GeoNameID = uint(geoID)
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "postal":
			for postalKey, postalErr := range d.DecodeMap() {
				if postalErr != nil {
					return postalErr
				}
				switch string(postalKey) {
				case "code":
					code, codeErr := d.DecodeString()
					if codeErr != nil {
						return codeErr
					}
					c.Postal.Code = code
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "continent":
			for contKey, contErr := range d.DecodeMap() {
				if contErr != nil {
					return contErr
				}
				switch string(contKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.Continent.Names = names
				case "code":
					code, codeErr := d.DecodeString()
					if codeErr != nil {
						return codeErr
					}
					c.Continent.Code = code
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.Continent.GeoNameID = uint(geoID)
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "subdivisions":
			// This is the complex part - handle slice of structs
			var subdivisions []struct {
				Names     map[string]string `maxminddb:"names"`
				ISOCode   string            `maxminddb:"iso_code"`
				GeoNameID uint              `maxminddb:"geoname_id"`
			}

			for sliceErr := range d.DecodeSlice() {
				if sliceErr != nil {
					return sliceErr
				}

				var subdivision struct {
					Names     map[string]string `maxminddb:"names"`
					ISOCode   string            `maxminddb:"iso_code"`
					GeoNameID uint              `maxminddb:"geoname_id"`
				}

				for subKey, subErr := range d.DecodeMap() {
					if subErr != nil {
						return subErr
					}
					switch string(subKey) {
					case "names":
						names := make(map[string]string)
						for nameKey, nameErr := range d.DecodeMap() {
							if nameErr != nil {
								return nameErr
							}
							value, valueErr := d.DecodeString()
							if valueErr != nil {
								return valueErr
							}
							names[string(nameKey)] = value
						}
						subdivision.Names = names
					case "iso_code":
						isoCode, isoErr := d.DecodeString()
						if isoErr != nil {
							return isoErr
						}
						subdivision.ISOCode = isoCode
					case "geoname_id":
						geoID, geoErr := d.DecodeUInt32()
						if geoErr != nil {
							return geoErr
						}
						subdivision.GeoNameID = uint(geoID)
					default:
						if err := d.SkipValue(); err != nil {
							return err
						}
					}
				}
				subdivisions = append(subdivisions, subdivision)
			}
			c.Subdivisions = subdivisions
		case "represented_country":
			for repKey, repErr := range d.DecodeMap() {
				if repErr != nil {
					return repErr
				}
				switch string(repKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.RepresentedCountry.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					c.RepresentedCountry.ISOCode = isoCode
				case "type":
					typeStr, typeErr := d.DecodeString()
					if typeErr != nil {
						return typeErr
					}
					c.RepresentedCountry.Type = typeStr
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.RepresentedCountry.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.RepresentedCountry.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "country":
			for countryKey, countryErr := range d.DecodeMap() {
				if countryErr != nil {
					return countryErr
				}
				switch string(countryKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.Country.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					c.Country.ISOCode = isoCode
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.Country.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Country.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "registered_country":
			for regKey, regErr := range d.DecodeMap() {
				if regErr != nil {
					return regErr
				}
				switch string(regKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.RegisteredCountry.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					c.RegisteredCountry.ISOCode = isoCode
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.RegisteredCountry.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.RegisteredCountry.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "location":
			for locKey, locErr := range d.DecodeMap() {
				if locErr != nil {
					return locErr
				}
				switch string(locKey) {
				case "time_zone":
					tz, tzErr := d.DecodeString()
					if tzErr != nil {
						return tzErr
					}
					c.Location.TimeZone = tz
				case "latitude":
					lat, latErr := d.DecodeFloat64()
					if latErr != nil {
						return latErr
					}
					c.Location.Latitude = lat
				case "longitude":
					lon, lonErr := d.DecodeFloat64()
					if lonErr != nil {
						return lonErr
					}
					c.Location.Longitude = lon
				case "metro_code":
					metro, metroErr := d.DecodeUInt16()
					if metroErr != nil {
						return metroErr
					}
					c.Location.MetroCode = uint(metro)
				case "accuracy_radius":
					radius, radiusErr := d.DecodeUInt16()
					if radiusErr != nil {
						return radiusErr
					}
					c.Location.AccuracyRadius = radius
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "traits":
			for traitsKey, traitsErr := range d.DecodeMap() {
				if traitsErr != nil {
					return traitsErr
				}
				switch string(traitsKey) {
				case "is_anonymous_proxy":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Traits.IsAnonymousProxy = val
				case "is_anycast":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Traits.IsAnycast = val
				case "is_satellite_provider":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Traits.IsSatelliteProvider = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The Country struct corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
type Country struct {
	Continent struct {
		Names     map[string]string `maxminddb:"names"`
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
	} `maxminddb:"continent"`
	Country struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		Names             map[string]string `maxminddb:"names"`
		ISOCode           string            `maxminddb:"iso_code"`
		Type              string            `maxminddb:"type"`
		GeoNameID         uint              `maxminddb:"geoname_id"`
		IsInEuropeanUnion bool              `maxminddb:"is_in_european_union"`
	} `maxminddb:"represented_country"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsAnycast           bool `maxminddb:"is_anycast"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for Country.
func (c *Country) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "continent":
			for contKey, contErr := range d.DecodeMap() {
				if contErr != nil {
					return contErr
				}
				switch string(contKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.Continent.Names = names
				case "code":
					code, codeErr := d.DecodeString()
					if codeErr != nil {
						return codeErr
					}
					c.Continent.Code = code
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.Continent.GeoNameID = uint(geoID)
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "country":
			for countryKey, countryErr := range d.DecodeMap() {
				if countryErr != nil {
					return countryErr
				}
				switch string(countryKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.Country.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					c.Country.ISOCode = isoCode
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.Country.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Country.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "registered_country":
			for regKey, regErr := range d.DecodeMap() {
				if regErr != nil {
					return regErr
				}
				switch string(regKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.RegisteredCountry.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					c.RegisteredCountry.ISOCode = isoCode
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.RegisteredCountry.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.RegisteredCountry.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "represented_country":
			for repKey, repErr := range d.DecodeMap() {
				if repErr != nil {
					return repErr
				}
				switch string(repKey) {
				case "names":
					names := make(map[string]string)
					for nameKey, nameErr := range d.DecodeMap() {
						if nameErr != nil {
							return nameErr
						}
						value, valueErr := d.DecodeString()
						if valueErr != nil {
							return valueErr
						}
						names[string(nameKey)] = value
					}
					c.RepresentedCountry.Names = names
				case "iso_code":
					isoCode, isoErr := d.DecodeString()
					if isoErr != nil {
						return isoErr
					}
					c.RepresentedCountry.ISOCode = isoCode
				case "type":
					typeStr, typeErr := d.DecodeString()
					if typeErr != nil {
						return typeErr
					}
					c.RepresentedCountry.Type = typeStr
				case "geoname_id":
					geoID, geoErr := d.DecodeUInt32()
					if geoErr != nil {
						return geoErr
					}
					c.RepresentedCountry.GeoNameID = uint(geoID)
				case "is_in_european_union":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.RepresentedCountry.IsInEuropeanUnion = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		case "traits":
			for traitsKey, traitsErr := range d.DecodeMap() {
				if traitsErr != nil {
					return traitsErr
				}
				switch string(traitsKey) {
				case "is_anonymous_proxy":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Traits.IsAnonymousProxy = val
				case "is_anycast":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Traits.IsAnycast = val
				case "is_satellite_provider":
					val, valErr := d.DecodeBool()
					if valErr != nil {
						return valErr
					}
					c.Traits.IsSatelliteProvider = val
				default:
					if err := d.SkipValue(); err != nil {
						return err
					}
				}
			}
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The AnonymousIP struct corresponds to the data in the GeoIP2
// Anonymous IP database.
type AnonymousIP struct {
	IsAnonymous        bool `maxminddb:"is_anonymous"`
	IsAnonymousVPN     bool `maxminddb:"is_anonymous_vpn"`
	IsHostingProvider  bool `maxminddb:"is_hosting_provider"`
	IsPublicProxy      bool `maxminddb:"is_public_proxy"`
	IsResidentialProxy bool `maxminddb:"is_residential_proxy"`
	IsTorExitNode      bool `maxminddb:"is_tor_exit_node"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for AnonymousIP.
func (a *AnonymousIP) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "is_anonymous":
			val, valErr := d.DecodeBool()
			if valErr != nil {
				return valErr
			}
			a.IsAnonymous = val
		case "is_anonymous_vpn":
			val, valErr := d.DecodeBool()
			if valErr != nil {
				return valErr
			}
			a.IsAnonymousVPN = val
		case "is_hosting_provider":
			val, valErr := d.DecodeBool()
			if valErr != nil {
				return valErr
			}
			a.IsHostingProvider = val
		case "is_public_proxy":
			val, valErr := d.DecodeBool()
			if valErr != nil {
				return valErr
			}
			a.IsPublicProxy = val
		case "is_residential_proxy":
			val, valErr := d.DecodeBool()
			if valErr != nil {
				return valErr
			}
			a.IsResidentialProxy = val
		case "is_tor_exit_node":
			val, valErr := d.DecodeBool()
			if valErr != nil {
				return valErr
			}
			a.IsTorExitNode = val
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The ASN struct corresponds to the data in the GeoLite2 ASN database.
type ASN struct {
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for ASN.
func (a *ASN) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "autonomous_system_organization":
			org, orgErr := d.DecodeString()
			if orgErr != nil {
				return orgErr
			}
			a.AutonomousSystemOrganization = org
		case "autonomous_system_number":
			asn, asnErr := d.DecodeUInt32()
			if asnErr != nil {
				return asnErr
			}
			a.AutonomousSystemNumber = uint(asn)
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The ConnectionType struct corresponds to the data in the GeoIP2
// Connection-Type database.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for ConnectionType.
func (c *ConnectionType) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "connection_type":
			connType, connErr := d.DecodeString()
			if connErr != nil {
				return connErr
			}
			c.ConnectionType = connType
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The Domain struct corresponds to the data in the GeoIP2 Domain database.
type Domain struct {
	Domain string `maxminddb:"domain"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for Domain.
func (d *Domain) UnmarshalMaxMindDB(decoder *maxminddb.Decoder) error {
	for key, err := range decoder.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "domain":
			domain, domainErr := decoder.DecodeString()
			if domainErr != nil {
				return domainErr
			}
			d.Domain = domain
		default:
			if err := decoder.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

// The ISP struct corresponds to the data in the GeoIP2 ISP database.
type ISP struct {
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization"`
	ISP                          string `maxminddb:"isp"`
	MobileCountryCode            string `maxminddb:"mobile_country_code"`
	MobileNetworkCode            string `maxminddb:"mobile_network_code"`
	Organization                 string `maxminddb:"organization"`
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number"`
}

// UnmarshalMaxMindDB implements maxminddb.Unmarshaler for ISP.
func (i *ISP) UnmarshalMaxMindDB(d *maxminddb.Decoder) error {
	for key, err := range d.DecodeMap() {
		if err != nil {
			return err
		}

		switch string(key) {
		case "autonomous_system_organization":
			org, orgErr := d.DecodeString()
			if orgErr != nil {
				return orgErr
			}
			i.AutonomousSystemOrganization = org
		case "isp":
			isp, ispErr := d.DecodeString()
			if ispErr != nil {
				return ispErr
			}
			i.ISP = isp
		case "mobile_country_code":
			mcc, mccErr := d.DecodeString()
			if mccErr != nil {
				return mccErr
			}
			i.MobileCountryCode = mcc
		case "mobile_network_code":
			mnc, mncErr := d.DecodeString()
			if mncErr != nil {
				return mncErr
			}
			i.MobileNetworkCode = mnc
		case "organization":
			org, orgErr := d.DecodeString()
			if orgErr != nil {
				return orgErr
			}
			i.Organization = org
		case "autonomous_system_number":
			asn, asnErr := d.DecodeUInt32()
			if asnErr != nil {
				return asnErr
			}
			i.AutonomousSystemNumber = uint(asn)
		default:
			if err := d.SkipValue(); err != nil {
				return err
			}
		}
	}
	return nil
}

type databaseType int

const (
	isAnonymousIP = 1 << iota
	isASN
	isCity
	isConnectionType
	isCountry
	isDomain
	isEnterprise
	isISP
)

// Reader holds the maxminddb.Reader struct. It can be created using the
// Open and FromBytes functions.
type Reader struct {
	mmdbReader   *maxminddb.Reader
	databaseType databaseType
}

// InvalidMethodError is returned when a lookup method is called on a
// database that it does not support. For instance, calling the ISP method
// on a City database.
type InvalidMethodError struct {
	Method       string
	DatabaseType string
}

func (e InvalidMethodError) Error() string {
	return fmt.Sprintf(`geoip2: the %s method does not support the %s database`,
		e.Method, e.DatabaseType)
}

// UnknownDatabaseTypeError is returned when an unknown database type is
// opened.
type UnknownDatabaseTypeError struct {
	DatabaseType string
}

func (e UnknownDatabaseTypeError) Error() string {
	return fmt.Sprintf(`geoip2: reader does not support the %q database type`,
		e.DatabaseType)
}

// Open takes a string path to a file and returns a Reader struct or an error.
// The database file is opened using a memory map. Use the Close method on the
// Reader object to return the resources to the system.
func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}

// FromBytes takes a byte slice corresponding to a GeoIP2/GeoLite2 database
// file and returns a Reader struct or an error. Note that the byte slice is
// used directly; any modification of it after opening the database will result
// in errors while reading from the database.
func FromBytes(bytes []byte) (*Reader, error) {
	reader, err := maxminddb.FromBytes(bytes)
	if err != nil {
		return nil, err
	}
	dbType, err := getDBType(reader)
	return &Reader{reader, dbType}, err
}

func getDBType(reader *maxminddb.Reader) (databaseType, error) {
	switch reader.Metadata.DatabaseType {
	case "GeoIP2-Anonymous-IP":
		return isAnonymousIP, nil
	case "DBIP-ASN-Lite (compat=GeoLite2-ASN)",
		"GeoLite2-ASN":
		return isASN, nil
	// We allow City lookups on Country for back compat
	case "DBIP-City-Lite",
		"DBIP-Country-Lite",
		"DBIP-Country",
		"DBIP-Location (compat=City)",
		"GeoLite2-City",
		"GeoIP2-City",
		"GeoIP2-City-Africa",
		"GeoIP2-City-Asia-Pacific",
		"GeoIP2-City-Europe",
		"GeoIP2-City-North-America",
		"GeoIP2-City-South-America",
		"GeoIP2-Precision-City",
		"GeoLite2-Country",
		"GeoIP2-Country":
		return isCity | isCountry, nil
	case "GeoIP2-Connection-Type":
		return isConnectionType, nil
	case "GeoIP2-Domain":
		return isDomain, nil
	case "DBIP-ISP (compat=Enterprise)",
		"DBIP-Location-ISP (compat=Enterprise)",
		"GeoIP2-Enterprise":
		return isEnterprise | isCity | isCountry, nil
	case "GeoIP2-ISP", "GeoIP2-Precision-ISP":
		return isISP | isASN, nil
	default:
		return 0, UnknownDatabaseTypeError{reader.Metadata.DatabaseType}
	}
}

// Enterprise takes an IP address as a netip.Addr and returns an Enterprise
// struct and/or an error. This is intended to be used with the GeoIP2
// Enterprise database.
func (r *Reader) Enterprise(ipAddress netip.Addr) (*Enterprise, error) {
	if isEnterprise&r.databaseType == 0 {
		return nil, InvalidMethodError{"Enterprise", r.Metadata().DatabaseType}
	}
	var enterprise Enterprise
	err := r.mmdbReader.Lookup(ipAddress).Decode(&enterprise)
	return &enterprise, err
}

// City takes an IP address as a netip.Addr and returns a City struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 City databases.
func (r *Reader) City(ipAddress netip.Addr) (*City, error) {
	if isCity&r.databaseType == 0 {
		return nil, InvalidMethodError{"City", r.Metadata().DatabaseType}
	}
	var city City
	err := r.mmdbReader.Lookup(ipAddress).Decode(&city)
	return &city, err
}

// Country takes an IP address as a netip.Addr and returns a Country struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 Country
// databases.
func (r *Reader) Country(ipAddress netip.Addr) (*Country, error) {
	if isCountry&r.databaseType == 0 {
		return nil, InvalidMethodError{"Country", r.Metadata().DatabaseType}
	}
	var country Country
	err := r.mmdbReader.Lookup(ipAddress).Decode(&country)
	return &country, err
}

// AnonymousIP takes an IP address as a netip.Addr and returns a
// AnonymousIP struct and/or an error.
func (r *Reader) AnonymousIP(ipAddress netip.Addr) (*AnonymousIP, error) {
	if isAnonymousIP&r.databaseType == 0 {
		return nil, InvalidMethodError{"AnonymousIP", r.Metadata().DatabaseType}
	}
	var anonIP AnonymousIP
	err := r.mmdbReader.Lookup(ipAddress).Decode(&anonIP)
	return &anonIP, err
}

// ASN takes an IP address as a netip.Addr and returns a ASN struct and/or
// an error.
func (r *Reader) ASN(ipAddress netip.Addr) (*ASN, error) {
	if isASN&r.databaseType == 0 {
		return nil, InvalidMethodError{"ASN", r.Metadata().DatabaseType}
	}
	var val ASN
	err := r.mmdbReader.Lookup(ipAddress).Decode(&val)
	return &val, err
}

// ConnectionType takes an IP address as a netip.Addr and returns a
// ConnectionType struct and/or an error.
func (r *Reader) ConnectionType(ipAddress netip.Addr) (*ConnectionType, error) {
	if isConnectionType&r.databaseType == 0 {
		return nil, InvalidMethodError{"ConnectionType", r.Metadata().DatabaseType}
	}
	var val ConnectionType
	err := r.mmdbReader.Lookup(ipAddress).Decode(&val)
	return &val, err
}

// Domain takes an IP address as a netip.Addr and returns a
// Domain struct and/or an error.
func (r *Reader) Domain(ipAddress netip.Addr) (*Domain, error) {
	if isDomain&r.databaseType == 0 {
		return nil, InvalidMethodError{"Domain", r.Metadata().DatabaseType}
	}
	var val Domain
	err := r.mmdbReader.Lookup(ipAddress).Decode(&val)
	return &val, err
}

// ISP takes an IP address as a netip.Addr and returns a ISP struct and/or
// an error.
func (r *Reader) ISP(ipAddress netip.Addr) (*ISP, error) {
	if isISP&r.databaseType == 0 {
		return nil, InvalidMethodError{"ISP", r.Metadata().DatabaseType}
	}
	var val ISP
	err := r.mmdbReader.Lookup(ipAddress).Decode(&val)
	return &val, err
}

// Metadata takes no arguments and returns a struct containing metadata about
// the MaxMind database in use by the Reader.
func (r *Reader) Metadata() maxminddb.Metadata {
	return r.mmdbReader.Metadata
}

// Close unmaps the database file from virtual memory and returns the
// resources to the system.
func (r *Reader) Close() error {
	return r.mmdbReader.Close()
}
