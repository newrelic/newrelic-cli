//go:build unit

package workload

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestWorkloadsGetCollection(t *testing.T) {
	command := CmdWorkloadsGetCollection

	assert.Equal(t, "get-collection", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"accountId", "guid"})
}
