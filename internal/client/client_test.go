//go:build integration

package client

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apiaccess"
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

func TestClientGetLicenseKey(t *testing.T) {

	keys := []apiaccess.APIKey{
		{
			APIAccessKey: apiaccess.APIAccessKey{
				Key:       "key1",
				CreatedAt: 1,
			},
		},
		{
			APIAccessKey: apiaccess.APIAccessKey{
				Key:       "key2",
				CreatedAt: 2,
			},
		},
	}

	key := getPreferredLicenseKey(keys)

	require.Equal(t, "key1", key, "Get License Key should return earlist if no prefer key name")

	keys = append(keys, apiaccess.APIKey{
		APIAccessKey: apiaccess.APIAccessKey{
			Key:       "preferKey1",
			CreatedAt: 3,
			Name:      PreferredIngestKeyName,
		},
	})

	keys = append(keys, apiaccess.APIKey{
		APIAccessKey: apiaccess.APIAccessKey{
			Key:       "preferKey2",
			CreatedAt: 4,
			Name:      PreferredIngestKeyName,
		},
	})

	key = getPreferredLicenseKey(keys)

	require.Equal(t, "preferKey1", key, "Get License Key should return earlist key with prefer key name")
}
