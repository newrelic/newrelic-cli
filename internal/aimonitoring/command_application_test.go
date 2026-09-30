//go:build unit

package aimonitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestAIMonitoringApplication(t *testing.T) {
	assert.Equal(t, "application", cmdApplication.Name())

	testcobra.CheckCobraMetadata(t, cmdApplication)
	testcobra.CheckCobraRequiredFlags(t, cmdApplication, []string{})
}

func TestAIMonitoringApplicationGet(t *testing.T) {
	assert.Equal(t, "get", cmdApplicationGet.Name())

	testcobra.CheckCobraMetadata(t, cmdApplicationGet)
	// guid is required, but Persisted Flags are not supported by this check
	testcobra.CheckCobraRequiredFlags(t, cmdApplicationGet, []string{})
}

func TestAIMonitoringApplicationSearch(t *testing.T) {
	assert.Equal(t, "search", cmdApplicationSearch.Name())

	testcobra.CheckCobraMetadata(t, cmdApplicationSearch)
	testcobra.CheckCobraRequiredFlags(t, cmdApplicationSearch, []string{})
}

func TestExtractEntityGUIDs(t *testing.T) {
	result := extractEntityGUIDs(nil)
	assert.Nil(t, result)

	result = extractEntityGUIDs([]nrdb.NRDBResult{
		{"guids": []interface{}{"MTIzfEFQTXxBUFBMSUNBVElPTnwxMjM", "", 42}},
	})
	assert.Equal(t, []common.EntityGUID{"MTIzfEFQTXxBUFBMSUNBVElPTnwxMjM"}, result)

	result = extractEntityGUIDs([]nrdb.NRDBResult{
		{"other": "value"},
	})
	assert.Nil(t, result)
}
