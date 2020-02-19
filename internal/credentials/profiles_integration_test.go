package credentials

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCredentials(t *testing.T) {
	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	// defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadCredentials(f)
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 0)

	// Create an initial profile to work with
	err = c.AddProfile("testCase1", "us", "apiKeyGoesHere", "anotherApiKeyGoesHere")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 1)

	// Set the default profile to the only one we've got
	err = c.SetDefaultProfile("testCase1")
	assert.NoError(t, err)
	assert.Equal(t, c.DefaultProfile, "testCase1")

	// Adding a profile with the same name should result in an error
	err = c.AddProfile("testCase1", "us", "foot", "hand")
	assert.Error(t, err)
	assert.Equal(t, len(c.Profiles), 1)

	// Create a second profile to work with
	err = c.AddProfile("testCase2", "us", "apiKeyGoesHere", "anotherApiKeyGoesHere")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 2)
	assert.Equal(t, c.Profiles["testCase2"].Region, "us")
	assert.Equal(t, c.Profiles["testCase2"].PersonalAPIKey, "apiKeyGoesHere")
	assert.Equal(t, c.Profiles["testCase2"].AdminAPIKey, "anotherApiKeyGoesHere")

	// Set the default profile to the new one
	err = c.SetDefaultProfile("testCase2")
	assert.NoError(t, err)
	assert.Equal(t, c.DefaultProfile, "testCase2")

	// Delete the initial profile
	err = c.RemoveProfile("testCase1")
	assert.NoError(t, err)
	assert.Equal(t, len(c.Profiles), 1)

}
