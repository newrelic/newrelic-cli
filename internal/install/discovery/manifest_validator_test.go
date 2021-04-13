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
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, result.Error(), "darwin")
	require.Contains(t, result.Error(), operatingSystemNotSupportedSuffix)
}

func Test_ShouldFailWindowsVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("1.0")

	result := NewManifestValidator().Execute(discover.GetManifest())
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), windowsVersionNoLongerSupported)
}
