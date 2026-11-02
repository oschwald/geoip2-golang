package geoip2

import (
	"net/netip"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShieldDatabases(t *testing.T) {
	openers := []struct {
		open func(string) (*Reader, error)
		name string
	}{
		{name: "Open", open: func(file string) (*Reader, error) { return Open(file) }},
		{name: "OpenBytes", open: func(file string) (*Reader, error) {
			data, err := os.ReadFile(file)
			if err != nil {
				return nil, err
			}
			return OpenBytes(data)
		}},
	}
	methods := []struct {
		lookup func(*Reader, netip.Addr) (any, error)
		name   string
	}{
		{name: "City", lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.City(ip) }},
		{
			name:   "Country",
			lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.Country(ip) },
		},
		{
			name:   "Enterprise",
			lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.Enterprise(ip) },
		},
		{name: "ASN", lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.ASN(ip) }},
		{name: "ISP", lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.ISP(ip) }},
		{
			name:   "AnonymousIP",
			lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.AnonymousIP(ip) },
		},
		{
			name:   "AnonymousPlus",
			lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.AnonymousPlus(ip) },
		},
		{
			name:   "ConnectionType",
			lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.ConnectionType(ip) },
		},
		{
			name:   "Domain",
			lookup: func(r *Reader, ip netip.Addr) (any, error) { return r.Domain(ip) },
		},
	}
	for _, family := range []struct {
		name string
		ipv4 string
		ipv6 string
	}{
		{"City", "81.2.69.160", "2001:218::"},
		{"Country", "81.2.69.160", "2001:218::"},
		{"Enterprise", "74.209.24.0", "2001:480::"},
	} {
		t.Run(family.name, func(t *testing.T) {
			base, err := Open("testdata/test-data/GeoIP2-" + family.name + "-Test.mmdb")
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, base.Close()) })
			for _, opener := range openers {
				t.Run(opener.name, func(t *testing.T) {
					file := "testdata/test-data/GeoIP2-" + family.name + "-Shield-Test.mmdb"
					shield, err := opener.open(file)
					if shield != nil {
						t.Cleanup(func() { require.NoError(t, shield.Close()) })
					}
					require.NoError(t, err)
					require.Equal(
						t,
						"GeoIP2-"+family.name+"-Shield",
						shield.Metadata().DatabaseType,
					)
					for _, method := range methods {
						t.Run(method.name, func(t *testing.T) {
							for _, address := range []string{
								family.ipv4,
								"::ffff:" + family.ipv4,
								family.ipv6,
								"192.168.1.1",
							} {
								t.Run(address, func(t *testing.T) {
									ip := netip.MustParseAddr(address)
									want, wantErr := method.lookup(base, ip)
									got, err := method.lookup(shield, ip)
									if wantErr != nil {
										var unsupported InvalidMethodError
										require.ErrorAs(t, wantErr, &unsupported)
										require.ErrorAs(t, err, &unsupported)
										require.Equal(t, method.name, unsupported.Method)
										require.Equal(
											t,
											shield.Metadata().DatabaseType,
											unsupported.DatabaseType,
										)
										require.Nil(t, got)
										return
									}
									require.NoError(t, err)
									require.Equal(t, want, got)
									result, ok := got.(interface{ HasData() bool })
									require.True(t, ok)
									require.Equal(t, address != "192.168.1.1", result.HasData())
								})
							}
						})
					}
				})
			}
		})
	}
}

func TestShieldSupportDoesNotEnableIPRisk(t *testing.T) {
	reader, err := Open("testdata/test-data/GeoIP2-IP-Risk-Test.mmdb")
	if reader != nil {
		t.Cleanup(func() { require.NoError(t, reader.Close()) })
	}
	var unknown UnknownDatabaseTypeError
	require.ErrorAs(t, err, &unknown)
	require.Equal(t, "GeoIP2-IP-Risk", unknown.DatabaseType)
}
