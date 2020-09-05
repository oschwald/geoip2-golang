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
	"net"

	"github.com/oschwald/maxminddb-golang"
)

// Names are the localized names for the location.
type Names struct {
	De   string `maxminddb:"de" json:"de,omitempty"`
	En   string `maxminddb:"en" json:"en,omitempty"`
	Es   string `maxminddb:"es" json:"es,omitempty"`
	Fr   string `maxminddb:"fr" json:"fr,omitempty"`
	Ja   string `maxminddb:"ja" json:"ja,omitempty"`
	PtBR string `maxminddb:"pt-BR" json:"pt-BR,omitempty"`
	Ru   string `maxminddb:"ru" json:"ru,omitempty"`
	ZhCN string `maxminddb:"zh-CN" json:"zh-CN,omitempty"`
}

// The Enterprise struct corresponds to the data in the GeoIP2 Enterprise
// database.
type Enterprise struct {
	City struct {
		Confidence uint8 `maxminddb:"confidence" json:"confidence,omitempty"`
		GeoNameID  uint  `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		Names      Names `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"city" json:"city,omitempty"`
	Continent struct {
		Code      string `maxminddb:"code" json:"code,omitempty"`
		GeoNameID uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		Names     Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"continent" json:"continent,omitempty"`
	Country struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
		Confidence        uint8  `maxminddb:"confidence" json:"confidence,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
	} `maxminddb:"country" json:"country,omitempty"`
	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius" json:"accuracy_radius,omitempty"`
		Latitude       float64 `maxminddb:"latitude" json:"latitude,omitempty"`
		Longitude      float64 `maxminddb:"longitude" json:"longitude,omitempty"`
		MetroCode      uint    `maxminddb:"metro_code" json:"metro_code,omitempty"`
		TimeZone       string  `maxminddb:"time_zone" json:"time_zone,omitempty"`
	} `maxminddb:"location" json:"location,omitempty"`
	Postal struct {
		Code       string `maxminddb:"code" json:"code,omitempty"`
		Confidence uint8  `maxminddb:"confidence" json:"confidence,omitempty"`
	} `maxminddb:"postal" json:"postal,omitempty"`
	RegisteredCountry struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
		Confidence        uint8  `maxminddb:"confidence" json:"confidence,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
	} `maxminddb:"registered_country" json:"registered_country,omitempty"`
	RepresentedCountry struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
		Type              string `maxminddb:"type" json:"type,omitempty"`
	} `maxminddb:"represented_country" json:"represented_country,omitempty"`
	Subdivisions []struct {
		Confidence uint8  `maxminddb:"confidence" json:"confidence,omitempty"`
		GeoNameID  uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsoCode    string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names      Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"subdivisions" json:"subdivisions,omitempty"`
	Traits struct {
		AutonomousSystemNumber       uint    `maxminddb:"autonomous_system_number" json:"autonomous_system_number,omitempty"`
		AutonomousSystemOrganization string  `maxminddb:"autonomous_system_organization" json:"autonomous_system_organization,omitempty"`
		ConnectionType               string  `maxminddb:"connection_type" json:"connection_type,omitempty"`
		Domain                       string  `maxminddb:"domain" json:"domain,omitempty"`
		IsAnonymousProxy             bool    `maxminddb:"is_anonymous_proxy" json:"is_anonymous_proxy,omitempty"`
		IsLegitimateProxy            bool    `maxminddb:"is_legitimate_proxy" json:"is_legitimate_proxy,omitempty"`
		IsSatelliteProvider          bool    `maxminddb:"is_satellite_provider" json:"is_satellite_provider,omitempty"`
		ISP                          string  `maxminddb:"isp" json:"isp,omitempty"`
		StaticIPScore                float64 `maxminddb:"static_ip_score" json:"static_ip_score,omitempty"`
		Organization                 string  `maxminddb:"organization" json:"organization,omitempty"`
		UserType                     string  `maxminddb:"user_type" json:"user_type,omitempty"`
	} `maxminddb:"traits" json:"traits,omitempty"`
}

