package geoip2

import (
	"github.com/oschwald/maxminddb-golang"
	"net"
)

type City struct {
	City struct {
		GeoNameID int `maxminddb:"geoname_id"`
		Names     map[string]string
	}
	Continent struct {
		Code      string
		GeoNameID int `maxminddb:"geoname_id"`
		Names     map[string]string
	}
	Country struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
	}
	Location struct {
		Latitude  float64
		Longitude float64
		MetroCode int    `maxminddb:"metro_code"`
		TimeZone  string `maxminddb:"time_zone"`
	}
	Postal struct {
		Code string
	}
	RegisteredCountry struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
		Type      string
	} `maxminddb:"represented_country"`
	Subdivisions []struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
	}
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	}
}

type Country struct {
	Continent struct {
		Code      string
		GeoNameID int `maxminddb:"geoname_id"`
		Names     map[string]string
	}
	Country struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
	}
	RegisteredCountry struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
	} `maxminddb:"registered_country"`
	RepresentedCountry struct {
		GeoNameID int    `maxminddb:"geoname_id"`
		IsoCode   string `maxminddb:"iso_code"`
		Names     map[string]string
		Type      string
	} `maxminddb:"represented_country"`
	Traits struct {
		IsAnonymousProxy    bool `maxminddb:"is_anonymous_proxy"`
		IsSatelliteProvider bool `maxminddb:"is_satellite_provider"`
	}
}

type Reader struct {
	mmdbReader *maxminddb.Reader
}

func Open(file string) (*Reader, error) {
	reader, err := maxminddb.Open(file)
	return &Reader{mmdbReader: reader}, err
}

func (r *Reader) City(ipAddress net.IP) (*City, error) {
	var city City
	err := r.mmdbReader.Unmarshal(ipAddress, &city)
	return &city, err
}

func (r *Reader) Country(ipAddress net.IP) (*Country, error) {
	var country Country
	err := r.mmdbReader.Unmarshal(ipAddress, &country)
	return &country, err
}

func (r *Reader) Close() {
	r.mmdbReader.Close()
}
