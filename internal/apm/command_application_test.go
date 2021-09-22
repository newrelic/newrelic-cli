//go:build unit
// +build unit

package apm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestApmApp(t *testing.T) {
	assert.Equal(t, "application", cmdApp.Name())

	testcobra.CheckCobraMetadata(t, cmdApp)
	testcobra.CheckCobraRequiredFlags(t, cmdApp, []string{})
}

func TestApmAppGet(t *testing.T) {
	assert.Equal(t, "get", cmdAppGet.Name())

	testcobra.CheckCobraMetadata(t, cmdAppGet)
	// guid is required, but Persisted Flags are not supported by this check
	testcobra.CheckCobraRequiredFlags(t, cmdAppGet, []string{})
}

func TestApmAppSearch(t *testing.T) {
	assert.Equal(t, "search", cmdAppSearch.Name())

	testcobra.CheckCobraMetadata(t, cmdAppSearch)
	testcobra.CheckCobraRequiredFlags(t, cmdAppSearch, []string{})
}
