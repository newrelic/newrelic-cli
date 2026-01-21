//go:build unit

package nrql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestDropRulesQuery(t *testing.T) {
	assert.Equal(t, "droprules", cmdDropRules.Name())
	testcobra.CheckCobraMetadata(t, cmdDropRules)
}
