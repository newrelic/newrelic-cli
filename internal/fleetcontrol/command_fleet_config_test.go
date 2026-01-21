//go:build unit
// +build unit

package fleetcontrol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

// TestFleetGetConfigurationCommand tests the fleet get-configuration command
func TestFleetGetConfigurationCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "get-configuration", cmdFleetGetConfiguration.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetGetConfiguration)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetGetConfiguration, []string{"entity-guid"})
}

// TestFleetGetConfigurationCommandFlags tests that get-configuration command has all expected flags
func TestFleetGetConfigurationCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"entity-guid"}
	optionalFlags := []string{"organization-id", "mode", "version"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetGetConfiguration.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetGetConfiguration.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetAddVersionCommand tests the fleet add-version command
func TestFleetAddVersionCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "add-version", cmdFleetAddVersion.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetAddVersion)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetAddVersion, []string{"configuration-guid", "body"})
}

// TestFleetAddVersionCommandFlags tests that add-version command has all expected flags
func TestFleetAddVersionCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"configuration-guid", "body"}
	optionalFlags := []string{"organization-id"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetAddVersion.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetAddVersion.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetDeleteConfigurationCommand tests the fleet delete-configuration command
func TestFleetDeleteConfigurationCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "delete-configuration", cmdFleetDeleteConfiguration.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetDeleteConfiguration)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetDeleteConfiguration, []string{"configuration-guid"})
}

// TestFleetDeleteConfigurationCommandFlags tests that delete-configuration command has all expected flags
func TestFleetDeleteConfigurationCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"configuration-guid"}
	optionalFlags := []string{"organization-id"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetDeleteConfiguration.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetDeleteConfiguration.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetDeleteVersionCommand tests the fleet delete-version command
func TestFleetDeleteVersionCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "delete-version", cmdFleetDeleteVersion.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetDeleteVersion)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetDeleteVersion, []string{"version-guid"})
}

// TestFleetDeleteVersionCommandFlags tests that delete-version command has all expected flags
func TestFleetDeleteVersionCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"version-guid"}
	optionalFlags := []string{"organization-id"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetDeleteVersion.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetDeleteVersion.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}
