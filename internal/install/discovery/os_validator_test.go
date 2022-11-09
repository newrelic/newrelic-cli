package discovery

import (
	"context"
	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	assert.True(t, true)
}

func Test_FailsOnInvalidOs(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS: "freebsd",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	err := NewOsValidator().Validate(manifest)
	require.Error(t, err)
	require.Contains(t, err.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, err.Error(), "freebsd")
}

func Test_FailsOnNoOs(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS: "",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsValidator().Validate(manifest)
	require.Equal(t, noOperatingSystemDetected, result.Error())
}

func Test_DoesntFailForWindowsOs(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS: "windows",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}

func Test_DoesntFailForLinuxOs(t *testing.T) {
	context := context.Background()
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.EXPECT().Discover(context).Return(&types.DiscoveryManifest{
		OS: "linux",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context)

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}
