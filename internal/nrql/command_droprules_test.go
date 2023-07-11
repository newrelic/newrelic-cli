//go:build unit
// +build unit

package nrql

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestDropRulesQuery(t *testing.T) {
	assert.Equal(t, "droprules", cmdDropRules.Name())
	testcobra.CheckCobraMetadata(t, cmdDropRules)
}
