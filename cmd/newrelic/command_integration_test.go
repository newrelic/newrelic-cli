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

func TestInitializeProfile(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	defer os.RemoveAll(f)
	assert.NoError(t, err)
	config.DefaultConfigDirectory = f

	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	apiRegion := os.Getenv("NEW_RELIC_REGION")
	apiAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	accountID, err := strconv.Atoi(apiAccountID)
	assert.NoError(t, err)

	// Init without the necessary environment variables
	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
	os.Setenv("NEW_RELIC_REGION", "")
	initializeProfile()

	// Load credentials from disk
	c, err := credentials.LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c.Profiles))
	assert.Equal(t, f, c.ConfigDirectory)
	assert.Equal(t, "", c.DefaultProfile)

	// Init with environment
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
	os.Setenv("NEW_RELIC_REGION", apiRegion)
	initializeProfile()

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

	initializeProfile()
}
