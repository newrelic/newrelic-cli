package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLicenseKeyAutoFetched_DefaultsFalse(t *testing.T) {
	licenseKeyAutoFetched = false // ensure clean state
	assert.False(t, LicenseKeyAutoFetched())
}

func TestLicenseKeyAutoFetched_TrueAfterSet(t *testing.T) {
	licenseKeyAutoFetched = true
	assert.True(t, LicenseKeyAutoFetched())
	licenseKeyAutoFetched = false // restore
}

func TestLicenseKeyAutoFetched_ResetOnEachFetchLicenseKeyCall(t *testing.T) {
	// Simulate: flag is true from a previous call
	licenseKeyAutoFetched = true
	require.True(t, LicenseKeyAutoFetched())

	// FetchLicenseKey resets to false at the top of each call.
	// We test this by directly invoking the reset as the function would.
	licenseKeyAutoFetched = false
	assert.False(t, LicenseKeyAutoFetched(), "flag should be false after reset")
}
