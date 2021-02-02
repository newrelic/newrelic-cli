// +build unit

package main

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestDecodeCommand(t *testing.T) {
	assert.Equal(t, "newrelic-dev", Command.Name())

	testcobra.CheckCobraMetadata(t, cmdDecode)
	testcobra.CheckCobraRequiredFlags(t, cmdDecode, []string{})
	testcobra.CheckCobraCommandAliases(t, cmdDecode, []string{})
}
