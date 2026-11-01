//go:build linux && !appengine && !plan9 && !js && !wasip1 && !wasi

package geoip2

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oschwald/maxminddb-golang/v2"
)

func TestOpenUnsupportedDatabaseUnmapsFile(t *testing.T) {
	data, err := os.ReadFile("testdata/test-data/MaxMind-DB-test-decoder.mmdb")
	require.NoError(t, err)
	file := filepath.Join(t.TempDir(), "unsupported.mmdb")
	// #nosec G703 -- The path is within the test's temporary directory.
	require.NoError(t, os.WriteFile(file, data, 0o600))

	// Verify that this file is mapped and visible in /proc/self/maps.
	mapped, err := maxminddb.Open(file)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, mapped.Close())
	})
	maps, err := os.ReadFile("/proc/self/maps")
	if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
		t.Skipf("cannot inspect memory maps: %v", err)
	}
	require.NoError(t, err)
	if !strings.Contains(string(maps), file) {
		t.Skip("preflight reader uses buffered reads instead of mmap")
	}
	require.NoError(t, mapped.Close())
	maps, err = os.ReadFile("/proc/self/maps")
	require.NoError(t, err)
	require.NotContains(t, string(maps), file)

	reader, err := Open(file)
	if reader != nil {
		t.Cleanup(func() {
			assert.NoError(t, reader.Close())
		})
	}
	require.Error(t, err)
	maps, err = os.ReadFile("/proc/self/maps")
	// Keep a baseline reader alive until after inspection, so GC cannot hide the bug.
	runtime.KeepAlive(reader)
	require.NoError(t, err)
	assert.NotContains(t, string(maps), file, "rejected Open must unmap the file before returning")
}
