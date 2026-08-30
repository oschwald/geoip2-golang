package geoip2

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2/mmdbdata"
)

func TestModelVariableSizeFieldsHaveMaxSize(t *testing.T) {
	roots := []reflect.Type{
		reflect.TypeFor[City](),
		reflect.TypeFor[Enterprise](),
		reflect.TypeFor[Country](),
		reflect.TypeFor[AnonymousIP](),
		reflect.TypeFor[AnonymousPlus](),
		reflect.TypeFor[ASN](),
		reflect.TypeFor[ConnectionType](),
		reflect.TypeFor[Domain](),
		reflect.TypeFor[ISP](),
	}
	modelPackage := reflect.TypeFor[City]().PkgPath()
	seen := make(map[reflect.Type]bool)
	var checkType func(reflect.Type, string)
	checkType = func(model reflect.Type, path string) {
		for model.Kind() == reflect.Pointer {
			model = model.Elem()
		}
		if model.Kind() != reflect.Struct || model.PkgPath() != modelPackage || seen[model] {
			return
		}
		seen[model] = true
		for fieldIndex := range model.NumField() {
			field := model.Field(fieldIndex)
			if !field.IsExported() {
				continue
			}
			fieldPath := path + "." + field.Name
			tag := field.Tag.Get("maxminddb")
			if tag == "-" {
				continue
			}
			fieldType := field.Type
			for fieldType.Kind() == reflect.Pointer {
				fieldType = fieldType.Elem()
			}
			switch fieldType.Kind() {
			case reflect.String, reflect.Slice, reflect.Map:
				maximum, ok := maxSizeFromTag(tag)
				require.True(t, ok, "%s must declare maxsize", fieldPath)
				require.Positive(t, maximum, "%s maxsize", fieldPath)
				require.LessOrEqual(t, maximum, 4096, "%s maxsize", fieldPath)
			case reflect.Struct:
				checkType(fieldType, fieldPath)
			}
		}
	}
	for _, root := range roots {
		checkType(root, root.Name())
	}
}

func maxSizeFromTag(tag string) (int, bool) {
	for option := range strings.SplitSeq(tag, ",") {
		value, ok := strings.CutPrefix(option, "maxsize:")
		if !ok {
			continue
		}
		maximum, err := strconv.Atoi(value)
		return maximum, err == nil
	}
	return 0, false
}

func TestGeneratedModelMaxSizeBounds(t *testing.T) {
	t.Run("date", func(t *testing.T) {
		var exact Date
		_, err := exact.UnmarshalMaxMindDBCursor(
			mmdbdata.NewDecoder(append([]byte{0x4a}, "2026-08-29"...), 0).Cursor(),
		)
		require.NoError(t, err)
		require.Equal(t, "2026-08-29", exact.Format("2006-01-02"))

		over := exact
		_, err = over.UnmarshalMaxMindDBCursor(
			mmdbdata.NewDecoder(mmdbString(maxDateSize+1), 0).Cursor(),
		)
		requireMaxSizeError(t, err)
		require.Equal(t, exact, over)
	})

	t.Run("localized name", func(t *testing.T) {
		decode := func(size int, destination *Names) error {
			data := append([]byte{0xe1, 0x42, 'e', 'n'}, mmdbString(size)...)
			_, err := destination.UnmarshalMaxMindDBCursor(mmdbdata.NewDecoder(data, 0).Cursor())
			return err
		}

		var exact Names
		require.NoError(t, decode(1024, &exact))
		require.Len(t, exact.English, 1024)

		over := Names{English: "keep"}
		err := decode(1025, &over)
		requireMaxSizeError(t, err)
		require.Equal(t, "keep", over.English)
	})

	t.Run("subdivisions", func(t *testing.T) {
		decode := func(size int, destination *City) error {
			data := []byte{0xe1, 0x4c, 's', 'u', 'b', 'd', 'i', 'v', 'i', 's', 'i', 'o', 'n', 's'}
			data = append(data, mmdbSliceHeader(size)...)
			for range size {
				data = append(data, 0xe0)
			}
			_, err := destination.UnmarshalMaxMindDBCursor(mmdbdata.NewDecoder(data, 0).Cursor())
			return err
		}

		var exact City
		require.NoError(t, decode(32, &exact))
		require.Len(t, exact.Subdivisions, 32)

		over := City{Subdivisions: []CitySubdivision{{GeoNameID: 42}}}
		err := decode(33, &over)
		requireMaxSizeError(t, err)
		require.Equal(t, []CitySubdivision{{GeoNameID: 42}}, over.Subdivisions)
	})
}

func TestGeneratedModelsRejectDuplicateFields(t *testing.T) {
	sharedCity := []byte{0xe1, 0x45, 'n', 'a', 'm', 'e', 's', 0xe1, 0x42, 'e', 'n', 0x43, 'o', 'n', 'e'}
	root := len(sharedCity)
	data := append(sharedCity, mmdbMapHeader(100)...)
	for range 100 {
		data = append(data, 0x44, 'c', 'i', 't', 'y', 0x20, 0x00)
	}

	var got City
	decoder := mmdbdata.NewDecoder(data, uint(root))
	_, err := got.UnmarshalMaxMindDBCursor(decoder.Cursor())
	var invalidDatabase mmdbdata.InvalidDatabaseError
	require.ErrorAs(t, err, &invalidDatabase)
	require.ErrorContains(t, err, "duplicate map key")
	require.Equal(t, "one", got.City.Names.English)
}

func requireMaxSizeError(t *testing.T, err error) {
	t.Helper()
	var invalidDatabase mmdbdata.InvalidDatabaseError
	require.True(t, errors.As(err, &invalidDatabase), "error = %v", err)
	require.ErrorContains(t, err, "exceeds maxsize")
}

func mmdbString(size int) []byte {
	return append(mmdbSizedHeader(0x40, size), make([]byte, size)...)
}

func mmdbSliceHeader(size int) []byte {
	header := mmdbSizedHeader(0, size)
	return append([]byte{header[0], 0x04}, header[1:]...)
}

func mmdbMapHeader(size int) []byte {
	return mmdbSizedHeader(0xe0, size)
}

func mmdbSizedHeader(kind byte, size int) []byte {
	switch {
	case size < 29:
		return []byte{kind | byte(size)}
	case size < 285:
		return []byte{kind | 29, byte(size - 29)}
	case size < 65821:
		value := size - 285
		return []byte{kind | 30, byte(value >> 8), byte(value)}
	default:
		value := size - 65821
		return []byte{kind | 31, byte(value >> 16), byte(value >> 8), byte(value)}
	}
}
