// +build unit

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "version", cmdVersion.Name())

	testcobra.CheckCobraMetadata(t, cmdVersion)
	testcobra.CheckCobraRequiredFlags(t, cmdVersion, []string{})
	testcobra.CheckCobraCommandAliases(t, cmdVersion, []string{})
}
