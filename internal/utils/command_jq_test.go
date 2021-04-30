// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestJq(t *testing.T) {
	assert.Equal(t, "jq", cmdJq.Name())

	testcobra.CheckCobraMetadata(t, cmdJq)
}
