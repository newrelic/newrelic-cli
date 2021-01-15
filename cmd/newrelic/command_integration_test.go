// +build integration

package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/config"
)

func TestInitializeProfile(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	defer os.RemoveAll(f)
	assert.NoError(t, err)

	config.ConfigDir = f

	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")

	// Init without the necessary environment variables
	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
	initializeProfile()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(config.GetProfileNames()))
	assert.Equal(t, "", config.GetDefaultProfileName())

	// Init with environment
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
	os.Setenv("NEW_RELIC_ACCOUNT_ID", envAccountID)
	initializeProfile()

	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")

	actualAPIKey, err := config.GetActiveProfileValue(config.APIKey)
	assert.NoError(t, err)

	actualRegion, err := config.GetActiveProfileValue(config.Region)
	assert.NoError(t, err)

	actualAccountID, err := config.GetActiveProfileValue(config.AccountID)
	assert.NoError(t, err)

	actualLicenseKey, err := config.GetActiveProfileValue(config.LicenseKey)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(config.GetProfileNames()))
	assert.Equal(t, defaultProfileName, config.GetDefaultProfileName())
	assert.Equal(t, apiKey, actualAPIKey)
	assert.NotEmpty(t, actualRegion)
	assert.NotEmpty(t, actualLicenseKey)
	assert.NotEmpty(t, actualAccountID)

	initializeProfile()
}
