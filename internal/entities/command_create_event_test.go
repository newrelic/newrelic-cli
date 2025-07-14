package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestEntityCreateEventCommand(t *testing.T) {
	assert.Equal(t, "create-event", CmdEntityCreateEvent.Name())
	testcobra.CheckCobraMetadata(t, CmdEntityCreateEvent)
	testcobra.CheckCobraRequiredFlags(t, CmdEntityCreateEvent,
		[]string{"entitySearch", "category", "type"})
}
