// +build integration

package credentials

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentialsLoadCredentialsNoDirectory(t *testing.T) {
	c, err := LoadCredentials("/tmp/notexist")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 0)
	assert.Equal(t, c.DefaultProfile, "")
	assert.Equal(t, c.ConfigDirectory, "/tmp/notexist")
}

func TestCredentialsLoadCredentialsHomeDirectory(t *testing.T) {
	c, err := LoadCredentials("$HOME/.newrelictesting")
	assert.NoError(t, err)

	home := os.Getenv("HOME")
	filePath := fmt.Sprintf("%s/.newrelictesting", home)

	assert.Equal(t, len(c.Profiles), 0)
	assert.Equal(t, c.DefaultProfile, "")
	assert.Equal(t, c.ConfigDirectory, filePath)
}

func TestCredentials(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 0)
	assert.Equal(t, c.DefaultProfile, "")
	assert.Equal(t, c.ConfigDirectory, f)

	// Create an initial profile to work with
	err = c.AddProfile("testCase1", "us", "apiKeyGoesHere")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 1)
	assert.Equal(t, c.Profiles["testCase1"].Region, "us")
	assert.Equal(t, c.Profiles["testCase1"].PersonalAPIKey, "apiKeyGoesHere")
	assert.Equal(t, c.DefaultProfile, "")

	// Set the default profile to the only one we've got
	err = c.SetDefaultProfile("testCase1")
	assert.NoError(t, err)
	assert.Equal(t, c.DefaultProfile, "testCase1")

	// Adding a profile with the same name should result in an error
	err = c.AddProfile("testCase1", "us", "foot")
	assert.Error(t, err)
	assert.Equal(t, len(c.Profiles), 1)
	assert.True(t, c.profileExists("testCase1"))

	// Create a second profile to work with
	err = c.AddProfile("testCase2", "us", "apiKeyGoesHere")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 2)
	assert.Equal(t, c.Profiles["testCase2"].Region, "us")
	assert.Equal(t, c.Profiles["testCase2"].PersonalAPIKey, "apiKeyGoesHere")

	// Set the default profile to the new one
	err = c.SetDefaultProfile("testCase2")
	assert.NoError(t, err)
	assert.Equal(t, c.DefaultProfile, "testCase2")

	// Delete the initial profile
	err = c.RemoveProfile("testCase1")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 1)

	err = c.RemoveProfile("testCase1")
	assert.Error(t, err)
	assert.Equal(t, len(c.Profiles), 1)

	// Load the credentials again to verify json
	c2, err := LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, len(c2.Profiles), 1)
	assert.Equal(t, c2.DefaultProfile, "testCase2")
	assert.Equal(t, c2.ConfigDirectory, f)
	assert.False(t, c.profileExists("testCase1"))
}
