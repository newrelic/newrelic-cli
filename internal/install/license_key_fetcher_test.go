//go:build integration

package install

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func TestLicenseKeyFetcher_FetchLicenseKey(t *testing.T) {
	t.Parallel()

	testAccountAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	testAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if testAccountAPIKey == "" || testAccountID == "" {
		t.Skip("New Relic internal testing account required")
	}

	licenseKeyFetcher := NewServiceLicenseKeyFetcher(config.DefaultMaxTimeoutSeconds)

	result, err := licenseKeyFetcher.FetchLicenseKey(utils.SignalCtx)
	require.NoError(t, err)
	require.NotNil(t, result)
}
