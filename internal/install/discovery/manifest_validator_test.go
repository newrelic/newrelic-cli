//go:build unit
// +build unit

package discovery

import (
	"context"
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"

	"github.com/stretchr/testify/mock"

	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
)

func Test_ShouldSucceedManifestOsDarwin(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.14",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	err := NewManifestValidator().Validate(manifest)
	require.NoError(t, err)
}

func Test_ShouldFailManifestOsDarwinOldVersion(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.13",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	err := NewManifestValidator().Validate(manifest)
	require.Error(t, err)
	require.Contains(t, err.Error(), errorPrefix)
	require.Contains(t, err.Error(), "darwin")
}

func Test_ShouldFailManifestWindowsVersion(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "1",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	result := NewManifestValidator().Validate(manifest)
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "windows"))
}

func Test_ShouldFailManifestUbuntuVersion(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "10.14",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	result := NewManifestValidator().Validate(manifest)
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"))
}
