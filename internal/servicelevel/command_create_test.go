//go:build unit

package servicelevel

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	command := CmdCreate

	assert.Equal(t, "create", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"entityGuid", "name", "accountId", "facets"})
}
