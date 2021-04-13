// +build unit

package discovery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailWithoutVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsNoVersionMessage, result)
}

func Test_ShouldFailVeryOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("5.3.1")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result)
}

func Test_ShouldFailMinOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("6.0.0")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result)
}

func Test_ShouldFailMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("6.0")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result)
}

func Test_ShouldFailMinUnspecifiedVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("6")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result)
}

func Test_ShouldPassMinVersionFull(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("6.1")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, "", result)
}

func Test_ShouldPassMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("6.1.0")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, "", result)
}

func Test_ShouldPass(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	discover.PlatformVersion("10.0.14393")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, "", result)
}
