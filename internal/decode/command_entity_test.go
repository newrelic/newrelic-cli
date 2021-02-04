// +build unit

package decode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestDecodeEntityCommand(t *testing.T) {
	assert.Equal(t, "decode", Command.Name())

	testcobra.CheckCobraMetadata(t, cmdEntity)
	testcobra.CheckCobraRequiredFlags(t, cmdEntity, []string{})
	testcobra.CheckCobraCommandAliases(t, cmdEntity, []string{})
}
