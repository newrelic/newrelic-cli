// +build integration

package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/region"
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
	err = c.AddProfile("testCase1", "us", "apiKeyGoesHere", "insightsInsertKeyGoesHere", 0)
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 1)
	assert.Equal(t, region.US.String(), c.Profiles["testCase1"].Region)
	assert.Equal(t, "apiKeyGoesHere", c.Profiles["testCase1"].APIKey)
	assert.Equal(t, "insightsInsertKeyGoesHere", c.Profiles["testCase1"].InsightsInsertKey)
	assert.Equal(t, "", c.DefaultProfile)

	// Set the default profile to the only one we've got
	err = c.SetDefaultProfile("testCase1")
	assert.NoError(t, err)
	assert.Equal(t, c.DefaultProfile, "testCase1")

	// Adding a profile with the same name should result in an error
	err = c.AddProfile("testCase1", "us", "foot", "", 0)
	assert.Error(t, err)
	assert.Equal(t, len(c.Profiles), 1)
	assert.True(t, c.profileExists("testCase1"))

	// Create a second profile to work with
	err = c.AddProfile("testCase2", "us", "apiKeyGoesHere", "", 0)
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 2)
	assert.Equal(t, c.Profiles["testCase2"].Region, region.US.String())
	assert.Equal(t, c.Profiles["testCase2"].APIKey, "apiKeyGoesHere")

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

	// Remove the default profile and check the results
	_, err = os.Stat(fmt.Sprintf("%s/%s.json", f, "default-profile"))
	assert.NoError(t, err)

	err = c.RemoveProfile("testCase2")
	assert.NoError(t, err)
	assert.Equal(t, c.DefaultProfile, "")
	_, err = os.Stat(fmt.Sprintf("%s/%s.json", f, "default-profile"))
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestCredentialLowerCaseRegion(t *testing.T) {
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
	err = c.AddProfile("testCase1", "US", "apiKeyGoesHere", "insightsInsertKeyGoesHere", 0)
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 1)
	assert.Equal(t, region.US.String(), c.Profiles["testCase1"].Region)
	assert.Equal(t, "apiKeyGoesHere", c.Profiles["testCase1"].APIKey)
	assert.Equal(t, "insightsInsertKeyGoesHere", c.Profiles["testCase1"].InsightsInsertKey)
}

// TestCredentialCompatibilityNR1
func TestCredentialCompatibilityNR1(t *testing.T) {
	t.Parallel()

	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Custom struct to mirror the config of NR1, and bypass
	// any custom marshal / unmarshal code we have
	testCredentialData := map[string]struct {
		APIKey string
		Region string
	}{
		"test": {
			APIKey: "apiKeyGoesHere",
			Region: "us",
		},
		"testeu": {
			APIKey: "apiKeyEU",
			Region: "EU",
		},
	}
	file, jsonErr := json.MarshalIndent(testCredentialData, "", " ")
	assert.NoError(t, jsonErr)

	err = ioutil.WriteFile(f+"/credentials.json", file, 0600)
	assert.NoError(t, err)

	c, loadErr := LoadCredentials(f)
	assert.NoError(t, loadErr)
	assert.Equal(t, len(testCredentialData), len(c.Profiles))

	for k := range c.Profiles {
		assert.Equal(t, testCredentialData[k].APIKey, c.Profiles[k].APIKey)
		assert.Equal(t, testCredentialData[k].Region, c.Profiles[k].Region)
	}
}
