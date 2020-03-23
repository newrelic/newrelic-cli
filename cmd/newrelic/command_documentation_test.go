// +build unit

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestDocumentation(t *testing.T) {
	assert.Equal(t, "documentation", cmdDocumentation.Name())

	testcobra.CheckCobraMetadata(t, cmdDocumentation)
	testcobra.CheckCobraRequiredFlags(t, cmdDocumentation, []string{"outputDir"})
}
