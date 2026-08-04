package uninstall

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnsupportedPackageManager_AllOperationsReturnClearError(t *testing.T) {
	pm := &UnsupportedPackageManager{}

	require.ErrorContains(t, pm.StopService("newrelic-infra"), "not supported")
	require.ErrorContains(t, pm.DisableService("newrelic-infra"), "not supported")
	require.ErrorContains(t, pm.RemovePackage("newrelic-infra"), "not supported")
	require.ErrorContains(t, pm.RemovePath("/etc/newrelic-infra"), "not supported")
}
