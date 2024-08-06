//go:build unit

package install

import (
	"net/http"
	"net/http/httptest"
	"os"
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
	okCase := func() *types.DetailError { return nil }()
	// Error case
	// err := fetchLicenseKey()
	// assert.Equal(t, types.EventTypes.UnableToFetchLicenseKey, err.EventName)

	// licenseKey := os.Getenv("NEW_RELIC_LICENSE_KEY")
	// assert.Equal(t, "", licenseKey)

	// TODO: From API (mock?)

	// TODO: From profile

	// From environment variable
	expect := "0123456789abcdefABCDEF0123456789abcdNRAL"
	os.Setenv("NEW_RELIC_LICENSE_KEY", expect)

	err := fetchLicenseKey()
	assert.Equal(t, okCase, err)

	actual := os.Getenv("NEW_RELIC_LICENSE_KEY")
	assert.Equal(t, expect, actual)
}

func initSegmentMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`[]`))
	}))
	return server
}
