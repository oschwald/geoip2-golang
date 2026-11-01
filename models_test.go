package geoip2

import (
	"encoding/json"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDateJSONRoundTrip(t *testing.T) {
	for _, test := range []struct {
		name string
		date time.Time
		json string
	}{
		{"zero", time.Time{}, `null`},
		{"leap_day", time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC), `"2024-02-29"`},
		{"leap_century", time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC), `"2000-02-29"`},
		{"lower_boundary", time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), `"0000-01-01"`},
		{"upper_boundary", time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC), `"9999-12-31"`},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := Date{Time: test.date}

			data, err := json.Marshal(original)
			require.NoError(t, err)
			assert.JSONEq(t, test.json, string(data))

			var decoded Date
			require.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestDateUnmarshalJSONCompatibility(t *testing.T) {
	for _, input := range []string{
		`null`,
		`"2024-02-29T12:34:56Z"`,
		`"2024-02-29T12:34:56.123456789+05:30"`,
		`"2024-02-29T12:34:56-07:00"`,
	} {
		t.Run(input, func(t *testing.T) {
			for _, initial := range []time.Time{
				{},
				time.Date(2020, time.June, 15, 12, 30, 45, 123, time.UTC),
			} {
				expected := initial
				require.NoError(t, json.Unmarshal([]byte(input), &expected))

				decoded := Date{Time: initial}
				require.NoError(t, json.Unmarshal([]byte(input), &decoded))
				assert.Equal(t, expected, decoded.Time)
			}
		})
	}
}

func TestDateUnmarshalJSONInvalid(t *testing.T) {
	for _, input := range []string{
		`""`,
		`"2023-02-29"`,
		`"1900-02-29"`,
		`"2024-02-30"`,
		`"2024-00-01"`,
		`"2024-13-01"`,
		`"2024-01-00"`,
		`"2024-01-32"`,
		`"2024-2-29"`,
		`"2024-02-29T00:00:00"`,
		`"2024-02-30T00:00:00Z"`,
		`"2024-02-29T25:00:00Z"`,
		`"2024-02-29T00:00:00+25:00"`,
		`"not a date"`,
		`123`,
		`true`,
		`{}`,
		`[]`,
		`"2024-02-29`,
		`"2024-02-29" trailing`,
	} {
		t.Run(input, func(t *testing.T) {
			original := Date{Time: time.Date(2020, time.June, 15, 12, 30, 45, 123, time.UTC)}
			decoded := original

			require.Error(t, json.Unmarshal([]byte(input), &decoded))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestAnonymousPlusJSONRoundTrip(t *testing.T) {
	original := AnonymousPlus{
		IPAddress:            netip.MustParseAddr("192.0.2.1"),
		Network:              netip.MustParsePrefix("192.0.2.0/24"),
		NetworkLastSeen:      Date{Time: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC)},
		ProviderName:         "Example VPN",
		AnonymizerConfidence: 75,
		IsAnonymous:          true,
		IsAnonymousVPN:       true,
		IsHostingProvider:    true,
		IsPublicProxy:        true,
		IsResidentialProxy:   true,
		IsTorExitNode:        true,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded AnonymousPlus
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, original, decoded)
}
