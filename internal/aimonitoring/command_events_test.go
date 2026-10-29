//go:build unit

package aimonitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestAIMonitoringEvents(t *testing.T) {
	assert.Equal(t, "events", cmdEvents.Name())

	testcobra.CheckCobraMetadata(t, cmdEvents)
	testcobra.CheckCobraRequiredFlags(t, cmdEvents, []string{})
}

func TestSortedEventTypeAliasKeys(t *testing.T) {
	keys := sortedEventTypeAliasKeys()

	assert.Equal(t, []string{"agent", "embedding", "feedback", "message", "summary", "tool", "vectorsearch", "vectorsearchresult"}, keys)
}
