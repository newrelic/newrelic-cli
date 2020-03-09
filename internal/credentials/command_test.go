// +build unit

package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestCredentialsCommand(t *testing.T) {
	assert.Equal(t, "profiles", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
}

func TestCredentialsAdd(t *testing.T) {
	assert.Equal(t, "add", cmdAdd.Name())

	testcobra.CheckCobraMetadata(t, cmdAdd)
	testcobra.CheckCobraRequiredFlags(t, cmdAdd, []string{"profileName", "region"})
}

func TestCredentialsDefault(t *testing.T) {
	assert.Equal(t, "default", cmdDefault.Name())

	testcobra.CheckCobraMetadata(t, cmdDefault)
	testcobra.CheckCobraRequiredFlags(t, cmdDefault, []string{"profileName"})
}

func TestCredentialsList(t *testing.T) {
	assert.Equal(t, "list", cmdList.Name())

	testcobra.CheckCobraMetadata(t, cmdList)
	testcobra.CheckCobraRequiredFlags(t, cmdList, []string{})
}

func TestCredentialsRemove(t *testing.T) {
	assert.Equal(t, "remove", cmdRemove.Name())

	testcobra.CheckCobraMetadata(t, cmdRemove)
	testcobra.CheckCobraRequiredFlags(t, cmdRemove, []string{"profileName"})
}
