//go:build integration

package client

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
)

func TestClientFetchLicenseKey(t *testing.T) {
	t.Parallel()

	testAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if testAccountID == "" {
		t.Skipf("New Relic internal testing account required")
	}

	acctID, err := strconv.Atoi(testAccountID)
	if err != nil {
		t.Skipf("error converting NEW_RELIC_ACCOUNT_ID to integer")
	}

	maxTimeoutSeconds := config.DefaultMaxTimeoutSeconds
	result, err := FetchLicenseKey(acctID, "default", &maxTimeoutSeconds)
	require.NoError(t, err)
	require.NotNil(t, result)
}
