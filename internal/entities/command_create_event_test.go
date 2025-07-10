package entities

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestEntityCreateEventCommand(t *testing.T) {
	assert.Equal(t, "create-event", CmdEntityCreateEvent.Name())
	testcobra.CheckCobraMetadata(t, CmdEntityCreateEvent)
	testcobra.CheckCobraRequiredFlags(t, CmdEntityCreateEvent,
		[]string{"entityName", "category", "type"})
}
