package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestEntityChangeEventCommand(t *testing.T) {
	assert.Equal(t, "change-event", CmdEntityChangeEvent.Name())
	testcobra.CheckCobraMetadata(t, CmdEntityChangeEvent)
}

func TestEntityChangeEventCreateCommand(t *testing.T) {
	assert.Equal(t, "create", CmdEntityChangeEventCreate.Name())
	testcobra.CheckCobraMetadata(t, CmdEntityChangeEventCreate)
	testcobra.CheckCobraRequiredFlags(t, CmdEntityChangeEventCreate,
		[]string{"entitySearch", "category", "type"})
}
