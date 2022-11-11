//go:build unit
// +build unit

package discovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/newrelic/newrelic-cli/internal/install/discovery/mocks"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
)

func Test_FailsOnInvalidOs(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS: "freebsd",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	err := NewOsValidator().Validate(manifest)
	require.Error(t, err)
	require.Contains(t, err.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, err.Error(), "freebsd")
}

func Test_FailsOnNoOs(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS: "",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	result := NewOsValidator().Validate(manifest)
	require.Equal(t, noOperatingSystemDetected, result.Error())
}

func Test_DoesntFailForWindowsOs(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS: "windows",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}

func Test_DoesntFailForLinuxOs(t *testing.T) {
	mockDiscoverer := mocks.NewDiscoverer(t)
	mockDiscoverer.On("Discover", mock.Anything).Return(&types.DiscoveryManifest{
		OS: "linux",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(context.Background())

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}
