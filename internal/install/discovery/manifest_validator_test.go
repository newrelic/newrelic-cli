// +build unit

package discovery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailManifestOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("darwin")

	result := NewManifestValidator().Execute(discover.GetManifest())
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, result.Error(), "darwin")
	require.Contains(t, result.Error(), operatingSystemNotSupportedSuffix)
}

func Test_ShouldFailManifestWindowsVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("1")

	result := NewManifestValidator().Execute(discover.GetManifest())
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "windows"))
}

func Test_ShouldFailManifestUbuntuVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("12.04")

	result := NewManifestValidator().Execute(discover.GetManifest())
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"))
}
