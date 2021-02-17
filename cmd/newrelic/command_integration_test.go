// +build integration

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
)

func TestInitializeProfile(t *testing.T) {

	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if apiKey == "" || envAccountID == "" {
		t.Skipf("NEW_RELIC_API_KEY and NEW_RELIC_ACCOUNT_ID are required to run this test")
	}

	f, err := ioutil.TempDir("/tmp", "newrelic")
	defer os.RemoveAll(f)
	assert.NoError(t, err)
	config.DefaultConfigDirectory = f

	// Init without the necessary environment variables
	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
	initializeProfile()

	// Load credentials from disk
	c, err := credentials.LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c.Profiles))
	assert.Equal(t, f, c.ConfigDirectory)
	assert.Equal(t, "", c.DefaultProfile)

	// Init with environment
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
	os.Setenv("NEW_RELIC_ACCOUNT_ID", envAccountID)
	initializeProfile()

	// Initialize the new configuration directory
	c, err = credentials.LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.Profiles))
	assert.Equal(t, f, c.ConfigDirectory)
	assert.Equal(t, defaultProfileName, c.DefaultProfile)
	assert.Equal(t, apiKey, c.Profiles[defaultProfileName].APIKey)
	assert.NotEmpty(t, c.Profiles[defaultProfileName].Region)
	assert.NotEmpty(t, c.Profiles[defaultProfileName].AccountID)

	// Ensure that we don't Fatal out if the default profile already exists, but
	// was not specified in the default-profile.json.
	if err = os.Remove(fmt.Sprintf("%s/%s.json", f, credentials.DefaultProfileFile)); err != nil {
		t.Fatal(err)
	}

	initializeProfile()
}
