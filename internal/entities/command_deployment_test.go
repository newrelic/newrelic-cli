package entities

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"

	"github.com/stretchr/testify/assert"
)

func TestEntityDeployment(t *testing.T) {
	command := cmdEntityDeployment

	assert.Equal(t, "deployment", command.Name())
}

func TestEntityDeploymentCreate(t *testing.T) {
	assert.Equal(t, "create", cmdEntityDeploymentCreate.Name())

	testcobra.CheckCobraMetadata(t, cmdEntityDeploymentCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdEntityDeploymentCreate,
		[]string{"guid", "version"})
}
