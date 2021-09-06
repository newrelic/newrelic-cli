// +build unit

package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestAgent(t *testing.T) {
	assert.Equal(t, "agent", Command.Name())
	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
}

func TestAgentConfig(t *testing.T) {
	assert.Equal(t, "config", cmdConfig.Name())
	testcobra.CheckCobraMetadata(t, cmdConfig)
	testcobra.CheckCobraMetadata(t, cmdConfig)
}

func TestAgentConfigObfuscate(t *testing.T) {
	assert.Equal(t, "obfuscate", cmdConfigObfuscate.Name())
	testcobra.CheckCobraMetadata(t, cmdConfigObfuscate)
	testcobra.CheckCobraRequiredFlags(t, cmdConfigObfuscate, []string{})
}

func TestAgentConfigMigrate(t *testing.T) {
	assert.Equal(t, "migrateV3toV4", cmdMigrateV3toV4.Name())
	testcobra.CheckCobraMetadata(t, cmdMigrateV3toV4)
	testcobra.CheckCobraRequiredFlags(t, cmdMigrateV3toV4, []string{})
}
