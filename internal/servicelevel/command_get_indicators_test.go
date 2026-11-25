//go:build unit

package servicelevel

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestServiceLevelGetIndicators(t *testing.T) {
	command := CmdServiceLevelGetIndicators

	assert.Equal(t, "get-indicators", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"entityGuid"})
}
