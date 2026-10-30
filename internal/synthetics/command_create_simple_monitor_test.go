//go:build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateSimpleMonitor(t *testing.T) {
	command := CmdCreateSimpleMonitor

	assert.Equal(t, "create-simple-monitor", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"name", "period", "status", "uri"})
}
