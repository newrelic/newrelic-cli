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
	discover.SetOs("freebsd")
	err := NewOsValidator().Validate(discover.GetManifest())
	require.Error(t, err)
	require.Contains(t, err.Error(), operatingSystemNotSupportedPrefix)
	require.Contains(t, err.Error(), "freebsd")
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
