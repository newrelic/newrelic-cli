package changeevent

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestChangeTrackingCommand(t *testing.T) {
	command := CmdChangeTracking
	assert.Equal(t, "changetracking", command.Name())
}

func TestChangeTrackingCreateEventCommand(t *testing.T) {
	assert.Equal(t, "create-event", CmdChangeTrackingCreateEvent.Name())
	testcobra.CheckCobraMetadata(t, CmdChangeTrackingCreateEvent)
	testcobra.CheckCobraRequiredFlags(t, CmdChangeTrackingCreateEvent,
		[]string{"entityName", "category", "type"})
}
