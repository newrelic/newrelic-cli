package discovery

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldSucceedManifestOsDarwin(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.14",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	err := NewManifestValidator().Validate(manifest)
	require.NoError(t, err)
}

func Test_ShouldFailManifestOsDarwinOldVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.13",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	err := NewManifestValidator().Validate(manifest)
	require.Error(t, err)
	require.Contains(t, err.Error(), errorPrefix)
	require.Contains(t, err.Error(), "darwin")
}

func Test_ShouldFailManifestWindowsVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewManifestValidator().Validate(manifest)
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "windows"))
}

func Test_ShouldFailManifestUbuntuVersion(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "10.14",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewManifestValidator().Validate(manifest)
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"))
}
