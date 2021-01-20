// +build unit

package discovery

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestIsValidOpenInstallationPlatform(t *testing.T) {
	require.True(t, isValidOpenInstallationPlatform(string(types.OpenInstallationPlatformTypes.AMAZON)))
	require.False(t, isValidOpenInstallationPlatform("invalidValue"))
}

func TestIsValidOpenInstallationPlatformFamily(t *testing.T) {
	require.True(t, isValidOpenInstallationPlatformFamily(string(types.OpenInstallationPlatformFamilyTypes.SUSE)))
	require.False(t, isValidOpenInstallationPlatformFamily("invalidValue"))
}

func TestFilterValues_ValidPlatform(t *testing.T) {
	m := types.DiscoveryManifest{
		OS:             string(types.OpenInstallationOperatingSystemTypes.WINDOWS),
		Platform:       string(types.OpenInstallationPlatformTypes.AMAZON),
		PlatformFamily: string(types.OpenInstallationPlatformFamilyTypes.DEBIAN),
	}
	m = filterValues(m)

	require.Equal(t, string(types.OpenInstallationPlatformTypes.AMAZON), m.Platform)
	require.Equal(t, string(types.OpenInstallationPlatformFamilyTypes.DEBIAN), m.PlatformFamily)
}

func TestFilterValues_InvalidPlatform(t *testing.T) {
	m := types.DiscoveryManifest{
		OS:             string(types.OpenInstallationOperatingSystemTypes.WINDOWS),
		Platform:       "invalidValue",
		PlatformFamily: "invalidValue",
	}
	m = filterValues(m)

	require.Empty(t, m.Platform)
	require.Empty(t, m.PlatformFamily)
}
