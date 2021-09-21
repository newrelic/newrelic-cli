//go:build unit
// +build unit

package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestProfilesCommand(t *testing.T) {
	assert.Equal(t, "profile", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
	testcobra.CheckCobraCommandAliases(t, Command, []string{"profiles"}) // DEPRECATED: from nr1 cli
}

func TestCredentialsAdd(t *testing.T) {
	assert.Equal(t, "add", cmdAdd.Name())

	testcobra.CheckCobraMetadata(t, cmdAdd)
	testcobra.CheckCobraCommandAliases(t, cmdAdd, []string{"configure"})
}

func TestProfilesDefault(t *testing.T) {
	assert.Equal(t, "default", cmdDefault.Name())

	testcobra.CheckCobraMetadata(t, cmdDefault)
	testcobra.CheckCobraCommandAliases(t, cmdDefault, []string{})
}

func TestProfilesList(t *testing.T) {
	assert.Equal(t, "list", cmdList.Name())

	testcobra.CheckCobraMetadata(t, cmdList)
	testcobra.CheckCobraRequiredFlags(t, cmdList, []string{})
	testcobra.CheckCobraCommandAliases(t, cmdList, []string{"ls"}) // DEPRECATED: from nr1 cli
}

func TestProfilesDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdDelete)
	testcobra.CheckCobraCommandAliases(t, cmdDelete, []string{"remove", "rm"}) // DEPRECATED: from nr1 cli
}
