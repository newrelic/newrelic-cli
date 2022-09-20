//go:build unit
// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestYq(t *testing.T) {
	assert.Equal(t, "yq", cmdYq.Name())

	testcobra.CheckCobraMetadata(t, cmdYq)
}
