// +build unit

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestCompletion(t *testing.T) {
	assert.Equal(t, "completion", cmdCompletion.Name())

	testcobra.CheckCobraMetadata(t, cmdCompletion)
	testcobra.CheckCobraRequiredFlags(t, cmdCompletion, []string{"shell"})
	testcobra.CheckCobraCommandAliases(t, cmdCompletion, []string{})
}
