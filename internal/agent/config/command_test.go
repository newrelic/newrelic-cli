// +build unit

package apm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())

	for _, c := range Command.Commands() {
		assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
		assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
		assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
		assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())
	}
}
