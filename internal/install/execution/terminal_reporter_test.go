// +build unit

package execution

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestTerminalStatusReporter_interface(t *testing.T) {
	var r StatusSubscriber = NewTerminalStatusReporter()
	require.NotNil(t, r)
}

func TestGenerateEntityLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockSuccessLinkGenerator()
	r.successLinkGenerator = g

	status := &InstallStatus{}

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, g.GenerateEntityLinkCallCount)
	require.Equal(t, 0, g.GenerateExplorerLinkCallCount)
}

func TestGenerateExplorerLink(t *testing.T) {
	r := NewTerminalStatusReporter()
	g := NewMockSuccessLinkGenerator()
	r.successLinkGenerator = g

	status := &InstallStatus{}
	status.successLinkConfig = types.SuccessLinkConfig{
		Type:   "explorer",
		Filter: "\"`tags.language` = 'java'\"",
	}

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 0, g.GenerateEntityLinkCallCount)
	require.Equal(t, 1, g.GenerateExplorerLinkCallCount)
}
