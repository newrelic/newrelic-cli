package testcobra

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// testCmdRequiredFlags is a helper function to make sure that
// the Cobra command has certain required flags set
func CheckCobraRequiredFlags(t *testing.T, command *cobra.Command, requiredFlags []string) {
	for _, r := range requiredFlags {
		x := command.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

// CheckCobraMetadata requires metadata fields be set on the command sent and
// any sub-command it contains.  Only checks one level deep.
func CheckCobraMetadata(t *testing.T, command *cobra.Command) {
	assert.NotEmptyf(t, command.Use, "Need to set Command.Use on Command %s", command.CalledAs())
	assert.NotEmptyf(t, command.Short, "Need to set Command.Short on Command %s", command.CalledAs())

	for _, c := range command.Commands() {
		assert.NotEmptyf(t, c.Use, "Need to set Command.Use on Command %s", c.CommandPath())
		assert.NotEmptyf(t, c.Short, "Need to set Command.Short on Command %s", c.CommandPath())
		assert.NotEmptyf(t, c.Long, "Need to set Command.Long on Command %s", c.CommandPath())
		assert.NotEmptyf(t, c.Example, "Need to set Command.Example on Command %s", c.CommandPath())
	}
}
