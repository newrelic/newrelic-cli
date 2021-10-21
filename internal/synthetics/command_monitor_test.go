//go:build unit
// +build unit

package synthetics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestSyntheticsMonitor(t *testing.T) {
	assert.Equal(t, "monitor", cmdMon.Name())

	testcobra.CheckCobraMetadata(t, cmdMon)
	testcobra.CheckCobraRequiredFlags(t, cmdMon, []string{})
}

func TestSyntheticsMonitorGet(t *testing.T) {
	assert.Equal(t, "get", cmdMonGet.Name())

	testcobra.CheckCobraMetadata(t, cmdMonGet)
	testcobra.CheckCobraRequiredFlags(t, cmdMonGet, []string{})
}

func TestSyntheticsMonitorSearch(t *testing.T) {
	assert.Equal(t, "search", cmdMonSearch.Name())

	testcobra.CheckCobraMetadata(t, cmdMonSearch)
	testcobra.CheckCobraRequiredFlags(t, cmdMonSearch, []string{})
}

func TestSyntheticsMonitorList(t *testing.T) {
	assert.Equal(t, "list", cmdMonList.Name())

	testcobra.CheckCobraMetadata(t, cmdMonList)
	testcobra.CheckCobraRequiredFlags(t, cmdMonList, []string{})
}
