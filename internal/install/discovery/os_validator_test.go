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
	discover.Os("darwin")
	result := NewOsValidator().Execute(discover.GetManifest())
	require.NotEqual(t, "", result)
}

func Test_FailsOnNoOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("")
	result := NewOsValidator().Execute(discover.GetManifest())
	require.Equal(t, noOperatingSystemDetected, result)
}

func Test_DoesntFailForWindowsOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("windows")
	result := NewOsValidator().Execute(discover.GetManifest())
	require.Equal(t, "", result)
}

func Test_DoesntFailForLinuxOs(t *testing.T) {
	discover := NewMockDiscoverer()
	discover.Os("linux")
	result := NewOsValidator().Execute(discover.GetManifest())
	require.Equal(t, "", result)
}
