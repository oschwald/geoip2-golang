# Migrating to geoip2-golang v2

Version 2.0 modernizes the public API, improves performance, and tightens
parity with current MaxMind databases. This guide highlights the most important
changes and shows how to update existing v1 code.

## Upgrade Checklist

1. **Update imports**

   ```go
   // v1
   import "github.com/oschwald/geoip2-golang"

   // v2
   import "github.com/oschwald/geoip2-golang/v2"
   ```

2. **Switch to `netip.Addr`**

   Lookup methods now take `netip.Addr` instead of `net.IP`. Use
   `netip.ParseAddr` (errors on invalid input) or `netip.MustParseAddr` for
   trusted literals. This avoids ambiguity around IPv4-in-IPv6 addresses and
   removes heap allocations for most lookups.

   ```go
   // v1
   ip := net.ParseIP("81.2.69.142")
   record, err := db.City(ip)

   // v2
   ip, err := netip.ParseAddr("81.2.69.142")
   if err != nil {
       return err
   }
   record, err := db.City(ip)
   ```

   To convert an existing `net.IP`, use `netip.AddrFromSlice`, check its `ok`
   result, and call `Unmap()`. `net.ParseIP` returns IPv4 addresses in 16-byte
   form, so without `Unmap()` they stay IPv4-mapped IPv6 addresses: lookups in
   IPv4-only databases fail, and `IPAddress` and `Network` come back in mapped
   form from IPv6 databases.

   ```go
   addr, ok := netip.AddrFromSlice(ip) // ip is a net.IP
   if !ok {
       return fmt.Errorf("invalid IP address %v", ip)
   }
   record, err := db.City(addr.Unmap())
   ```

3. **Adopt the new `HasData()` helpers**

   All result structs have a `HasData()` method that ignores the `Network` and
   `IPAddress` fields, which every successful lookup fills. Use it instead of
   comparing against zero values.

   ```go
   record, err := db.City(ip)
   if err != nil {
       return err
   }
   if !record.HasData() {
       // handle missing GeoIP data
   }
   ```

4. **Use structured `Names` fields**

   Localized names are now strongly typed instead of map lookups. Access the
   language-specific struct field directly.

   ```go
   // v1
   cityEn := record.City.Names["en"]
   cityPtBr := record.City.Names["pt-BR"]

   // v2
   cityEn := record.City.Names.English
   cityPtBr := record.City.Names.BrazilianPortuguese
   ```

   Supported fields include `English`, `German`, `Spanish`, `French`,
   `Japanese`, `BrazilianPortuguese`, `Russian`, and `SimplifiedChinese`. Names
   in any other language are not decoded. If you need them, use maxminddb's
   `Lookup` with your own struct.

5. **Expect `Network` and `IPAddress` on every result**

   Each successful lookup populates the queried IP address and the enclosing
   network on every result struct, even if `HasData()` returns false. For
   `City`, `Country`, and `Enterprise` they are under `Traits`, such as
   `record.Traits.Network`. If the lookup returns an error, the record may be
   nil, and any record that is returned has them empty. Because these fields
   are set, a record with no data is no longer equal to its zero value. Code
   that compared a record against its zero value to detect missing data should
   use `HasData()` instead.

6. **Adjust renamed fields**

   `IsoCode` is now `ISOCode` across all structs. Update code accordingly.

   ```go
   // v1
   code := record.Country.IsoCode

   // v2
   code := record.Country.ISOCode
   ```

7. **Check for coordinates before using them**

   `Location.Latitude` and `Location.Longitude` changed from `float64` to
   `*float64`. A missing coordinate is now `nil` instead of `0`, so a real
   location at (0, 0) is no longer mistaken for missing data. Code that does
   arithmetic with them or assigns them to a `float64` no longer compiles, but
   other uses still compile and silently change behavior. Printing or logging
   them shows a pointer address instead of the coordinate, and `==` or `!=` on
   a coordinate or on a whole `Location` compares pointers, so two lookups of
   the same place with coordinates are never equal. Review every read of these
   fields, not only the ones the compiler flags. Dereferencing them without a
   check panics when the database has no coordinates for the address, so use
   `HasCoordinates()`, which returns true only when both are set. Missing
   coordinates are also omitted from JSON output.

   ```go
   // v1
   var lat, lon float64 = record.Location.Latitude, record.Location.Longitude

   // v2
   if record.Location.HasCoordinates() {
       lat := *record.Location.Latitude
       lon := *record.Location.Longitude
       fmt.Println(lat, lon)
   }
   ```

