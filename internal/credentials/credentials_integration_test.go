// +build integration

package credentials

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/pkg/region"
)

var overrideEnvVars = []string{
	"NEW_RELIC_API_KEY",
	"NEW_RELIC_REGION",
}

func TestApplyOverrides(t *testing.T) {
	// Do not run this in parallel, we are messing with the environment

	f, err := ioutil.TempDir("/tmp", "newrelic")
	assert.NoError(t, err)
	defer os.RemoveAll(f)

	// Initialize the new configuration directory
	c, err := LoadCredentials(f)
	assert.NoError(t, err)

	// Create an initial profile to work with
	testProfile := Profile{
		Region:            "us",
		APIKey:            "apiKeyGoesHere",
		InsightsInsertKey: "insightsInsertKeyGoesHere",
	}
	err = c.AddProfile("testCase1", testProfile)
	assert.NoError(t, err)
	p := c.Profiles["testCase1"]
	assert.NotNil(t, p)

	// Clear env vars we are going to mess with, and reset on exit
	for _, v := range overrideEnvVars {
		if val, ok := os.LookupEnv(v); ok {
			defer os.Setenv(v, val)
			os.Unsetenv(v)
		}
	}

	// Default case (no overrides)
	p2 := applyOverrides(&p)
	assert.NotNil(t, p2)
	assert.Equal(t, p.APIKey, p2.APIKey)
	assert.Equal(t, p.Region, p2.Region)

	// Override just the API Key
	os.Setenv("NEW_RELIC_API_KEY", "anotherAPIKey")
	p2 = applyOverrides(&p)
	assert.NotNil(t, p2)
	assert.Equal(t, "anotherAPIKey", p2.APIKey)
	assert.Equal(t, p.Region, p2.Region)

	// Both
	os.Setenv("NEW_RELIC_REGION", "US")
	p2 = applyOverrides(&p)
	assert.NotNil(t, p2)
	assert.Equal(t, "anotherAPIKey", p2.APIKey)
	assert.Equal(t, region.US.String(), p2.Region)

	// Override just the REGION (valid)
	os.Unsetenv("NEW_RELIC_API_KEY")
	p2 = applyOverrides(&p)
	assert.NotNil(t, p2)
	assert.Equal(t, p.APIKey, p2.APIKey)
	assert.Equal(t, region.US.String(), p2.Region)

	// Region lowercase
	os.Setenv("NEW_RELIC_REGION", "eu")
	p2 = applyOverrides(&p)
	assert.NotNil(t, p2)
	assert.Equal(t, p.APIKey, p2.APIKey)
	assert.Equal(t, region.EU.String(), p2.Region)
}