// The City struct corresponds to the data in the GeoIP2/GeoLite2 City
// databases.
type City struct {
	City struct {
		GeoNameID uint  `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		Names     Names `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"city" json:"city,omitempty"`
	Continent struct {
		Code      string `maxminddb:"code" json:"code,omitempty"`
		GeoNameID uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		Names     Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"continent" json:"continent,omitempty"`
	Country struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"country" json:"country,omitempty"`
	Location struct {
		AccuracyRadius uint16  `maxminddb:"accuracy_radius" json:"accuracy_radius,omitempty"`
		Latitude       float64 `maxminddb:"latitude" json:"latitude,omitempty"`
		Longitude      float64 `maxminddb:"longitude" json:"longitude,omitempty"`
		MetroCode      uint    `maxminddb:"metro_code" json:"metro_code,omitempty"`
		TimeZone       string  `maxminddb:"time_zone" json:"time_zone,omitempty"`
	} `maxminddb:"location" json:"location,omitempty"`
	Postal struct {
		Code string `maxminddb:"code" json:"code,omitempty"`
	} `maxminddb:"postal" json:"postal,omitempty"`
	RegisteredCountry struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"registered_country" json:"registered_country,omitempty"`
	RepresentedCountry struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
		Type              string `maxminddb:"type" json:"type,omitempty"`
	} `maxminddb:"represented_country" json:"represented_country,omitempty"`
	Subdivisions []struct {
		GeoNameID uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsoCode   string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names     Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"subdivisions" json:"subdivisions,omitempty"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy" json:"is_anonymous_proxy,omitempty"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider" json:"is_satellite_provider,omitempty"`
	} `maxminddb:"traits" json:"traits,omitempty"`
}

// The Country struct corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
type Country struct {
	Continent struct {
		Code      string `maxminddb:"code" json:"code,omitempty"`
		GeoNameID uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		Names     Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"continent" json:"continent,omitempty"`
	Country struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"country" json:"country,omitempty"`
	RegisteredCountry struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
	} `maxminddb:"registered_country" json:"registered_country,omitempty"`
	RepresentedCountry struct {
		GeoNameID         uint   `maxminddb:"geoname_id" json:"geoname_id,omitempty"`
		IsInEuropeanUnion bool   `maxminddb:"is_in_european_union" json:"is_in_european_union,omitempty"`
		IsoCode           string `maxminddb:"iso_code" json:"iso_code,omitempty"`
		Names             Names  `maxminddb:"names" json:"names,omitempty"`
		Type              string `maxminddb:"type" json:"type,omitempty"`
	} `maxminddb:"represented_country" json:"represented_country,omitempty"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy" json:"is_anonymous_proxy,omitempty"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider" json:"is_satellite_provider,omitempty"`
	} `maxminddb:"traits" json:"traits,omitempty"`
}

// The AnonymousIP struct corresponds to the data in the GeoIP2
// Anonymous IP database.
type AnonymousIP struct {
	IsAnonymous       bool `maxminddb:"is_anonymous" json:"is_anonymous,omitempty"`
	IsAnonymousVPN    bool `maxminddb:"is_anonymous_vpn" json:"is_anonymous_vpn,omitempty"`
	IsHostingProvider bool `maxminddb:"is_hosting_provider" json:"is_hosting_provider,omitempty"`
	IsPublicProxy     bool `maxminddb:"is_public_proxy" json:"is_public_proxy,omitempty"`
	IsTorExitNode     bool `maxminddb:"is_tor_exit_node" json:"is_tor_exit_node,omitempty"`
}

// The ASN struct corresponds to the data in the GeoLite2 ASN database.
type ASN struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number" json:"autonomous_system_number,omitempty"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization" json:"autonomous_system_organization,omitempty"`
}

// The ConnectionType struct corresponds to the data in the GeoIP2
// Connection-Type database.
type ConnectionType struct {
	ConnectionType string `maxminddb:"connection_type" json:"connection_type,omitempty"`
}

// The Domain struct corresponds to the data in the GeoIP2 Domain database.
type Domain struct {
	Domain string `maxminddb:"domain" json:"domain,omitempty"`
}

// The ISP struct corresponds to the data in the GeoIP2 ISP database.
type ISP struct {
	AutonomousSystemNumber       uint   `maxminddb:"autonomous_system_number" json:"autonomous_system_number,omitempty"`
	AutonomousSystemOrganization string `maxminddb:"autonomous_system_organization" json:"autonomous_system_organization,omitempty"`
	ISP                          string `maxminddb:"isp" json:"isp,omitempty"`
	Organization                 string `maxminddb:"organization" json:"organization,omitempty"`
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
	return fmt.Sprintf(`geoip2: reader does not support the "%s" database type`,
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
// use directly; any modification of it after opening the database will result
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
	case "GeoLite2-ASN":
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
	case "GeoIP2-ISP",
		"GeoIP2-Precision-ISP":
		return isISP | isASN, nil
	default:
		return 0, UnknownDatabaseTypeError{reader.Metadata.DatabaseType}
	}
}

// Enterprise takes an IP address as a net.IP struct and returns an Enterprise
// struct and/or an error. This is intended to be used with the GeoIP2
// Enterprise database.
func (r *Reader) Enterprise(ipAddress net.IP) (*Enterprise, error) {
	if isEnterprise&r.databaseType == 0 {
		return nil, InvalidMethodError{"Enterprise", r.Metadata().DatabaseType}
	}
	var enterprise Enterprise
	err := r.mmdbReader.Lookup(ipAddress, &enterprise)
	return &enterprise, err
}

// City takes an IP address as a net.IP struct and returns a City struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 City databases.
func (r *Reader) City(ipAddress net.IP) (*City, error) {
	if isCity&r.databaseType == 0 {
		return nil, InvalidMethodError{"City", r.Metadata().DatabaseType}
	}
	var city City
	err := r.mmdbReader.Lookup(ipAddress, &city)
	return &city, err
}

// Country takes an IP address as a net.IP struct and returns a Country struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 Country
// databases.
func (r *Reader) Country(ipAddress net.IP) (*Country, error) {
	if isCountry&r.databaseType == 0 {
		return nil, InvalidMethodError{"Country", r.Metadata().DatabaseType}
	}
	var country Country
	err := r.mmdbReader.Lookup(ipAddress, &country)
	return &country, err
}

// AnonymousIP takes an IP address as a net.IP struct and returns a
// AnonymousIP struct and/or an error.
func (r *Reader) AnonymousIP(ipAddress net.IP) (*AnonymousIP, error) {
	if isAnonymousIP&r.databaseType == 0 {
		return nil, InvalidMethodError{"AnonymousIP", r.Metadata().DatabaseType}
	}
	var anonIP AnonymousIP
	err := r.mmdbReader.Lookup(ipAddress, &anonIP)
	return &anonIP, err
}

// ASN takes an IP address as a net.IP struct and returns a ASN struct and/or
// an error
func (r *Reader) ASN(ipAddress net.IP) (*ASN, error) {
	if isASN&r.databaseType == 0 {
		return nil, InvalidMethodError{"ASN", r.Metadata().DatabaseType}
	}
	var val ASN
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// ConnectionType takes an IP address as a net.IP struct and returns a
// ConnectionType struct and/or an error
func (r *Reader) ConnectionType(ipAddress net.IP) (*ConnectionType, error) {
	if isConnectionType&r.databaseType == 0 {
		return nil, InvalidMethodError{"ConnectionType", r.Metadata().DatabaseType}
	}
	var val ConnectionType
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// Domain takes an IP address as a net.IP struct and returns a
// Domain struct and/or an error
func (r *Reader) Domain(ipAddress net.IP) (*Domain, error) {
	if isDomain&r.databaseType == 0 {
		return nil, InvalidMethodError{"Domain", r.Metadata().DatabaseType}
	}
	var val Domain
	err := r.mmdbReader.Lookup(ipAddress, &val)
	return &val, err
}

// ISP takes an IP address as a net.IP struct and returns a ISP struct and/or
// an error
func (r *Reader) ISP(ipAddress net.IP) (*ISP, error) {
	if isISP&r.databaseType == 0 {
		return nil, InvalidMethodError{"ISP", r.Metadata().DatabaseType}
	}
	var val ISP
	err := r.mmdbReader.Lookup(ipAddress, &val)
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