8. **Replace removed fields**

   `IsAnonymousProxy` and `IsSatelliteProvider` were removed from the `Traits`
   structs of `City`, `Country`, and `Enterprise` because MaxMind no longer
   includes them in its databases. For anonymity detection, look the address up
   in a GeoIP Anonymous IP database (`AnonymousIP`) or GeoIP Anonymous Plus
   database (`AnonymousPlus`) instead, and check the specific flags you need.
   Their `IsAnonymous` flag is broader than v1's `IsAnonymousProxy`: it is also
   set for hosting providers, residential proxies, and Tor exit nodes. In place
   of `IsSatelliteProvider`, check for a `ConnectionType` of `"Satellite"` from
   a GeoIP Connection Type database or from `Enterprise.Traits.ConnectionType`.

   `Enterprise.RegisteredCountry` now uses the shared `CountryRecord` type,
   which has no `Confidence` field. `Enterprise.Country.Confidence` is
   unchanged.

   Three fields that v1 had are still present but deprecated, so linters flag
   them: `Location.MetroCode`, which MaxMind no longer maintains, and
   `EnterpriseTraits.IsLegitimateProxy` and `EnterpriseTraits.StaticIPScore`,
   which will be removed in the next major release. `StaticIPScore` has never
   been populated.

   ```go
   // v1
   if record.Traits.IsAnonymousProxy {
       // ...
   }

   // v2, with anonDB opened from a GeoIP Anonymous IP database
   anon, err := anonDB.AnonymousIP(ip)
   if err != nil {
       return err
   }
   if anon.IsPublicProxy || anon.IsAnonymousVPN {
       // ...
   }

   // v1
   if record.Traits.IsSatelliteProvider {
       // ...
   }

   // v2, with ctDB opened from a GeoIP Connection Type database
   ct, err := ctDB.ConnectionType(ip)
   if err != nil {
       return err
   }
   if ct.ConnectionType == "Satellite" {
       // ...
   }
   ```

9. **Use the named nested types**

   The nested structs in `City`, `Country`, and `Enterprise` were anonymous in
   v1 and are now named types. Field paths such as `record.City.GeoNameID` are
   unchanged, but code that spelled out an anonymous struct type, for example
   to build a record in a test or to pass a subdivision to a helper, must use
   the named type instead.

   | Field                | `City`               | `Country`            | `Enterprise`              |
   | -------------------- | -------------------- | -------------------- | ------------------------- |
   | `City`               | `CityRecord`         |                      | `EnterpriseCityRecord`    |
   | `Continent`          | `Continent`          | `Continent`          | `Continent`               |
   | `Country`            | `CountryRecord`      | `CountryRecord`      | `EnterpriseCountryRecord` |
   | `Location`           | `Location`           |                      | `Location`                |
   | `Postal`             | `CityPostal`         |                      | `EnterprisePostal`        |
   | `RegisteredCountry`  | `CountryRecord`      | `CountryRecord`      | `CountryRecord`           |
   | `RepresentedCountry` | `RepresentedCountry` | `RepresentedCountry` | `RepresentedCountry`      |
   | `Subdivisions`       | `[]CitySubdivision`  |                      | `[]EnterpriseSubdivision` |
   | `Traits`             | `CityTraits`         | `CountryTraits`      | `EnterpriseTraits`        |

   ```go
   // v2
   func subdivisionName(s geoip2.CitySubdivision) string {
       return s.Names.English
   }

   func testRecord() *geoip2.City {
       return &geoip2.City{
           City:    geoip2.CityRecord{Names: geoip2.Names{English: "London"}},
           Country: geoip2.CountryRecord{ISOCode: "GB"},
       }
   }
   ```

10. **Update JSON consumers**

    v1 structs had no JSON tags, so `json.Marshal` used the Go field names,
    such as `Country` and `IsoCode`. v2 uses the database's names, such as
    `country` and `iso_code`, adds `ip_address` and `network` (under `traits`
    for `City`, `Country`, and `Enterprise`), and omits empty values. Anything
    that reads JSON you produce from these structs needs updating.

    Regenerate JSON that v1 wrote, for example in a cache or database, rather
    than decoding it into v2 types. `encoding/json` ignores case, so some keys
    such as `Country` and `Latitude` still match, but others such as `IsoCode`
    and `GeoNameID` do not and are silently left empty. v1 also wrote `0` for
    missing coordinates, which v2 decodes as a real location at (0, 0).

    ```go
    // v1 output includes: "Country":{...,"IsoCode":"GB",...}
    // v2 output includes: "country":{...,"iso_code":"GB",...}
    data, err := json.Marshal(record)
    ```

## Additional Notes

- v2 needs a newer Go toolchain than v1. Version 2.0.0 required Go 1.24 for
  `omitzero` JSON tags, and later releases require the version in `go.mod`.
- `Open` and `OpenBytes` accept optional `Option` values, allowing future
  extensions without new breaking releases.
- The legacy `FromBytes` helper is removed. Use `OpenBytes` instead.
- In releases after 2.5.0, `Open` and `OpenBytes` return a nil reader when the
  database type is unsupported. v1 and earlier v2 releases returned a reader
  along with the error. The error may be joined with a close error, so check
  for it with `errors.As` rather than a type assertion.
- Since 2.4.0, lookups enforce size limits on strings and subdivision arrays
  and reject duplicate fields, returning an `InvalidDatabaseError`. Custom or
  third-party databases that v1 read may fail these checks.
- `Reader.Metadata()` returns the `Metadata` type from
  `github.com/oschwald/maxminddb-golang/v2`. Code that names that type must
  import the v2 package.

Following the checklist above should cover most migrations. For edge cases,
review `models.go` for field definitions and `example_test.go` for idiomatic
usage with the v2 API.
