//go:build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	command := cmdCreate

	assert.Equal(t, "create", command.Name())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{})
}
