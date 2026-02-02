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
	accountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	region := os.Getenv("NEW_RELIC_REGION")
	licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")

	// Set up required environment variables for validateProfile()
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	os.Setenv("NEW_RELIC_API_KEY", "NRAK-TEST1234567890")
	os.Setenv("NEW_RELIC_REGION", "US")

	okCase := func() *types.DetailError { return nil }()
	// TODO: Error case

	// TODO: From API (mock?)

	// TODO: From profile

	// From environment variable
	expect := "0123456789abcdefABCDEF0123456789abcdNRAL"
	os.Setenv("NEW_RELIC_LICENSE_KEY", expect)

	err := fetchLicenseKey()
	assert.Equal(t, okCase, err)

	actual := os.Getenv("NEW_RELIC_LICENSE_KEY")
	assert.Equal(t, expect, actual)

	os.Setenv("NEW_RELIC_ACCOUNT_ID", accountID)
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
	os.Setenv("NEW_RELIC_REGION", region)
	os.Setenv("NEW_RELIC_LICENSE_KEY", licenseKey)
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
