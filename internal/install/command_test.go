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
	// TEST 1: POSITIVE CASE - License key provided and matches the account
	// ==================================================================================
	// Validates that when a license key is provided and BELONGS to the configured
	// account, the NEW validateLicenseKeyForAccount() logic accepts it.
	// ==================================================================================
	t.Run("LicenseKeyProvided_MatchingAccount", func(t *testing.T) {
		// Skip if no real credentials available
		if origAccountID == "" || origAPIKey == "" || origRegion == "" || origLicenseKey == "" {
			t.Skip("Skipping test: requires NEW_RELIC_ACCOUNT_ID, NEW_RELIC_API_KEY, NEW_RELIC_REGION, and NEW_RELIC_LICENSE_KEY environment variables")
		}

		// Use the real credentials - license key should match the account
		os.Setenv("NEW_RELIC_ACCOUNT_ID", origAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Setenv("NEW_RELIC_LICENSE_KEY", origLicenseKey)

		// Call fetchLicenseKey - it will validate the license key against the account
		err := fetchLicenseKey()

		// Assert: Should succeed - license key belongs to the account
		if err != nil {
			t.Errorf("Expected no error when license key matches account")
			t.Logf("  Account ID: %s", origAccountID)
			t.Logf("  API Key: %s", utils.Obfuscate(origAPIKey))
			t.Logf("  Region: %s", origRegion)
			t.Logf("  License Key: %s", utils.Obfuscate(origLicenseKey))
			t.Logf("  Error: %v", err)
		}

		actual := os.Getenv("NEW_RELIC_LICENSE_KEY")
		if actual != origLicenseKey {
			t.Errorf("License key should remain unchanged. Expected: %s, Got: %s",
				utils.Obfuscate(origLicenseKey), utils.Obfuscate(actual))
		}
	})

	// ==================================================================================
	// TEST 2: NEGATIVE CASE - License key provided but does NOT match the account
	// ==================================================================================
	// Validates that when a license key does NOT belong to the configured account,
	// the NEW validateLicenseKeyForAccount() logic rejects it with CredentialAccountMismatch.
	// ==================================================================================
	t.Run("LicenseKeyProvided_NonMatchingAccount", func(t *testing.T) {
		// Skip if no real credentials available
		if origAccountID == "" || origAPIKey == "" || origRegion == "" || origLicenseKey == "" {
			t.Skip("Skipping test: requires NEW_RELIC_ACCOUNT_ID, NEW_RELIC_API_KEY, NEW_RELIC_REGION, and NEW_RELIC_LICENSE_KEY environment variables")
		}

		// Use a different account ID that won't match the license key
		wrongAccountID := "9999999"

		// Set up environment with WRONG account ID
		os.Setenv("NEW_RELIC_ACCOUNT_ID", wrongAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Setenv("NEW_RELIC_LICENSE_KEY", origLicenseKey)

		// Call fetchLicenseKey - it will validate the license key against the WRONG account
		err := fetchLicenseKey()

		// Assert: Should fail with CredentialAccountMismatch error
		if err == nil {
			t.Errorf("Expected CredentialAccountMismatch error when license key doesn't match account")
			t.Logf("  Correct Account ID: %s", origAccountID)
			t.Logf("  Wrong Account ID (used): %s", wrongAccountID)
			t.Logf("  API Key: %s", utils.Obfuscate(origAPIKey))
			t.Logf("  Region: %s", origRegion)
			t.Logf("  License Key: %s (belongs to %s, NOT %s)", utils.Obfuscate(origLicenseKey), origAccountID, wrongAccountID)
		}

		if err != nil && err.EventName != types.EventTypes.CredentialAccountMismatch {
			t.Errorf("Expected CredentialAccountMismatch error type, got: %v", err.EventName)
			t.Logf("  Wrong Account ID: %s", wrongAccountID)
			t.Logf("  License Key: %s", utils.Obfuscate(origLicenseKey))
			t.Logf("  Error: %v", err)
		}

		if err != nil && !strings.Contains(err.Error(), "credential mismatch detected") {
			t.Errorf("Error message should mention 'credential mismatch detected'")
			t.Logf("  Error message: %s", err.Error())
		}

		if err != nil && !strings.Contains(err.Error(), wrongAccountID) {
			t.Errorf("Error message should contain wrong account ID: %s", wrongAccountID)
			t.Logf("  Error message: %s", err.Error())
		}
	})

	// ==================================================================================
	// TEST 3: FALLBACK CASE - No license key provided, fetch from API
	// ==================================================================================
	// Validates the 3-tier fallback: env var → profile → API fetch.
	// When no license key is provided, it should fetch from the API for the account.
	// ==================================================================================
	t.Run("NoLicenseKeyProvided_FetchesFromAPI", func(t *testing.T) {
		// Skip if no real credentials available
		if origAccountID == "" || origAPIKey == "" || origRegion == "" {
			t.Skip("Skipping test: requires NEW_RELIC_ACCOUNT_ID, NEW_RELIC_API_KEY, and NEW_RELIC_REGION environment variables")
		}

		// Set up environment WITHOUT license key
		os.Setenv("NEW_RELIC_ACCOUNT_ID", origAccountID)
		os.Setenv("NEW_RELIC_API_KEY", origAPIKey)
		os.Setenv("NEW_RELIC_REGION", origRegion)
		os.Unsetenv("NEW_RELIC_LICENSE_KEY") // Unset to simulate not being provided

		// Call fetchLicenseKey - it will fetch license key from API
		err := fetchLicenseKey()

		// Assert: Should succeed - fetches license key from API
		if err != nil {
			t.Errorf("Expected no error when fetching license key from API")
			t.Logf("  Account ID: %s", origAccountID)
			t.Logf("  API Key: %s", utils.Obfuscate(origAPIKey))
			t.Logf("  Region: %s", origRegion)
			t.Logf("  Error: %v", err)
		}

		// Assert: License key should now be set in environment
		fetchedKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
		if fetchedKey == "" {
			t.Errorf("License key should be fetched and set in environment")
			t.Logf("  Account ID: %s", origAccountID)
			t.Logf("  API Key: %s", utils.Obfuscate(origAPIKey))
		}

		if len(fetchedKey) != 40 {
			t.Errorf("License key should be 40 characters, got: %d", len(fetchedKey))
			t.Logf("  Fetched License Key: %s", utils.Obfuscate(fetchedKey))
		}

		if !strings.HasSuffix(fetchedKey, "NRAL") {
			t.Errorf("License key should end with NRAL")
			t.Logf("  Fetched License Key: %s", utils.Obfuscate(fetchedKey))
		}
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
