// +build unit

package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestAgentConfig(t *testing.T) {
	assert.Equal(t, "config", cmdConfig.Name())

	testcobra.CheckCobraMetadata(t, cmdConfig)
}

func TestAgentConfigObfuscate(t *testing.T) {
	assert.Equal(t, "obfuscate", cmdConfigObfuscate.Name())
	testcobra.CheckCobraMetadata(t, cmdConfigObfuscate)
	testcobra.CheckCobraRequiredFlags(t, cmdConfigObfuscate, []string{})
}

func TestInstallCommand(t *testing.T) {
	assert.Equal(t, "integrations", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
}
