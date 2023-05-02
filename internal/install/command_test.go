//go:build unit
// +build unit

package install

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestInstallCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "install", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{})
}
