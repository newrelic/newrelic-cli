//go:build unit

package changeTracking

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestChangeTrackingCreateDeployment(t *testing.T) {
	command := CmdChangeTrackingCreateDeployment

	assert.Equal(t, "deployment", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"entityGuid", "version", "validationFlags"})
}
