package geoip2

import (
	"encoding"
	"encoding/json"
	"encoding/xml"
	"net/netip"
	"strconv"
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

var (
	_ encoding.TextMarshaler   = Date{Time: time.Time{}}
	_ encoding.TextAppender    = Date{Time: time.Time{}}
	_ encoding.TextUnmarshaler = (*Date)(nil)
)

func TestDateTextRoundTrip(t *testing.T) {
	for _, test := range []struct {
		name string
		date time.Time
		text string
	}{
		{"zero", time.Time{}, `0001-01-01`},
		{"leap_day", time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC), `2024-02-29`},
		{"leap_century", time.Date(2000, time.February, 29, 0, 0, 0, 0, time.UTC), `2000-02-29`},
		{"lower_boundary", time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC), `0000-01-01`},
		{"upper_boundary", time.Date(9999, time.December, 31, 0, 0, 0, 0, time.UTC), `9999-12-31`},
	} {
		t.Run(test.name, func(t *testing.T) {
			original := Date{Time: test.date}

			text, err := original.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, test.text, string(text))

			appended, err := original.AppendText([]byte("prefix:"))
			require.NoError(t, err)
			assert.Equal(t, "prefix:"+test.text, string(appended))

			decoded := Date{Time: time.Date(2020, time.June, 15, 12, 30, 45, 123, time.UTC)}
			require.NoError(t, decoded.UnmarshalText(text))
			assert.Equal(t, original, decoded)
		})
	}
}

func TestDateTextZeroInOtherLocation(t *testing.T) {
	for _, offset := range []int{-5 * 3600, 5 * 3600} {
		loc := time.FixedZone("", offset)
		t.Run(strconv.Itoa(offset), func(t *testing.T) {
			zero := Date{Time: time.Time{}.In(loc)}
			require.True(t, zero.IsZero())

			text, err := zero.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, "0001-01-01", string(text))

			var decoded Date
			require.NoError(t, decoded.UnmarshalText(text))
			assert.True(t, decoded.IsZero())
		})
	}
}

func TestDateMarshalYearOutOfRange(t *testing.T) {
	for _, date := range []time.Time{
		time.Date(-1, time.December, 31, 0, 0, 0, 0, time.UTC),
		time.Date(10000, time.January, 1, 0, 0, 0, 0, time.UTC),
	} {
		t.Run(date.String(), func(t *testing.T) {
			value := Date{Time: date}

			_, err := value.MarshalText()
			require.Error(t, err)

			_, err = value.AppendText([]byte("prefix:"))
			require.Error(t, err)

			_, err = json.Marshal(value)
			require.Error(t, err)
		})
	}
}

func TestDateTextMatchesJSON(t *testing.T) {
	date := Date{Time: time.Date(2025, time.January, 2, 3, 0, 0, 0, time.UTC)}

	text, err := date.MarshalText()
	require.NoError(t, err)
	data, err := json.Marshal(date)
	require.NoError(t, err)
	assert.Equal(t, `"`+string(text)+`"`, string(data))
	assert.Equal(t, "2025-01-02", string(text))

	// encoding/json writes map keys with MarshalText but reads them back with
	// UnmarshalJSON, so this checks that the two agree, including for the zero
	// Date. UnmarshalText is exercised by the XML round trip below.
	keyed, err := json.Marshal(map[Date]int{date: 1, {Time: time.Time{}}: 2})
	require.NoError(t, err)
	assert.JSONEq(t, `{"2025-01-02": 1, "0001-01-01": 2}`, string(keyed))

	var decodedMap map[Date]int
	require.NoError(t, json.Unmarshal(keyed, &decodedMap))
	assert.Equal(
		t,
		map[Date]int{
			{Time: time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC)}: 1,
			{Time: time.Time{}}: 2,
		},
		decodedMap,
	)

	type xmlRecord struct {
		Attr Date `xml:"attr,attr"`
		Elem Date `xml:"elem"`
	}
	xmlData, err := xml.Marshal(xmlRecord{Attr: date, Elem: date})
	require.NoError(t, err)
	assert.Equal(
		t,
		`<xmlRecord attr="2025-01-02"><elem>2025-01-02</elem></xmlRecord>`,
		string(xmlData),
	)
	var decodedXML xmlRecord
	require.NoError(t, xml.Unmarshal(xmlData, &decodedXML))
	day := Date{Time: time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC)}
	assert.Equal(t, xmlRecord{Attr: day, Elem: day}, decodedXML)
}

func TestDateUnmarshalTextCompatibility(t *testing.T) {
	for _, input := range []string{
		`2024-02-29T12:34:56Z`,
		`2024-02-29T12:34:56.123456789+05:30`,
		`2024-02-29T12:34:56-07:00`,
	} {
		t.Run(input, func(t *testing.T) {
			var expected time.Time
			require.NoError(t, expected.UnmarshalText([]byte(input)))

			decoded := Date{Time: time.Date(2020, time.June, 15, 12, 30, 45, 123, time.UTC)}
			require.NoError(t, decoded.UnmarshalText([]byte(input)))
			assert.Equal(t, expected, decoded.Time)

			var fromJSON Date
			require.NoError(t, json.Unmarshal([]byte(`"`+input+`"`), &fromJSON))
			assert.Equal(t, fromJSON, decoded)
		})
	}
}

func TestDateUnmarshalTextInvalid(t *testing.T) {
	for _, input := range []string{
		`2023-02-29`,
		`1900-02-29`,
		`2024-02-30`,
		`2024-00-01`,
		`2024-13-01`,
		`2024-01-00`,
		`2024-01-32`,
		`2024-2-29`,
		`2024-02-29T00:00:00`,
		`2024-02-30T00:00:00Z`,
		`2024-02-29T25:00:00Z`,
		`2024-02-29T00:00:00+25:00`,
		``,
		`not a date`,
		`"2024-02-29"`,
		` 2024-02-29`,
		`null`,
	} {
		t.Run(input, func(t *testing.T) {
			original := Date{Time: time.Date(2020, time.June, 15, 12, 30, 45, 123, time.UTC)}
			decoded := original

			require.Error(t, decoded.UnmarshalText([]byte(input)))
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

var dateBenchmarkSink []byte

func BenchmarkDateMarshalText(b *testing.B) {
	date := Date{Time: time.Date(2025, time.April, 14, 0, 0, 0, 0, time.UTC)}
	b.ReportAllocs()
	for b.Loop() {
		text, err := date.MarshalText()
		if err != nil {
			b.Fatal(err)
		}
		dateBenchmarkSink = text
	}
}

func BenchmarkDateMarshalJSON(b *testing.B) {
	date := Date{Time: time.Date(2025, time.April, 14, 0, 0, 0, 0, time.UTC)}
	b.ReportAllocs()
	for b.Loop() {
		data, err := date.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
		dateBenchmarkSink = data
	}
}
