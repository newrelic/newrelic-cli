//go:build unit

package servicelevel

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	command := Command

	assert.Equal(t, "servicelevel", command.Name())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{})
}
