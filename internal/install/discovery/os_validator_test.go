//go:build unit
// +build unit

package discovery

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
)

func Test_FailsOnInvalidOs(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS: "freebsd",
	}

	err := NewOsValidator().Validate(manifest)

	require.Error(t, err)
	require.Contains(t, err.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, err.Error(), "freebsd")
}

func Test_FailsOnNoOs(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS: "",
	}

	result := NewOsValidator().Validate(manifest)
	require.Equal(t, noOperatingSystemDetected, result.Error())
}

func Test_DoesntFailForWindowsOs(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS: "windows",
	}

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}

func Test_DoesntFailForLinuxOs(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS: "linux",
	}

	result := NewOsValidator().Validate(manifest)
	require.NoError(t, result)
}
