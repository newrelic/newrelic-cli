// +build unit

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "version", versionCmd.Name())

	c := versionCmd

	assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
	assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
	assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
	assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())
}
