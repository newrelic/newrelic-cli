//go:build unit
// +build unit

package discovery

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestValidateOpenInstallationPlatform(t *testing.T) {
	require.Equal(t, validateOpenInstallationPlatform(string(types.OpenInstallationPlatformTypes.AMAZON)), "amazon")
	require.Equal(t, validateOpenInstallationPlatform(strings.ToLower(string(types.OpenInstallationPlatformTypes.AMAZON))), "amazon")
	require.Equal(t, validateOpenInstallationPlatform("invalidValue"), "")
}

func TestIsValidOpenInstallationPlatformFamily(t *testing.T) {
	require.True(t, isValidOpenInstallationPlatformFamily(string(types.OpenInstallationPlatformFamilyTypes.SUSE)))
	require.True(t, isValidOpenInstallationPlatformFamily(strings.ToLower(string(types.OpenInstallationPlatformFamilyTypes.SUSE))))
	require.False(t, isValidOpenInstallationPlatformFamily("invalidValue"))
}

func TestFilterValues_ValidPlatform(t *testing.T) {
	m := types.DiscoveryManifest{
		OS:             string(types.OpenInstallationOperatingSystemTypes.WINDOWS),
		Platform:       string(types.OpenInstallationPlatformTypes.AMAZON),
		PlatformFamily: string(types.OpenInstallationPlatformFamilyTypes.DEBIAN),
	}
	m = filterValues(m)

	require.Equal(t, strings.ToLower(string(types.OpenInstallationPlatformTypes.AMAZON)), m.Platform)
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
