package install

import (
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
func TestCommandValiProfile(t *testing.T) {

	accountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	region := os.Getenv("NEW_RELIC_REGION")

	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
	err := validateProfile(5)
	assert.Error(t, err)
	assert.Equal(t, types.EventTypes.AccountIDMissing, err.EventName)

	if accountID == "" {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", "67890")
	} else {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", accountID)
	}

	os.Setenv("NEW_RELIC_API_KEY", "")
	err = validateProfile(5)
	assert.Equal(t, types.EventTypes.APIKeyMissing, err.EventName)
	if apiKey == "" {
		os.Setenv("NEW_RELIC_API_KEY", "67890")
	} else {
		os.Setenv("NEW_RELIC_API_KEY", apiKey)
	}

	os.Setenv("NEW_RELIC_REGION", "")
	err = validateProfile(5)
	assert.Equal(t, types.EventTypes.RegionMissing, err.EventName)

	os.Setenv("NEW_RELIC_ACCOUNT_ID", accountID)
	os.Setenv("NEW_RELIC_REGION", region)
	os.Setenv("NEW_RELIC_API_KEY", apiKey)
}
