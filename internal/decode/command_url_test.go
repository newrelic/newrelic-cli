// +build unit

package decode

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestDecodeURLCommand(t *testing.T) {
	assert.Equal(t, "url", Command.Name())

	testcobra.CheckCobraMetadata(t, cmdDecode)
	testcobra.CheckCobraRequiredFlags(t, cmdDecode, []string{})
	testcobra.CheckCobraCommandAliases(t, cmdDecode, []string{})
}
