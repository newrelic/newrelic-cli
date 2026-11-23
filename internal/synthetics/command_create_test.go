//go:build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestSyntheticsCreate(t *testing.T) {
	command := CmdSyntheticsCreate

	assert.Equal(t, "create", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"accountId", "description", "name"})
}
