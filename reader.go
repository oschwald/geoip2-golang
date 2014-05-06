// Package geoip2 provides a wrapper around the maxminddb package for
// easy use with the MaxMind GeoIP2 and GeoLite2 databases. The records for
// the IP address is returned from this package as well-formed structures
// that match the internal layout of data from MaxMind.
package geoip2

import (
	"github.com/oschwald/maxminddb-golang"
	"net"
)

// The City structure corresponds to the data in the GeoIP2/GeoLite2 City
// databases.
type City struct {
	City struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Continent struct {
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
		MetroCode uint    `maxminddb:"metro_code"`
		TimeZone  string  `maxminddb:"time_zone"`
	} `maxminddb:"location"`
	Postal struct {
		Code string `maxminddb:"code"`
	} `maxminddb:"postal"`
	RegisteredCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
		Type      string            `maxminddb:"type"`
	} `maxminddb:"represented_country"`
	Subdivisions []struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"subdivisions"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// The Country structure corresponds to the data in the GeoIP2/GeoLite2
// Country databases.
type Country struct {
	Continent struct {
		Code      string            `maxminddb:"code"`
		GeoNameID uint              `maxminddb:"geoname_id"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"continent"`
	Country struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"country"`
	RegisteredCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID uint              `maxminddb:"geoname_id"`
		IsoCode   string            `maxminddb:"iso_code"`
		Names     map[string]string `maxminddb:"names"`
		Type      string            `maxminddb:"type"`
	} `maxminddb:"represented_country"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	} `maxminddb:"traits"`
}

// Reader holds the maxminddb.Reader structure. It should be created
// using the Open function.
type Reader struct {
	mmdbReader *maxminddb.Reader
}

// Open takes a string path to a file and returns a Reader structure or an
// error. The database file is opened using a memory map. Use the Close method
// on the Reader object to return the resources to the system.
func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	return &Reader{mmdbReader: reader}, err
}

// FromBytes takes a byte slice corresponding to a GeoIP2/GeoLite2 database
// file and returns a Reader structure or an error.
func FromBytes(bytes []byte) (*Reader, error) {
	reader, err := maxminddb.FromBytes(bytes)
	return &Reader{mmdbReader: reader}, err
}

// City takes an IP address as a net.IP struct and returns a City struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 City databases.
func (r *Reader) City(ipAddress net.IP) (*City, error) {
	var city City
	err := r.mmdbReader.Lookup(ipAddress, &city)
	return &city, err
}

// Country takes an IP address as a net.IP struct and returns a Country struct
// and/or an error. Although this can be used with other databases, this
// method generally should be used with the GeoIP2 or GeoLite2 Country
// databases.
func (r *Reader) Country(ipAddress net.IP) (*Country, error) {
	var country Country
	err := r.mmdbReader.Lookup(ipAddress, &country)
	return &country, err
}

// Close unmaps the database file from virtual memory and returns the
// resources to the system.
func (r *Reader) Close() {
	r.mmdbReader.Close()
}
