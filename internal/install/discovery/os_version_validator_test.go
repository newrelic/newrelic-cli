//go:build unit
// +build unit

package discovery

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailUbuntuWithoutVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:       "linux",
		Platform: "ubuntu",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(noVersionMessage, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuVeryOldVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "12.04",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinOldVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.03",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinUnspecifiedVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldPassUbuntuMinVersionFull(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.04.0",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassUbuntuMinVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.04",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassUbuntu(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "20.04",
	}

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldFailWindowsWithoutVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS: "windows",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(noVersionMessage, "windows"), result.Error())
}

func Test_ShouldFailWindowsVeryOldVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "5.3.1",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsOldWithAnyPlatform(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		Platform:        "Anything possible",
		PlatformVersion: "5.3.1",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows/Anything possible"), result.Error())
}

func Test_ShouldFailWindowsMinOldVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.1.0",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.1",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinUnspecifiedVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldPassWindowsMinVersionFull(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.2",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassWindowsMinVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.2.0",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassWindows(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "10.0.14393",
	}

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.NoError(t, result)
}
