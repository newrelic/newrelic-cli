// +build unit

package discovery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailWithoutVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsNoVersionMessage, result.Error())
}

func Test_ShouldFailVeryOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("5.3.1")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result.Error())
}

func Test_ShouldFailMinOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.0.0")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result.Error())
}

func Test_ShouldFailMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.0")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result.Error())
}

func Test_ShouldFailMinUnspecifiedVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.Equal(t, windowsVersionNoLongerSupported, result.Error())
}

func Test_ShouldPassMinVersionFull(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.1")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.NoError(t, result)
}

func Test_ShouldPassMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.1.0")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.NoError(t, result)
}

func Test_ShouldPass(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("10.0.14393")

	result := NewOsWindowsValidator().Execute(discover.GetManifest())
	require.NoError(t, result)
}
