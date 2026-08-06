//go:build unit

package uninstall

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestUninstallCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "uninstall", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{"recipe"})
}
