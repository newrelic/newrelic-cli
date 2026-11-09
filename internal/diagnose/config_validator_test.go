package diagnose

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
)

// TestSearchForConfiguredLicenseKey_SkipsWhenKeyAutoFetched verifies that
// searchForConfiguredLicenseKey returns nil immediately without making any
// API call when keyWasAutoFetched reports true.
func TestSearchForConfiguredLicenseKey_SkipsWhenKeyAutoFetched(t *testing.T) {
	cv := &ConfigValidator{
		keyWasAutoFetched: func() bool { return true },
	}
	err := cv.searchForConfiguredLicenseKey(context.Background())
	assert.NoError(t, err, "should return nil without any API call when key was auto-fetched")
}

// TestSearchForConfiguredLicenseKey_RunsFullScanWhenKeyUserProvided verifies
// that searchForConfiguredLicenseKey makes the key search API call and returns
// an InvalidIngestKey-typed error when the configured key is absent from the
// (empty) result list returned by the server.
func TestSearchForConfiguredLicenseKey_RunsFullScanWhenKeyUserProvided(t *testing.T) {
	// Mock server returns an empty key list (key not found).
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"data":{"actor":{"apiAccess":{"keySearch":{"keys":[],"nextCursor":""}}}}}`)
	}))
	defer ts.Close()

	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey("test-api-key"),
		newrelic.ConfigLogLevel("none"),
		newrelic.ConfigNerdGraphBaseURL(ts.URL),
	)
	require.NoError(t, err)

	cv := &ConfigValidator{
		client:            nrClient,
		keyWasAutoFetched: func() bool { return false },
	}

	err = cv.searchForConfiguredLicenseKey(context.Background())
	// Expect an error — the key was not found in the (empty) list.
	assert.Error(t, err)
	// The error details come from types.ErrLicenseKey.
	assert.Contains(t, err.Error(), "the configured license key is invalid")
}
