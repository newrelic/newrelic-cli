//go:build unit
// +build unit

package events

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestPostFile(t *testing.T) {
	assert.Equal(t, "postFile", cmdPostFile.Name())

	testcobra.CheckCobraMetadata(t, cmdPostFile)
	testcobra.CheckCobraRequiredFlags(t, cmdPostFile, []string{})
}
