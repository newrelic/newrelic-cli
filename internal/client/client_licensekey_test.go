package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
)

func TestLicenseKeyAutoFetched_DefaultsFalse(t *testing.T) {
	licenseKeyAutoFetched = false
	assert.False(t, LicenseKeyAutoFetched())
}

func TestLicenseKeyAutoFetched_TrueAfterSet(t *testing.T) {
	licenseKeyAutoFetched = true
	assert.True(t, LicenseKeyAutoFetched())
	licenseKeyAutoFetched = false // restore
}

// TestFetchLicenseKey_SetsFlag verifies that FetchLicenseKey sets
// licenseKeyAutoFetched to true when a key is successfully found.
func TestFetchLicenseKey_SetsFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"actor":{"apiAccess":{"keySearch":{"keys":[{"key":"test-license-key","name":"Installer Ingest License Key","type":"INGEST"}],"nextCursor":""}}}}}`)
	}))
	defer ts.Close()

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey("test-api-key"),
		newrelic.ConfigLogLevel("none"),
		newrelic.ConfigNerdGraphBaseURL(ts.URL),
	)
	require.NoError(t, err)

	licenseKeyAutoFetched = false
	NRClient = nrClient
	t.Cleanup(func() { NRClient = nil; licenseKeyAutoFetched = false })

	maxTimeout := 5
	key, err := FetchLicenseKey(12345, "", &maxTimeout)
	require.NoError(t, err)
	assert.Equal(t, "test-license-key", key)
	assert.True(t, LicenseKeyAutoFetched(), "flag must be true after a successful FetchLicenseKey")
}

// TestFetchLicenseKey_ResetsFlag verifies that FetchLicenseKey resets
// licenseKeyAutoFetched to false at the start of each call, even when the
// flag was true from a previous successful run.
// Uses a bad profile name to trigger an immediate early-exit error before
// the retry loop — avoiding any sleep delay in the test.
func TestFetchLicenseKey_ResetsFlag(t *testing.T) {
	// Pre-set to true to simulate a prior successful run.
	licenseKeyAutoFetched = true
	t.Cleanup(func() { licenseKeyAutoFetched = false })

	maxTimeout := 5
	// A non-empty profileName that has no credentials triggers NewClient to fail
	// immediately — no retry loop, no sleep.
	_, err := FetchLicenseKey(12345, "nonexistent-profile-for-test", &maxTimeout)
	assert.Error(t, err)
	assert.False(t, LicenseKeyAutoFetched(), "flag must be reset to false at the start of FetchLicenseKey")
}
