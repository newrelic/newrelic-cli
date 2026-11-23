//go:build unit

package synthetics

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/stretchr/testify/assert"
)

func TestCreateCertCheckMonitor(t *testing.T) {
	command := cmdCreateCertCheckMonitor

	assert.Equal(t, "create-cert-check-monitor", command.Name())
	assert.True(t, command.HasFlags())
	testcobra.CheckCobraMetadata(t, command)
	testcobra.CheckCobraRequiredFlags(t, command, []string{"name", "domain", "numberDaysToFailBeforeCertExpires", "period", "status"})
}
