package discovery

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldFailUbuntuWithoutVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:       "linux",
		Platform: "ubuntu",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(noVersionMessage, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuVeryOldVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "12.04",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinOldVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.03",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldFailUbuntuMinUnspecifiedVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"), result.Error())
}

func Test_ShouldPassUbuntuMinVersionFull(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.04.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassUbuntuMinVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "16.04",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassUbuntu(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "20.04",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("linux", "ubuntu", 16, 04).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldFailWindowsWithoutVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS: "windows",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(noVersionMessage, "windows"), result.Error())
}

func Test_ShouldFailWindowsVeryOldVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "5.3.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsOldWithAnyPlatform(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		Platform:        "Anything possible",
		PlatformVersion: "5.3.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows/Anything possible"), result.Error())
}

func Test_ShouldFailWindowsMinOldVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.1.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldFailWindowsMinUnspecifiedVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.Equal(t, fmt.Sprintf(versionNoLongerSupported, "windows"), result.Error())
}

func Test_ShouldPassWindowsMinVersionFull(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.2",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassWindowsMinVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "6.2.0",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.NoError(t, result)
}

func Test_ShouldPassWindows(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "10.0.14393",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsVersionValidator("windows", "", 6, 2).Validate(manifest)

	require.NoError(t, result)
}
