//go:build unit

package servicelevel

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestServiceLevelDelete(t *testing.T) {
	command := CmdServiceLevelDelete

	assert.Equal(t, "delete", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"guid"})
}
