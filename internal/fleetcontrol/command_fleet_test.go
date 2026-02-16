//go:build unit
// +build unit

package fleetcontrol

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// TestParseTags tests the tag parsing functionality
func TestParseTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []string
		expected    []fleetcontrol.FleetControlTagInput
		expectError bool
	}{
		{
			name:  "single tag with single value",
			input: []string{"env:production"},
			expected: []fleetcontrol.FleetControlTagInput{
				{Key: "env", Values: []string{"production"}},
			},
			expectError: false,
		},
		{
			name:  "single tag with multiple values",
			input: []string{"env:production,staging"},
			expected: []fleetcontrol.FleetControlTagInput{
				{Key: "env", Values: []string{"production", "staging"}},
			},
			expectError: false,
		},
		{
			name:  "multiple tags",
			input: []string{"env:production", "team:platform,devops"},
			expected: []fleetcontrol.FleetControlTagInput{
				{Key: "env", Values: []string{"production"}},
				{Key: "team", Values: []string{"platform", "devops"}},
			},
			expectError: false,
		},
		{
			name:  "tag with whitespace",
			input: []string{"env: production , staging "},
			expected: []fleetcontrol.FleetControlTagInput{
				{Key: "env", Values: []string{"production", "staging"}},
			},
			expectError: false,
		},
		{
			name:        "invalid format - no colon",
			input:       []string{"envproduction"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid format - empty key",
			input:       []string{":production"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid format - empty value",
			input:       []string{"env:"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty input",
			input:       []string{},
			expected:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := ParseTags(tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestFleetCreateCommand tests the fleet create command
func TestFleetCreateCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "create", cmdFleetCreate.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetCreate, []string{"name", "managed-entity-type"})
}

// TestFleetCreateCommandFlags tests that create command has all expected flags
func TestFleetCreateCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"name", "managed-entity-type"}
	optionalFlags := []string{"description", "product", "organization-id", "operating-system", "tags"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetCreate.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetCreate.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetUpdateCommand tests the fleet update command
func TestFleetUpdateCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "update", cmdFleetUpdate.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetUpdate)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetUpdate, []string{"id"})
}

// TestFleetUpdateCommandFlags tests that update command has all expected flags
func TestFleetUpdateCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"id"}
	optionalFlags := []string{"name", "description", "tags"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetUpdate.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetUpdate.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetDeleteCommand tests the fleet delete command
func TestFleetDeleteCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "delete", cmdFleetDelete.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetDelete, []string{"id"})
}

// TestFleetDeleteCommandFlags tests that delete command has the expected flag
func TestFleetDeleteCommandFlags(t *testing.T) {
	t.Parallel()

	flag := cmdFleetDelete.Flags().Lookup("id")
	require.NotNil(t, flag, "id flag should exist")
}

// TestFleetAddMembersCommand tests the fleet add-members command
func TestFleetAddMembersCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "add-members", cmdFleetAddMembers.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetAddMembers)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetAddMembers, []string{"fleet-id", "ring", "entity-ids"})
}

// TestFleetAddMembersCommandFlags tests that add-members command has all expected flags
func TestFleetAddMembersCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"fleet-id", "ring", "entity-ids"}

	for _, flagName := range requiredFlags {
		t.Run("has_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetAddMembers.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetRemoveMembersCommand tests the fleet remove-members command
func TestFleetRemoveMembersCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "remove-members", cmdFleetRemoveMembers.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetRemoveMembers)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetRemoveMembers, []string{"fleet-id", "ring", "entity-ids"})
}

// TestFleetRemoveMembersCommandFlags tests that remove-members command has all expected flags
func TestFleetRemoveMembersCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"fleet-id", "ring", "entity-ids"}

	for _, flagName := range requiredFlags {
		t.Run("has_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetRemoveMembers.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestFleetCommandStructure tests that all expected subcommands exist under fleet command
func TestFleetCommandStructure(t *testing.T) {
	t.Parallel()

	expectedCommands := []string{
		"create",
		"update",
		"delete",
		"add-members",
		"remove-members",
		"create-configuration",
		"get-configuration",
		"add-version",
		"delete-configuration",
		"delete-version",
	}

	for _, cmdName := range expectedCommands {
		t.Run("has_"+cmdName+"_command", func(t *testing.T) {
			cmd, _, err := cmdFleet.Find([]string{cmdName})
			require.NoError(t, err)
			assert.Equal(t, cmdName, cmd.Name())
		})
	}
}

// TestFleetCreateConfigurationCommand tests the fleet create-configuration command
func TestFleetCreateConfigurationCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "create-configuration", cmdFleetCreateConfiguration.Name())
	testcobra.CheckCobraMetadata(t, cmdFleetCreateConfiguration)
	testcobra.CheckCobraRequiredFlags(t, cmdFleetCreateConfiguration, []string{"entity-name", "agent-type", "managed-entity-type", "body"})
}

// TestFleetCreateConfigurationCommandFlags tests that create-configuration command has all expected flags
func TestFleetCreateConfigurationCommandFlags(t *testing.T) {
	t.Parallel()

	requiredFlags := []string{"entity-name", "agent-type", "managed-entity-type", "body"}
	optionalFlags := []string{"organization-id"}

	// Test required flags exist
	for _, flagName := range requiredFlags {
		t.Run("has_required_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetCreateConfiguration.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}

	// Test optional flags exist
	for _, flagName := range optionalFlags {
		t.Run("has_optional_flag_"+flagName, func(t *testing.T) {
			flag := cmdFleetCreateConfiguration.Flags().Lookup(flagName)
			require.NotNil(t, flag, "flag %s should exist", flagName)
		})
	}
}

// TestParseAttributes tests the attribute parsing functionality
func TestParseAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []string
		expected    map[string]string
		expectError bool
	}{
		{
			name:  "single attribute",
			input: []string{"attribute1:value1"},
			expected: map[string]string{
				"attribute1": "value1",
			},
			expectError: false,
		},
		{
			name:  "multiple attributes",
			input: []string{"attribute1:value1", "attribute2:value2"},
			expected: map[string]string{
				"attribute1": "value1",
				"attribute2": "value2",
			},
			expectError: false,
		},
		{
			name:  "attribute with whitespace",
			input: []string{"attribute1: value1 "},
			expected: map[string]string{
				"attribute1": "value1",
			},
			expectError: false,
		},
		{
			name:  "attribute with special characters",
			input: []string{"env.config:production-v1.2"},
			expected: map[string]string{
				"env.config": "production-v1.2",
			},
			expectError: false,
		},
		{
			name:        "invalid format - no colon",
			input:       []string{"attribute1value1"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid format - empty key",
			input:       []string{":value1"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "invalid format - empty value",
			input:       []string{"attribute1:"},
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty input",
			input:       []string{},
			expected:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := parseAttributes(tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
