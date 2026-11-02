package geoip2

import (
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
