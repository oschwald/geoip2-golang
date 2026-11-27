package geoip2

import (
	"bytes"
	"encoding/json"
	"net/netip"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenUnsupportedDatabase(t *testing.T) {
	const file = "testdata/test-data/MaxMind-DB-test-decoder.mmdb"
	data, err := os.ReadFile(file)
	require.NoError(t, err)

	for _, test := range []struct {
		open func() (*Reader, error)
		name string
	}{
		{func() (*Reader, error) { return Open(file) }, "Open"},
		{func() (*Reader, error) { return OpenBytes(data) }, "OpenBytes"},
	} {
		t.Run(test.name, func(t *testing.T) {
			reader, err := test.open()
			if reader != nil {
				t.Cleanup(func() {
					assert.NoError(t, reader.Close())
				})
			}

			assert.Nil(t, reader, "unsupported databases must not return a reader")
			var unknownType UnknownDatabaseTypeError
			require.ErrorAs(t, err, &unknownType)
			assert.Equal(t, "MaxMind DB Decoder Test", unknownType.DatabaseType)
			assert.EqualError(t, err,
				`geoip2: reader does not support the "MaxMind DB Decoder Test" database type`)
		})
	}
}

// TestClose checks the behavior that the Reader and Close docs promise:
// lookup results stay valid after Close and do not reference the database,
// lookups after Close return an error, and Close leaves a slice passed to
// OpenBytes unchanged. The Open case covers the unmapped-file case only on
// platforms where Open uses a memory map.
func TestClose(t *testing.T) {
	const file = "testdata/test-data/GeoIP2-City-Test.mmdb"
	data, err := os.ReadFile(file)
	require.NoError(t, err)

	for _, test := range []struct {
		open func(t *testing.T) (*Reader, []byte)
		name string
	}{
		{
			name: "Open",
			open: func(t *testing.T) (*Reader, []byte) {
				reader, err := Open(file)
				require.NoError(t, err)
				return reader, nil
			},
		},
		{
			name: "OpenBytes",
			open: func(t *testing.T) (*Reader, []byte) {
				buf := bytes.Clone(data)
				reader, err := OpenBytes(buf)
				require.NoError(t, err)
				return reader, buf
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			reader, buf := test.open(t)
			closed := false
			t.Cleanup(func() {
				if !closed {
					assert.NoError(t, reader.Close())
				}
			})

			ip := netip.MustParseAddr("81.2.69.160")
			record, err := reader.City(ip)
			require.NoError(t, err)
			before, err := json.Marshal(record)
			require.NoError(t, err)

			err = reader.Close()
			closed = true
			require.NoError(t, err)

			if buf != nil {
				assert.Equal(t, data, buf, "Close must not modify the OpenBytes slice")
				// Overwrite the database bytes. Results must not reference them.
				clear(buf)
			}

			assert.Equal(t, "London", record.City.Names.English)
			assert.Equal(t, netip.MustParsePrefix("81.2.69.160/27"), record.Traits.Network)
			after, err := json.Marshal(record)
			require.NoError(t, err)
			assert.JSONEq(t, string(before), string(after))

			_, err = reader.City(ip)
			require.ErrorContains(t, err, "closed database",
				"lookups after Close must return the closed-database error")
		})
	}
}
