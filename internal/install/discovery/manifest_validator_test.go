package discovery

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-cli/internal/install/discovery/discoveryfakes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldSucceedManifestOsDarwin(t *testing.T) {
	discoverer := &discoveryfakes.FakeDiscoverer{}
	discoverer.DiscoverReturns(&types.DiscoveryManifest{
		OS:       "darwin",
		Platform: "10.14",
	}, nil)
	manifest, _ := discoverer.Discover(context.Background())

	err := NewManifestValidator().Validate(manifest)
	require.NoError(t, err)
}

func Test_ShouldFailManifestOsDarwinOldVersion(t *testing.T) {
	discoverer := &discoveryfakes.FakeDiscoverer{}
	discoverer.DiscoverReturns(&types.DiscoveryManifest{
		OS:       "darwin",
		Platform: "10.13",
	}, nil)
	manifest, _ := discoverer.Discover(context.Background())

	err := NewManifestValidator().Validate(manifest)
	require.Error(t, err)
	require.Contains(t, err.Error(), errorPrefix)
	require.Contains(t, err.Error(), "darwin")
}

func Test_ShouldFailManifestWindowsVersion(t *testing.T) {
	discoverer := &discoveryfakes.FakeDiscoverer{}
	discoverer.DiscoverReturns(&types.DiscoveryManifest{
		OS:       "windows",
		Platform: "1",
	}, nil)
	manifest, _ := discoverer.Discover(context.Background())

	result := NewManifestValidator().Validate(manifest)
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "windows"))
}

func Test_ShouldFailManifestUbuntuVersion(t *testing.T) {
	discoverer := &discoveryfakes.FakeDiscoverer{}
	discoverer.DiscoverReturns(&types.DiscoveryManifest{
		OS:             "linux",
		PlatformFamily: "ubuntu",
		Platform:       "12.04",
	}, nil)
	manifest, _ := discoverer.Discover(context.Background())

	result := NewManifestValidator().Validate(manifest)
	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"))
}
