// +build unit

package discovery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("darwin")

	result := NewManifestValidator().Execute(discover.GetManifest())
	require.Contains(t, result, errorPrefix)
	require.Contains(t, result, operatingSystemNotSupportedPrefix)
	require.Contains(t, result, "darwin")
	require.Contains(t, result, operatingSystemNotSupportedSuffix)
}

func Test_ShouldFailWindowsVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("1.0")

	result := NewManifestValidator().Execute(discover.GetManifest())
	require.Contains(t, result, errorPrefix)
	require.Contains(t, result, windowsVersionNoLongerSupported)
}
