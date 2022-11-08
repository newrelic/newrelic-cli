package discovery

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_FailsOnInvalidOs(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS: "freebsd",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	err := NewOsValidator().Validate(manifest)
	require.Error(t, err)
	require.Contains(t, err.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, err.Error(), "freebsd")
}

func Test_FailsOnNoOs(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS: "",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsValidator().Validate(manifest)
	require.Equal(t, noOperatingSystemDetected, result.Error())
}

func Test_DoesntFailForWindowsOs(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS: "windows",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}

func Test_DoesntFailForLinuxOs(t *testing.T) {
	setup(t)
	mockDiscoverer.EXPECT().Discover(mockContext).Return(&types.DiscoveryManifest{
		OS: "linux",
	}, nil)
	manifest, _ := mockDiscoverer.Discover(mockContext)

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}
