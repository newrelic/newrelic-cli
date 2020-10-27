// +build integration

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())
}

func TestInitialization(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)
	config.DefaultConfigDirectory = f

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
	testingAPIKey := "mysupersecretAPIKey"
	testingRegion := "us"
	os.Setenv("NEW_RELIC_API_KEY", testingAPIKey)
	os.Setenv("NEW_RELIC_REGION", testingRegion)
	initializeCLI(Command, []string{})

	// // Initialize the new configuration directory
	c, err = credentials.LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.Profiles))
	assert.Equal(t, f, c.ConfigDirectory)
	assert.Equal(t, defaultProfileName, c.DefaultProfile)
	assert.Equal(t, testingAPIKey, c.Profiles[defaultProfileName].APIKey)
	assert.Equal(t, testingRegion, c.Profiles[defaultProfileName].Region)

	// Ensure that we don't Fatal out if the default profile already exists, but
	// was not specified in the default-profile.json.
	if err = os.Remove(fmt.Sprintf("%s/%s.json", f, credentials.DefaultProfileFile)); err != nil {
		t.Fatal(err)
	}

	initializeCLI(Command, []string{})
}
