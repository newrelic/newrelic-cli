//go:build unit

package servicelevel

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestServiceLevelGet(t *testing.T) {
	command := CmdServiceLevelGet

	assert.Equal(t, "get", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"entityGuid"})
}
