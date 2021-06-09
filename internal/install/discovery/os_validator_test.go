// +build unit

package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	assert.True(t, true)
}

func Test_FailsOnInvalidOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("darwin")
	result := NewOsValidator().Validate(discover.GetManifest())
	require.Contains(t, result.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, result.Error(), "darwin")
	require.Contains(t, result.Error(), operatingSystemNotSupportedSuffix)
}

func Test_FailsOnNoOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("")
	result := NewOsValidator().Validate(discover.GetManifest())
	require.Equal(t, noOperatingSystemDetected, result.Error())
}

func Test_DoesntFailForWindowsOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("windows")
	result := NewOsValidator().Validate(discover.GetManifest())
	require.NoError(t, result)
}

func Test_DoesntFailForLinuxOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.SetOs("linux")
	result := NewOsValidator().Validate(discover.GetManifest())
	require.NoError(t, result)
}
