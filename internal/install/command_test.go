package install

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
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

	err := validateProfile(5)
	assert.Error(t, err)
	assert.Equal(t, types.EventTypes.AccountIDMissing, err.EventName)

	os.Setenv("NEW_RELIC_ACCOUNT_ID", "67890")
	os.Setenv("NEW_RELIC_API_KEY", "")
	APIKey := configAPI.GetActiveProfileString(config.APIKey)
	log.Info(APIKey)

	err = validateProfile(5)
	assert.Equal(t, types.EventTypes.APIKeyMissing, err.EventName)

	os.Setenv("NEW_RELIC_API_KEY", "12345")
	os.Setenv("NEW_RELIC_REGION", "")
	err = validateProfile(5)
	assert.Equal(t, types.EventTypes.RegionMissing, err.EventName)

	os.Setenv("NEW_RELIC_API_KEY", "")
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "")
}
