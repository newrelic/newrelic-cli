package entities

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"

	"github.com/stretchr/testify/assert"
)

func TestEntityChangetracker(t *testing.T) {
	command := cmdEntityChangetracker

	assert.Equal(t, "changetracker", command.Name())
}

func TestEntityChangetrackerCreate(t *testing.T) {
	assert.Equal(t, "create", cmdEntityChangetrackerCreate.Name())

	testcobra.CheckCobraMetadata(t, cmdEntityChangetrackerCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdEntityChangetrackerCreate,
		[]string{"guid", "version"})

}
