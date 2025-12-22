package geoip2

import (
	"fmt"
	"log"
	"net/netip"
)

// Example provides a basic example of using the API. Use of the Country
// method is analogous to that of the City method.
func Example() {
	db, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, use netip.ParseAddr and check for errors
	ip, err := netip.ParseAddr("81.2.69.142")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.City(ip)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names.BrazilianPortuguese)
	fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names.English)
	fmt.Printf("Russian country name: %v\n", record.Country.Names.Russian)
	fmt.Printf("ISO country code: %v\n", record.Country.ISOCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	if record.Location.HasCoordinates() {
		fmt.Printf("Coordinates: %v, %v\n", *record.Location.Latitude, *record.Location.Longitude)
	} else {
		fmt.Println("Coordinates: unavailable")
	}
	// Output:
	// Portuguese (BR) city name: Londres
	// English subdivision name: England
	// Russian country name: Великобритания
	// ISO country code: GB
	// Time zone: Europe/London
	// Coordinates: 51.5142, -0.0931
}

// ExampleReader_City demonstrates how to use the City database.
func ExampleReader_City() {
	db, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("81.2.69.142")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.City(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("City: %v\n", record.City.Names.English)
	fmt.Printf("Country: %v (%v)\n", record.Country.Names.English, record.Country.ISOCode)
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	// Output:
	// City: London
	// Country: United Kingdom (GB)
	// Time zone: Europe/London
}

// ExampleReader_Country demonstrates how to use the Country database.
func ExampleReader_Country() {
	db, err := Open("test-data/test-data/GeoIP2-City-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("81.2.69.142")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.Country(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("Country: %v (%v)\n", record.Country.Names.English, record.Country.ISOCode)
	fmt.Printf("Continent: %v (%v)\n", record.Continent.Names.English, record.Continent.Code)
	// Output:
	// Country: United Kingdom (GB)
	// Continent: Europe (EU)
}

// ExampleReader_ASN demonstrates how to use the ASN database.
func ExampleReader_ASN() {
	db, err := Open("test-data/test-data/GeoLite2-ASN-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("1.128.0.0")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.ASN(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("ASN: %v\n", record.AutonomousSystemNumber)
	fmt.Printf("Organization: %v\n", record.AutonomousSystemOrganization)
	// Output:
	// ASN: 1221
	// Organization: Telstra Pty Ltd
}

// ExampleReader_AnonymousIP demonstrates how to use the Anonymous IP database.
func ExampleReader_AnonymousIP() {
	db, err := Open("test-data/test-data/GeoIP2-Anonymous-IP-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("1.2.0.0")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.AnonymousIP(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("Is Anonymous: %v\n", record.IsAnonymous)
	fmt.Printf("Is Anonymous VPN: %v\n", record.IsAnonymousVPN)
	fmt.Printf("Is Public Proxy: %v\n", record.IsPublicProxy)
	// Output:
	// Is Anonymous: true
	// Is Anonymous VPN: true
	// Is Public Proxy: false
}

// ExampleReader_AnonymousPlus demonstrates how to use the Anonymous Plus database.
func ExampleReader_AnonymousPlus() {
	db, err := Open("test-data/test-data/GeoIP-Anonymous-Plus-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("1.2.0.1")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.AnonymousPlus(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("Is Anonymous: %v\n", record.IsAnonymous)
	fmt.Printf("Is Anonymous VPN: %v\n", record.IsAnonymousVPN)
	fmt.Printf("Anonymizer Confidence: %v\n", record.AnonymizerConfidence)
	fmt.Printf("Provider Name: %v\n", record.ProviderName)
	fmt.Printf("Network Last Seen: %v\n", record.NetworkLastSeen.Format("2006-01-02"))
	// Output:
	// Is Anonymous: true
	// Is Anonymous VPN: true
	// Anonymizer Confidence: 30
	// Provider Name: foo
	// Network Last Seen: 2025-04-14
}

// ExampleReader_Enterprise demonstrates how to use the Enterprise database.
func ExampleReader_Enterprise() {
	db, err := Open("test-data/test-data/GeoIP2-Enterprise-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("74.209.24.0")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.Enterprise(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("City: %v\n", record.City.Names.English)
	fmt.Printf("Country: %v (%v)\n", record.Country.Names.English, record.Country.ISOCode)
	fmt.Printf("ISP: %v\n", record.Traits.ISP)
	fmt.Printf("Organization: %v\n", record.Traits.Organization)
	// Output:
	// City: Chatham
	// Country: United States (US)
	// ISP: Fairpoint Communications
	// Organization: Fairpoint Communications
}

// ExampleReader_ISP demonstrates how to use the ISP database.
func ExampleReader_ISP() {
	db, err := Open("test-data/test-data/GeoIP2-ISP-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("1.128.0.0")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.ISP(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("ISP: %v\n", record.ISP)
	fmt.Printf("Organization: %v\n", record.Organization)
	fmt.Printf("ASN: %v\n", record.AutonomousSystemNumber)
	// Output:
	// ISP: Telstra Internet
	// Organization: Telstra Internet
	// ASN: 1221
}

// ExampleReader_Domain demonstrates how to use the Domain database.
func ExampleReader_Domain() {
	db, err := Open("test-data/test-data/GeoIP2-Domain-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("1.2.0.0")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.Domain(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("Domain: %v\n", record.Domain)
	// Output:
	// Domain: maxmind.com
}

// ExampleReader_ConnectionType demonstrates how to use the Connection Type database.
func ExampleReader_ConnectionType() {
	db, err := Open("test-data/test-data/GeoIP2-Connection-Type-Test.mmdb")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	ip, err := netip.ParseAddr("1.0.128.0")
	if err != nil {
		log.Panic(err)
	}
	record, err := db.ConnectionType(ip)
	if err != nil {
		log.Panic(err)
	}

	if !record.HasData() {
		fmt.Println("No data found for this IP")
		return
	}

	fmt.Printf("Connection Type: %v\n", record.ConnectionType)
	// Output:
	// Connection Type: Cable/DSL
}
