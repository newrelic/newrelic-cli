// +build unit

package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgentCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())
}
