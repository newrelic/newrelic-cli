// +build integration

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/credentials"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
	mock "github.com/newrelic/newrelic-client-go/pkg/testhelpers"
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
	assert.NotEmpty(t, c.Profiles[defaultProfileName].LicenseKey)
	assert.NotEmpty(t, c.Profiles[defaultProfileName].InsightsInsertKey)

	// Ensure that we don't Fatal out if the default profile already exists, but
	// was not specified in the default-profile.json.
	if err = os.Remove(fmt.Sprintf("%s/%s.json", f, credentials.DefaultProfileFile)); err != nil {
		t.Fatal(err)
	}

	initializeProfile()
}

func TestFetchLicenseKey_missingKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(missingLicenseKey))
	})

	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)

	nr := &newrelic.NewRelic{
		NerdGraph: nerdgraph.New(tc),
	}

	response, err := fetchLicenseKey(nr, 0)
	require.Error(t, err)
	require.Empty(t, response)
}

func TestFetchLicenseKey_withKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(withLicenseKey))
	})

	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)

	nr := &newrelic.NewRelic{
		NerdGraph: nerdgraph.New(tc),
	}

	response, err := fetchLicenseKey(nr, 0)
	require.NoError(t, err)
	require.Equal(t, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", response)
}

var missingLicenseKey string = `
{
  "data": {
    "actor": {
      "account": {
      }
    }
  }
}
`

var withLicenseKey string = `
{
  "data": {
    "actor": {
      "account": {
        "licenseKey": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
      }
    }
  }
}
`
