# GeoIP2 Reader for Go #

[![Build Status](https://travis-ci.org/oschwald/geoip2-golang.png?branch=master)](https://travis-ci.org/oschwald/geoip2-golang)
[![GoDoc](https://godoc.org/github.com/oschwald/geoip2-golang?status.png)](https://godoc.org/github.com/oschwald/geoip2-golang)

This library reads MaxMind [GeoLite2](http://dev.maxmind.com/geoip/geoip2/geolite2/)
and [GeoIP2](http://www.maxmind.com/en/geolocation_landing) databases.

This is not an official MaxMind API.

## Installation ##

```
go get github.com/oschwald/geoip2-golang
```

## Usage ##

[See GoDoc](http://godoc.org/github.com/oschwald/geoip2-golang) for
documentation and examples.

## Example ##

```go
package main

import (
    "fmt"
    "github.com/oschwald/geoip2-golang"
    "net"
)

func main() {
    db, err := geoip2.Open("GeoIP2-City.mmdb")
    if err != nil {
            panic(err)
    }
    // If you are using strings that may be invalid, check that ip is not nil
    ip := net.ParseIP("81.2.69.142")
    record, err := db.City(ip)
    if err != nil {
            panic(err)
    }
    fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
    fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
    fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
    fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)
    fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
    fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)

    db.Close()
    // Output:
    // Portuguese (BR) city name: Londres
    // English subdivision name: England
    // Russian country name: Великобритания
    // ISO country code: GB
    // Time zone: Europe/London
    // Coordinates: 51.5142, -0.0931
}
```

## Notes ##

This uses the [maxminddb library](https://github.com/oschwald/maxminddb-golang).

## Contributing ##

Contributions welcome! Please fork the repository and open a pull request
with your changes.

## License ##

This is free software, licensed under the Apache License, Version 2.0.
