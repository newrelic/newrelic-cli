//go:build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateBrokenLinksMonitor(t *testing.T) {
	command := cmdCreateBrokenLinksMonitor

	assert.Equal(t, "create-broken-links-monitor", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"name", "uri", "period", "status"})
}
