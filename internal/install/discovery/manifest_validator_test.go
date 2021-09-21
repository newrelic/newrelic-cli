//go:build unit
// +build unit

package discovery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldSucceedManifestOsDarwin(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("darwin")
	discover.SetPlatformVersion("10.14")

	err := NewManifestValidator().Validate(discover.GetManifest())
	require.NoError(t, err)
}

func Test_ShouldFailManifestOsDarwinOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("darwin")
	discover.SetPlatformVersion("10.13")

	err := NewManifestValidator().Validate(discover.GetManifest())
	require.Error(t, err)
	require.Contains(t, err.Error(), errorPrefix)
	require.Contains(t, err.Error(), "darwin")
}

func Test_ShouldFailManifestWindowsVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("1")

	result := NewManifestValidator().Validate(discover.GetManifest())
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "windows"))
}

func Test_ShouldFailManifestUbuntuVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("12.04")

	result := NewManifestValidator().Validate(discover.GetManifest())
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"))
}
