package geoip2

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2/mmdbdata"
)

func TestDateUnmarshalMaxMindDB(t *testing.T) {
	methods := []struct {
		decode func(*testing.T, *Date, *mmdbdata.Decoder) error
		name   string
	}{
		{
			name: "UnmarshalMaxMindDBCursor",
			decode: func(t *testing.T, date *Date, decoder *mmdbdata.Decoder) error {
				next, err := date.UnmarshalMaxMindDBCursor(decoder.Cursor())
				if err != nil {
					return err
				}
				require.Zero(t, decoder.Offset(), "cursor must not advance the decoder")
				value, _, err := next.ReadString()
				require.NoError(t, err)
				require.Equal(t, "next", value)
				require.NoError(t, decoder.Advance(next), "successor must be valid for the decoder")
				return nil
			},
		},
		{
			name: "UnmarshalMaxMindDB",
			decode: func(t *testing.T, date *Date, decoder *mmdbdata.Decoder) error {
				if err := date.UnmarshalMaxMindDB(decoder); err != nil {
					return err
				}
				value, err := decoder.ReadString()
				require.NoError(t, err)
				require.Equal(t, "next", value)
				return nil
			},
		},
	}
	original := time.Date(2025, 4, 14, 12, 34, 56, 789, time.UTC)
	cases := []struct {
		want    time.Time
		name    string
		value   string
		invalid bool
	}{
		{name: "empty", value: "", want: time.Time{}, invalid: false},
		{
			name:    "valid",
			value:   "2024-02-29",
			want:    time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			invalid: false,
		},
		{name: "invalid syntax", value: "not-a-date", want: original, invalid: true},
		{name: "invalid calendar date", value: "2025-02-29", want: original, invalid: true},
	}
	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			for _, test := range cases {
				t.Run(test.name, func(t *testing.T) {
					data := append(mmdbSizedHeader(0x40, len(test.value)), test.value...)
					data = append(data, 0x44, 'n', 'e', 'x', 't')
					decoder := mmdbdata.NewDecoder(data, 0)
					date := Date{Time: original}
					err := method.decode(t, &date, decoder)
					if test.invalid {
						require.Error(t, err)
						require.Zero(
							t,
							decoder.Offset(),
							"invalid dates must not advance the decoder",
						)
					} else {
						require.NoError(t, err)
					}
					require.Equal(t, test.want, date.Time)
					if test.value == "" {
						require.True(
							t,
							date.IsZero(),
							"an explicitly empty date must clear the destination",
						)
					}
				})
			}
		})
	}
}

func TestAnonymousPlusMissingDatePreservesDestination(t *testing.T) {
	original := Date{Time: time.Date(2025, 4, 14, 0, 0, 0, 0, time.UTC)}
	t.Run("UnmarshalMaxMindDBCursor", func(t *testing.T) {
		var record AnonymousPlus
		record.NetworkLastSeen = original
		decoder := mmdbdata.NewDecoder([]byte{0xe0, 0x44, 'n', 'e', 'x', 't'}, 0)
		next, err := record.UnmarshalMaxMindDBCursor(decoder.Cursor())
		require.NoError(t, err)
		require.Equal(t, original, record.NetworkLastSeen)
		value, _, err := next.ReadString()
		require.NoError(t, err)
		require.Equal(t, "next", value)
	})
	t.Run("UnmarshalMaxMindDB", func(t *testing.T) {
		var record AnonymousPlus
		record.NetworkLastSeen = original
		decoder := mmdbdata.NewDecoder([]byte{0xe0, 0x44, 'n', 'e', 'x', 't'}, 0)
		require.NoError(t, record.UnmarshalMaxMindDB(decoder))
		require.Equal(t, original, record.NetworkLastSeen)
		value, err := decoder.ReadString()
		require.NoError(t, err)
		require.Equal(t, "next", value)
	})
}
