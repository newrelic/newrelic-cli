package discovery

import (
	"fmt"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailUbuntuWithoutVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:       "linux",
		Platform: "ubuntu",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.Equal(t, fmt.Sprintf(noVersionMessage, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuVeryOldVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "12.04",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinOldVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.03",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinUnspecifiedVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldPassUbuntuMinVersionFull(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.04.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.NoError(t, result)
}

func Test_ShouldPassUbuntuMinVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.04",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.NoError(t, result)
}

func Test_ShouldPassUbuntu(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "20.04",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)
	require.NoError(t, result)
}

func Test_ShouldFailWindowsWithoutVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS: "windows",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.Equal(t, fmt.Sprintf(noVersionMessage, "windows"), result.Error())
}

func Test_ShouldFailWindowsVeryOldVersion(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "3.0.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsOldWithAnyPlatform(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		Platform:        "anything-possible",
		PlatformVersion: "3.0.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows/anything-possible"), result.Error())
}

func Test_ShouldFailWindowsMinOldVersion(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.1.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinVersion(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinUnspecifiedVersion(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.0.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldPassWindowsMinVersionFull(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.2",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.NoError(t, result)
}

func Test_ShouldPassWindowsMinVersion(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.2.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.NoError(t, result)
}

func Test_ShouldPassWindows(t *testing.T) {
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "10.0.14393",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)
	require.NoError(t, result)
}
