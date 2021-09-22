//go:build unit
// +build unit

package nerdgraph

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestQuery(t *testing.T) {
	assert.Equal(t, "query", cmdQuery.Name())

	testcobra.CheckCobraMetadata(t, cmdQuery)
	testcobra.CheckCobraRequiredFlags(t, cmdQuery, []string{})
}
