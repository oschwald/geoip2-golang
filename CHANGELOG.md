# 2.0.0-beta.1 - 2025-01-XX

* **BREAKING CHANGE**: Updated to use `maxminddb-golang/v2` which provides
  significant performance improvements and a more modern API.
* **BREAKING CHANGE**: All lookup methods now accept `netip.Addr` instead of
  `net.IP`. This provides better performance and aligns with modern Go
  networking practices.
* Updated module path to `github.com/oschwald/geoip2-golang/v2` to follow
  Go's semantic versioning guidelines for breaking changes.
* Updated examples and documentation to demonstrate proper error handling
  with `netip.ParseAddr()`.
* Updated linting rules to support both v1 and v2 import paths during the
  transition period.

## Migration Guide

To migrate from v1 to v2:

1. Update your import path:
   ```go
   // Old
   import "github.com/oschwald/geoip2-golang"

   // New
   import "github.com/oschwald/geoip2-golang/v2"
   ```

2. Replace `net.IP` with `netip.Addr`:
   ```go
   // Old
   ip := net.ParseIP("81.2.69.142")
   record, err := db.City(ip)

   // New
   ip, err := netip.ParseAddr("81.2.69.142")
   if err != nil {
       // handle error
   }
   record, err := db.City(ip)
   ```

# 1.11.0 - 2024-06-03

* Go 1.21 or greater is now required.
* The new `is_anycast` output is now supported on the GeoIP2 Country, City,
  and Enterprise databases. [#119](https://github.com/oschwald/geoip2-golang/issues/119).

Note: 1.10.0 was accidentally skipped.

# 1.9.0 - 2023-06-18

* Rearrange fields in structs to reduce memory usage. Although this
  does reduce readability, these structs are often created at very
  rates, making the trade-off worth it.

# 1.8.0 - 2022-08-07

* Set Go version to 1.18 in go.mod.

# 1.7.0 - 2022-03-26

* Set the minimum Go version in the go.mod file to 1.17.
* Updated dependencies.

# 1.6.1 - 2022-01-28

* This is a re-release with the changes that were supposed to be in 1.6.0.

# 1.6.0 - 2022-01-28

* Add support for new `mobile_country_code` and `mobile_network_code` outputs
  on GeoIP2 ISP and GeoIP2 Enterprise.

# 1.5.0 - 2021-02-20

* Add `StaticIPScore` field to Enterprise. Pull request by Pierre
  Bonzel. GitHub [#54](https://github.com/oschwald/geoip2-golang/issues/54).
* Add `IsResidentialProxy` field to `AnonymousIP`. Pull request by
  Brendan Boyle. GitHub [#72](https://github.com/oschwald/geoip2-golang/issues/72).
* Support DBIP-ASN-Lite database. Requested by Muhammad Hussein
  Fattahizadeh. GitHub [#69](https://github.com/oschwald/geoip2-golang/issues/69).

# 1.4.0 - 2019-12-25

* This module now uses Go modules. Requested by Axel Etcheverry.
  GitHub [#52](https://github.com/oschwald/geoip2-golang/issues/52).
* DBIP databases are now supported. Requested by jaw0. GitHub [#45](https://github.com/oschwald/geoip2-golang/issues/45).
* Allow using the ASN method with the GeoIP2 ISP database. Pull request
  by lspgn. GitHub [#47](https://github.com/oschwald/geoip2-golang/issues/47).
* The example in the `README.md` now checks the length of the
  subdivision slice before using it. GitHub [#51](https://github.com/oschwald/geoip2-golang/issues/51).

# 1.3.0 - 2019-08-28

* Added support for the GeoIP2 Enterprise database.

# 1.2.1 - 2018-02-25

* HTTPS is now used for the test data submodule rather than the Git
  protocol

# 1.2.0 - 2018-02-19

* The country structs for `geoip2.City` and `geoip2.Country` now have an
  `IsInEuropeanUnion` boolean field. This is true when the associated
  country is a member state of the European Union. This requires a
  database built on or after February 13, 2018.
* Switch from Go Check to Testify. Closes [#27](https://github.com/oschwald/geoip2-golang/issues/27)

# 1.1.0 - 2017-04-23

* Add support for the GeoLite2 ASN database.
* Add support for the GeoIP2 City by Continent databases. GitHub [#26](https://github.com/oschwald/geoip2-golang/issues/26).


# 1.0.0 - 2016-11-09

New release for those using tagged releases. Closes [#21](https://github.com/oschwald/geoip2-golang/issues/21).
