//go:build unit

package nrql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestMacrosSubcommand(t *testing.T) {
	assert.Equal(t, "macros", cmdMacros.Name())
	testcobra.CheckCobraMetadata(t, cmdMacros)
}

func TestMacrosCreateCommand(t *testing.T) {
	assert.Equal(t, "create", cmdMacrosCreate.Name())
	testcobra.CheckCobraMetadata(t, cmdMacrosCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdMacrosCreate, []string{"name", "expression"})
}

func TestMacrosGetCommand(t *testing.T) {
	assert.Equal(t, "get", cmdMacrosGet.Name())
	testcobra.CheckCobraMetadata(t, cmdMacrosGet)
	testcobra.CheckCobraRequiredFlags(t, cmdMacrosGet, []string{"id"})
}

func TestMacrosListCommand(t *testing.T) {
	assert.Equal(t, "list", cmdMacrosList.Name())
	testcobra.CheckCobraMetadata(t, cmdMacrosList)
}

func TestMacrosUpdateCommand(t *testing.T) {
	assert.Equal(t, "update", cmdMacrosUpdate.Name())
	testcobra.CheckCobraMetadata(t, cmdMacrosUpdate)
	testcobra.CheckCobraRequiredFlags(t, cmdMacrosUpdate, []string{"id"})
}

func TestMacrosDeleteCommand(t *testing.T) {
	assert.Equal(t, "delete", cmdMacrosDelete.Name())
	testcobra.CheckCobraMetadata(t, cmdMacrosDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdMacrosDelete, []string{"id"})
}
