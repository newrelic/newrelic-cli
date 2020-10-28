// +build integration

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
)

func TestInitialization(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	// defer os.RemoveAll(f)
	config.DefaultConfigDirectory = f

	// Save the creds while we have them.
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	apiRegion := os.Getenv("NEW_RELIC_REGION")
	apiAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	accountID, err := strconv.Atoi(apiAccountID)
	assert.NoError(t, err)

	// Init without the logical environment variables
	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_REGION", "")
	initializeCLI(Command, []string{})

	// // Initialize the new configuration directory
	c, err := credentials.LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c.Profiles))
	assert.Equal(t, f, c.ConfigDirectory)
	assert.Equal(t, "", c.DefaultProfile)

	// Init with environment
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
	os.Setenv("NEW_RELIC_REGION", apiRegion)
	initializeCLI(Command, []string{})

	// // Initialize the new configuration directory
	c, err = credentials.LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.Profiles))
	assert.Equal(t, f, c.ConfigDirectory)
	assert.Equal(t, defaultProfileName, c.DefaultProfile)
	assert.Equal(t, apiKey, c.Profiles[defaultProfileName].APIKey)
	assert.True(t, strings.EqualFold(apiRegion, c.Profiles[defaultProfileName].Region))
	assert.Equal(t, accountID, c.Profiles[defaultProfileName].AccountID)

	// Ensure that we don't Fatal out if the default profile already exists, but
	// was not specified in the default-profile.json.
	if err = os.Remove(fmt.Sprintf("%s/%s.json", f, credentials.DefaultProfileFile)); err != nil {
		t.Fatal(err)
	}

	initializeCLI(Command, []string{})
}
