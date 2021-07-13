// +build unit

package command

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestConfigCommand(t *testing.T) {
	assert.Equal(t, "config", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
}

func TestCmdReset(t *testing.T) {
	assert.Equal(t, "reset", cmdReset.Name())

	testcobra.CheckCobraMetadata(t, cmdReset)
	testcobra.CheckCobraRequiredFlags(t, cmdReset, []string{"key"})
}

func TestCmdGet(t *testing.T) {
	assert.Equal(t, "get", cmdGet.Name())

	testcobra.CheckCobraMetadata(t, cmdGet)
	testcobra.CheckCobraRequiredFlags(t, cmdGet, []string{"key"})
}

func TestCmdSet(t *testing.T) {
	assert.Equal(t, "set", cmdSet.Name())

	testcobra.CheckCobraMetadata(t, cmdSet)
	testcobra.CheckCobraRequiredFlags(t, cmdSet, []string{"key", "value"})
}
