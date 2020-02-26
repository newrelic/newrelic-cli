package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	assert.Equal(t, "completion", completionCmd.Name())

	c := completionCmd

	assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
	assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
	assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
	assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())

	requiredFlags := []string{"shell"}

	for _, r := range requiredFlags {
		x := completionCmd.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestVersion(t *testing.T) {
	assert.Equal(t, "version", versionCmd.Name())

	c := versionCmd

	assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
	assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
	assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
	assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())
}
