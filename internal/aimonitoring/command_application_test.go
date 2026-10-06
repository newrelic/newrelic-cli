//go:build unit

package aimonitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestAIMonitoringApplication(t *testing.T) {
	assert.Equal(t, "application", cmdApplication.Name())

	testcobra.CheckCobraMetadata(t, cmdApplication)
	testcobra.CheckCobraRequiredFlags(t, cmdApplication, []string{})
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

func TestFilterEntitiesByTags(t *testing.T) {
	enabled := &entities.ApmApplicationEntity{
		Tags: []entities.EntityTag{
			{Key: "aiEnabledApp", Values: []string{"true"}},
			{Key: "env", Values: []string{"prod", "staging"}},
		},
	}
	ents := []entities.EntityInterface{enabled}

	// No tags given: no filtering applied.
	result := filterEntitiesByTags(ents, nil)
	assert.Equal(t, ents, result)

	// Single matching tag.
	result = filterEntitiesByTags(ents, []map[string]string{{"key": "aiEnabledApp", "value": "true"}})
	assert.Equal(t, []entities.EntityInterface{enabled}, result)

	// Tag value matched among multiple values.
	result = filterEntitiesByTags(ents, []map[string]string{{"key": "env", "value": "staging"}})
	assert.Equal(t, []entities.EntityInterface{enabled}, result)

	// Multiple tags must all match (AND semantics).
	result = filterEntitiesByTags(ents, []map[string]string{
		{"key": "aiEnabledApp", "value": "true"},
		{"key": "env", "value": "prod"},
	})
	assert.Equal(t, []entities.EntityInterface{enabled}, result)

	// If any one of multiple tags doesn't match, the entity is excluded.
	result = filterEntitiesByTags(ents, []map[string]string{
		{"key": "aiEnabledApp", "value": "true"},
		{"key": "env", "value": "qa"},
	})
	assert.Equal(t, []entities.EntityInterface{}, result)

	// No entity matches.
	result = filterEntitiesByTags(ents, []map[string]string{{"key": "aiEnabledApp", "value": "maybe"}})
	assert.Equal(t, []entities.EntityInterface{}, result)
}

func TestToApplicationSearchResults(t *testing.T) {
	result := toApplicationSearchResults(nil)
	assert.Equal(t, []applicationSearchResult{}, result)

	app := &entities.ApmApplicationEntity{
		Name: "checkout-service",
		GUID: "MTIzfEFQTXxBUFBMSUNBVElPTnwxMjM",
	}

	result = toApplicationSearchResults([]entities.EntityInterface{app})
	assert.Equal(t, []applicationSearchResult{
		{Name: "checkout-service", GUID: "MTIzfEFQTXxBUFBMSUNBVElPTnwxMjM"},
	}, result)
}
