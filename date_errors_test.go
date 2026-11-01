package geoip2

import (
	"bytes"
	"fmt"
	"net/netip"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2"
	"github.com/oschwald/maxminddb-golang/v2/mmdbdata"
)

func TestDateWrongMMDBKind(t *testing.T) {
	for _, test := range []struct {
		name string
		data []byte
		kind mmdbdata.Kind
	}{
		{name: "bool", data: []byte{0x01, 0x07}, kind: mmdbdata.KindBool},
		{name: "bytes", data: []byte{0x81, 'x'}, kind: mmdbdata.KindBytes},
		{name: "uint16", data: []byte{0xa1, 42}, kind: mmdbdata.KindUint16},
		{name: "map", data: []byte{0xe0}, kind: mmdbdata.KindMap},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, entryPoint := range []string{"cursor", "decoder"} {
				t.Run(entryPoint, func(t *testing.T) {
					decoder := mmdbdata.NewDecoder(append([]byte{0}, test.data...), 1)
					want := Date{Time: time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)}
					got := want
					var err error
					if entryPoint == "cursor" {
						var next mmdbdata.Cursor
						next, err = got.UnmarshalMaxMindDBCursor(decoder.Cursor())
						require.Zero(t, next)
					} else {
						err = got.UnmarshalMaxMindDB(decoder)
					}
					require.Equal(t, want, got)
					require.Equal(t, uint(1), decoder.Offset())
					require.ErrorContains(t, err, "at offset 1:")
					requireDateTypeError(t, err, test.kind)
				})
			}
		})
	}
}

func TestAnonymousPlusWrongDateMMDBKind(t *testing.T) {
	data, err := os.ReadFile("testdata/test-data/GeoIP-Anonymous-Plus-Test.mmdb")
	require.NoError(t, err)
	encodedDate := append([]byte{0x4a}, "2025-04-14"...)
	require.Equal(t, 1, bytes.Count(data, encodedDate))
	offset := bytes.Index(data, encodedDate)
	// Change only the string kind to bytes. Keep the fixture layout intact.
	data[offset] = 0x8a
	reader, err := OpenBytes(data)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, reader.Close()) })
	metadata := reader.Metadata()
	dataStart := int(metadata.NodeCount)*int(metadata.RecordSize)/4 + 16
	address := netip.MustParseAddr("1.2.0.1")

	record, err := reader.AnonymousPlus(address)
	require.NotNil(t, record)
	require.True(t, record.NetworkLastSeen.IsZero())
	require.ErrorContains(t, err, "decoding field network_last_seen:")
	require.ErrorContains(t, err, fmt.Sprintf("at offset %d:", offset-dataStart))
	requireDateTypeError(t, err, mmdbdata.KindBytes)

	want := Date{Time: time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)}
	var seeded AnonymousPlus
	seeded.NetworkLastSeen = want
	err = reader.mmdbReader.Lookup(address).Decode(&seeded)
	require.Equal(t, want, seeded.NetworkLastSeen)
	require.ErrorContains(t, err, "decoding field network_last_seen:")
	requireDateTypeError(t, err, mmdbdata.KindBytes)
}

func TestDateInvalidMMDBData(t *testing.T) {
	for _, test := range []struct {
		name string
		data []byte
	}{
		{name: "truncated string", data: []byte{0x4a, '2'}},
		{name: "truncated header", data: []byte{0x5d}},
		{name: "invalid pointer", data: []byte{0x20, 0xff}},
		{name: "oversized string", data: mmdbString(maxDateSize + 1)},
	} {
		t.Run(test.name, func(t *testing.T) {
			for _, entryPoint := range []string{"cursor", "decoder"} {
				t.Run(entryPoint, func(t *testing.T) {
					decoder := mmdbdata.NewDecoder(test.data, 0)
					want := Date{Time: time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)}
					got := want
					var err error
					if entryPoint == "cursor" {
						var next mmdbdata.Cursor
						next, err = got.UnmarshalMaxMindDBCursor(decoder.Cursor())
						require.Zero(t, next)
					} else {
						err = got.UnmarshalMaxMindDB(decoder)
					}
					require.Equal(t, want, got)
					require.Zero(t, decoder.Offset())
					var invalidDatabase maxminddb.InvalidDatabaseError
					require.ErrorAs(t, err, &invalidDatabase)
					var typeError maxminddb.UnmarshalTypeError
					require.NotErrorAs(t, err, &typeError)
					var mismatch mmdbdata.UnexpectedKindError
					require.NotErrorAs(t, err, &mismatch)
				})
			}
		})
	}
}

func requireDateTypeError(t *testing.T, err error, kind mmdbdata.Kind) {
	t.Helper()
	var typeError maxminddb.UnmarshalTypeError
	require.ErrorAs(t, err, &typeError)
	require.Equal(t, reflect.TypeFor[Date](), typeError.Type)
	require.Equal(t, fmt.Sprintf("%v (%T)", kind, kind), typeError.Value)
	var invalidDatabase maxminddb.InvalidDatabaseError
	require.NotErrorAs(t, err, &invalidDatabase)
	var mismatch mmdbdata.UnexpectedKindError
	require.NotErrorAs(t, err, &mismatch)
}
