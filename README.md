# GeoIP2 Reader for Go #

This library reads MaxMind [GeoLite2](http://dev.maxmind.com/geoip/geoip2/geolite2/) and
[GeoIP2](http://www.maxmind.com/en/geolocation_landing) databases.

# Example #

```go
package main

import (
    "fmt"
    "github.com/oschwald/geoip2-golang"
    "log"
    "net"
)

func main() {
    db, err := geoip2.Open("GeoLite2-City.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    ip := net.ParseIP("1.1.1.1")
    record, err := db.City(ip)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(record)
    db.Close()
}
```

## Contributing ##

Contributions welcome! Please fork the repository and open a pull request
with your changes.

## License ##

This is free software, licensed under the Apache License, Version 2.0.
