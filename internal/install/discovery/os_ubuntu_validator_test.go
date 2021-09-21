//go:build unit
// +build unit

package discovery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailUbuntuWithoutVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(noVersionMessage, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuVeryOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("12.04")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("16")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("16.03")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinUnspecifiedVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("16")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldPassUbuntuMinVersionFull(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("16.04.0")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.NoError(t, result)
}

func Test_ShouldPassUbuntuMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("16.04")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.NoError(t, result)
}

func Test_ShouldPassUbuntu(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	discover.SetPlatform("ubuntu")
	discover.SetPlatformVersion("20.04")

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(discover.GetManifest())
	require.NoError(t, result)
}
