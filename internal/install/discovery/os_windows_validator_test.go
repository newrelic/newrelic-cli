// +build unit

package discovery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailWindowsWithoutVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(noVersionMessage, "windows"), result.Error())
}

func Test_ShouldFailWindowsVeryOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("5.3.1")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsOldWithAnyPlatform(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatform("Anything possible")
	discover.SetPlatformVersion("5.3.1")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows/Anything possible"), result.Error())
}

func Test_ShouldFailWindowsMinOldVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.1.0")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.1")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinUnspecifiedVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldPassWindowsMinVersionFull(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.2")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.NoError(t, result)
}

func Test_ShouldPassWindowsMinVersion(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("6.2.0")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.NoError(t, result)
}

func Test_ShouldPassWindows(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	discover.SetPlatformVersion("10.0.14393")

	result := NewOsVersionValidator("windows", "", 6, 2).Execute(discover.GetManifest())
	require.NoError(t, result)
}
