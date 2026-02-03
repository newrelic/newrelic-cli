package changeTracking

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestChangeTrackingCommand(t *testing.T) {
	assert.Equal(t, "changeTracking", Command.Name())
	testcobra.CheckCobraMetadata(t, Command)
}

func TestChangeTrackingCreateCommand(t *testing.T) {
	assert.Equal(t, "create", CmdChangeTrackingCreate.Name())
	testcobra.CheckCobraMetadata(t, CmdChangeTrackingCreate)
	testcobra.CheckCobraRequiredFlags(t, CmdChangeTrackingCreate,
		[]string{"entitySearch", "category", "type"})
}
