//go:build unit
// +build unit

package discovery

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
)

func Test_ShouldSucceedManifestOsDarwin(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.14",
	}

	err := NewManifestValidator().Validate(manifest)

	require.NoError(t, err)
}

func Test_ShouldFailManifestOsDarwinOldVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "darwin",
		PlatformVersion: "10.13",
	}

	err := NewManifestValidator().Validate(manifest)

	require.Error(t, err)
	require.Contains(t, err.Error(), errorPrefix)
	require.Contains(t, err.Error(), "darwin")
}

func Test_ShouldFailManifestWindowsVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "windows",
		PlatformVersion: "1",
	}

	result := NewManifestValidator().Validate(manifest)

	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "windows"))
}

func Test_ShouldFailManifestUbuntuVersion(t *testing.T) {
	manifest := &types.DiscoveryManifest{
		OS:              "linux",
		Platform:        "ubuntu",
		PlatformVersion: "10.14",
	}

	result := NewManifestValidator().Validate(manifest)

	require.Contains(t, result.Error(), errorPrefix)
	require.Contains(t, result.Error(), fmt.Sprintf(versionNoLongerSupported, "linux/ubuntu"))
}
