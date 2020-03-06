// +build unit

package credentials

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCredentialsCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())

	for _, c := range Command.Commands() {
		assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
		assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
		assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
		assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())
	}
}

func TestCredentialsAdd(t *testing.T) {
	command := credentialsAdd
	assert.Equal(t, "add", command.Name())

	requiredFlags := []string{"profileName", "region"}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestCredentialsDefault(t *testing.T) {
	command := credentialsDefault
	assert.Equal(t, "default", command.Name())

	requiredFlags := []string{"profileName"}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestCredentialsList(t *testing.T) {
	command := credentialsList
	assert.Equal(t, "list", command.Name())

	requiredFlags := []string{}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestCredentialsRemove(t *testing.T) {
	command := credentialsRemove
	assert.Equal(t, "remove", command.Name())

	requiredFlags := []string{"profileName"}

	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}
