//go:build unit

package install

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func TestInstallCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "install", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
}

func TestValidateProfile(t *testing.T) {
	accountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	region := os.Getenv("NEW_RELIC_REGION")

	server := initSegmentMockServer()
	defer server.Close()

	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
	err := validateProfile()
	assert.Error(t, err)
	assert.Equal(t, types.EventTypes.AccountIDMissing, err.EventName)

	if accountID == "" {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", "67890")
	} else {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", accountID)
	}

	os.Setenv("NEW_RELIC_API_KEY", "")
	err = validateProfile()
	assert.Equal(t, types.EventTypes.APIKeyMissing, err.EventName)

	os.Setenv("NEW_RELIC_API_KEY", "67890")
	err = validateProfile()
	assert.Equal(t, types.EventTypes.InvalidUserAPIKeyFormat, err.EventName)

	if apiKey == "" {
		os.Setenv("NEW_RELIC_API_KEY", "NRAK-67890")
	} else {
		os.Setenv("NEW_RELIC_API_KEY", apiKey)
	}

	os.Setenv("NEW_RELIC_REGION", "au")
	err = validateProfile()
	assert.Equal(t, types.EventTypes.InvalidRegion, err.EventName)

	os.Setenv("NEW_RELIC_REGION", "")
	err = validateProfile()
	assert.Equal(t, types.EventTypes.RegionMissing, err.EventName)

	os.Setenv("NEW_RELIC_ACCOUNT_ID", accountID)
	os.Setenv("NEW_RELIC_REGION", region)
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
}

func TestFetchLicenseKey(t *testing.T) {
	// TODO: Error case
	// TODO: From API (mock?)
	// TODO: From profile

	// Save original environment to restore later
	origAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	origAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	origRegion := os.Getenv("NEW_RELIC_REGION")
	origLicenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")

	// Restore environment after test
	defer func() {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", origAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Setenv("NEW_RELIC_LICENSE_KEY", origLicenseKey)
	}()

	// ==================================================================================
	// TEST 1: License key matches account - NEW validation accepts it
	// ==================================================================================
	t.Run("LicenseKeyProvided_MatchingAccount", func(t *testing.T) {
		if origAccountID == "" || origAPIKey == "" || origRegion == "" || origLicenseKey == "" {
			t.Skip("Skipping: requires NEW_RELIC_ACCOUNT_ID, NEW_RELIC_API_KEY, NEW_RELIC_REGION, NEW_RELIC_LICENSE_KEY")
		}

		os.Setenv("NEW_RELIC_ACCOUNT_ID", origAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Setenv("NEW_RELIC_LICENSE_KEY", origLicenseKey)

		err := fetchLicenseKey()

		assert.Nil(t, err, "License key belonging to account should be accepted. Account: %s, Key: %s",
			origAccountID, utils.Obfuscate(origLicenseKey))
		assert.Equal(t, origLicenseKey, os.Getenv("NEW_RELIC_LICENSE_KEY"))
	})

	// ==================================================================================
	// TEST 2: License key doesn't match account - NEW validation rejects it
	// ==================================================================================
	t.Run("LicenseKeyProvided_NonMatchingAccount", func(t *testing.T) {
		if origAccountID == "" || origAPIKey == "" || origRegion == "" || origLicenseKey == "" {
			t.Skip("Skipping: requires NEW_RELIC_ACCOUNT_ID, NEW_RELIC_API_KEY, NEW_RELIC_REGION, NEW_RELIC_LICENSE_KEY")
		}

		wrongAccountID := "9999999"
		os.Setenv("NEW_RELIC_ACCOUNT_ID", wrongAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Setenv("NEW_RELIC_LICENSE_KEY", origLicenseKey)

		err := fetchLicenseKey()

		assert.NotNil(t, err, "License key not belonging to account should be rejected")
		assert.Equal(t, types.EventTypes.CredentialAccountMismatch, err.EventName,
			"Expected CredentialAccountMismatch. Account: %s, Key: %s", wrongAccountID, utils.Obfuscate(origLicenseKey))
		assert.Contains(t, err.Error(), "credential mismatch detected")
		assert.Contains(t, err.Error(), wrongAccountID)
	})

	// ==================================================================================
	// TEST 3: No license key provided - API fetch should succeed
	// ==================================================================================
	t.Run("NoLicenseKeyProvided_FetchesFromAPI", func(t *testing.T) {
		if origAccountID == "" || origAPIKey == "" || origRegion == "" {
			t.Skip("Skipping: requires NEW_RELIC_ACCOUNT_ID, NEW_RELIC_API_KEY, NEW_RELIC_REGION")
		}

		os.Setenv("NEW_RELIC_ACCOUNT_ID", origAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Unsetenv("NEW_RELIC_LICENSE_KEY")

		err := fetchLicenseKey()

		assert.Nil(t, err, "API fetch should succeed. Account: %s, API Key: %s",
			origAccountID, utils.Obfuscate(origAPIKey))
		assert.NotEmpty(t, os.Getenv("NEW_RELIC_LICENSE_KEY"), "License key should be fetched and set")
	})
}

func initSegmentMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`[]`))
	}))
	return server
}

func TestProxyNetwork(t *testing.T) {

	proxyConfig := struct {
		HTTPSProxy string
		HTTPProxy  string
	}{
		HTTPSProxy: "http://localhost:3128",
		HTTPProxy:  "http://localhost:8080",
	}

	// Validate HTTPSProxy
	if strings.HasPrefix(proxyConfig.HTTPSProxy, "http://") {
		t.Log("New Relic CLI exclusively supports https proxy, not http for security reasons.")
	} else if strings.HasPrefix(proxyConfig.HTTPSProxy, "https://") {
		// Do nothing
	} else {
		t.Log("Invalid proxy provided")
	}

	// Validate HTTPProxy
	if strings.HasPrefix(proxyConfig.HTTPProxy, "http://") {
		t.Log("If you need to use a proxy, consider setting the HTTPS_PROXY environment variable, then try again. New Relic CLI exclusively supports https proxy.")
	}

}
