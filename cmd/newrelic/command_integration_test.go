// +build integration

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
)

func TestInitializeProfile(t *testing.T) {
	envUserKey := os.Getenv("NEW_RELIC_API_KEY")
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if envUserKey == "" || envAccountID == "" {
		t.Skipf("NEW_RELIC_API_KEY and NEW_RELIC_ACCOUNT_ID are required to run this test")
	}

	f, err := ioutil.TempDir("/tmp", "newrelic")
	defer os.RemoveAll(f)
	require.NoError(t, err)

	config.ConfigDir = f

	// Init without the necessary environment variables
	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
	initializeDefaultProfile()

	require.NoError(t, err)
	require.Equal(t, 0, len(config.GetProfileNames()))
	require.Equal(t, "", config.GetDefaultProfileName())

	// Init with environment
	os.Setenv("NEW_RELIC_API_KEY", envUserKey)
	os.Setenv("NEW_RELIC_ACCOUNT_ID", envAccountID)
	initializeDefaultProfile()

	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")

	actualUserKey, err := config.GetActiveProfileValue(config.UserKey)
	require.NoError(t, err)

	actualRegion, err := config.GetActiveProfileValue(config.Region)
	require.NoError(t, err)

	actualAccountID, err := config.GetActiveProfileValue(config.AccountID)
	require.NoError(t, err)

	actualLicenseKey, err := config.GetActiveProfileValue(config.LicenseKey)
	require.NoError(t, err)

	require.Equal(t, 1, len(config.GetProfileNames()))
	require.Equal(t, config.DefaultDefaultProfileName, config.GetDefaultProfileName())
	require.Equal(t, envUserKey, actualUserKey)
	require.NotEmpty(t, actualRegion)
	require.NotEmpty(t, actualLicenseKey)
	require.NotEmpty(t, actualAccountID)

	initializeDefaultProfile()
}
