// +build unit

package profiles

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

func TestProfilesAdd(t *testing.T) {
	assert.Equal(t, "add", cmdAdd.Name())

	testcobra.CheckCobraMetadata(t, cmdAdd)
	testcobra.CheckCobraRequiredFlags(t, cmdAdd, []string{"name", "region"})
	testcobra.CheckCobraCommandAliases(t, cmdAdd, []string{})
}

func TestProfilesDefault(t *testing.T) {
	assert.Equal(t, "default", cmdDefault.Name())

	testcobra.CheckCobraMetadata(t, cmdDefault)
	testcobra.CheckCobraRequiredFlags(t, cmdDefault, []string{"name"})
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
	testcobra.CheckCobraRequiredFlags(t, cmdDelete, []string{"name"})
	testcobra.CheckCobraCommandAliases(t, cmdDelete, []string{"remove", "rm"}) // DEPRECATED: from nr1 cli
}
