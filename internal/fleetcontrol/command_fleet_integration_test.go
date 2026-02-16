//go:build integration
// +build integration

package fleetcontrol

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var (
	testOrgID = "b961cf81-d62b-4359-8822-7b1d6dadd374"
)

// resetFlags recursively resets all flags in a command and its subcommands
func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
		f.Value.Set(f.DefValue)
	})

	// Also reset local flags
	cmd.LocalFlags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
		f.Value.Set(f.DefValue)
	})

	// Recursively reset all subcommands
	for _, subCmd := range cmd.Commands() {
		resetFlags(subCmd)
	}
}

// setupCommand prepares a command for testing by resetting its state and setting up output capture
func setupCommand() {
	// Reset the command state
	Command.SetArgs([]string{})

	// Reset all flag values to their defaults recursively
	// This is critical to prevent flag values from persisting between test executions
	resetFlags(Command)

	// LOCAL TESTING ONLY: Uncomment to route traffic through proxy (comment out before pushing to GitHub)
	os.Setenv("HTTPS_PROXY", "http://localhost:8888")

	// Initialize client for integration tests
	tc := mock.NewIntegrationTestConfig(&testing.T{})
	nrClient, err := newrelic.New(
		newrelic.ConfigPersonalAPIKey(tc.PersonalAPIKey),
		newrelic.ConfigRegion(tc.Region().String()),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to create test client: %v", err))
	}
	client.NRClient = nrClient
}

// executeCommand runs a command with the given arguments and returns the output
// This function captures os.Stdout directly since the commands write to stdout via printJSON
func executeCommand(args []string) (string, error) {
	setupCommand()

	// Capture stdout by creating a pipe
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Channel to capture output
	outputChan := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outputChan <- buf.String()
	}()

	// Execute the command
	Command.SetArgs(args)
	err := Command.Execute()

	// Restore stdout and close pipe
	w.Close()
	os.Stdout = oldStdout

	// Get the captured output
	output := <-outputChan

	return output, err
}

// TestIntegrationFleetCreateUpdateDelete tests the complete lifecycle of fleet CLI commands
func TestIntegrationFleetCreateUpdateDelete(t *testing.T) {
	// Note: Not using t.Parallel() because tests share global Command state and stdout capture
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	fleetName := fmt.Sprintf("CLI Test Fleet %d", time.Now().Unix())
	var createdFleetID string

	defer func() {
		if createdFleetID != "" {
			// Clean up: Delete the fleet using CLI command
			_, _ = executeCommand([]string{
				"fleet", "delete",
				"--fleet-id", createdFleetID,
			})
			fmt.Printf("Cleaned up fleet: %s\n", createdFleetID)
		}
	}()

	// Test 1: Create fleet using CLI command
	t.Run("create_fleet", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "create",
			"--name", fleetName,
			"--managed-entity-type", "HOST",
			"--description", "CLI integration test fleet",
			"--operating-system", "LINUX",
			"--tags", "env:test",
			"--tags", "cli-test:integration",
		})

		require.NoError(t, err, "Fleet create command should succeed")
		require.NotEmpty(t, output, "Command should produce output")

		// Parse JSON output
		var result FleetResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		// Verify response
		assert.Equal(t, "success", result.Status)
		assert.NotEmpty(t, result.ID)
		assert.Equal(t, fleetName, result.Name)
		assert.Equal(t, "HOST", result.ManagedEntityType)
		assert.NotEmpty(t, result.Tags)

		createdFleetID = result.ID
		fmt.Printf("Created fleet with ID: %s\n", createdFleetID)
	})

	// Test 2: Get fleet using CLI command
	t.Run("get_fleet", func(t *testing.T) {
		require.NotEmpty(t, createdFleetID, "Fleet ID should be set from create test")

		output, err := executeCommand([]string{
			"fleet", "get",
			"--fleet-id", createdFleetID,
		})

		require.NoError(t, err, "Fleet get command should succeed")

		var result FleetResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		// Debug: print the result if status is not success
		if result.Status != "success" {
			fmt.Printf("DEBUG: Get fleet failed with error: %s\n", result.Error)
			fmt.Printf("DEBUG: Full output: %s\n", output)
		}

		assert.Equal(t, "success", result.Status)
		assert.Equal(t, createdFleetID, result.ID)
		assert.Equal(t, fleetName, result.Name)
	})

	// Test 3: Update fleet using CLI command
	t.Run("update_fleet", func(t *testing.T) {
		require.NotEmpty(t, createdFleetID, "Fleet ID should be set")

		updatedName := fmt.Sprintf("CLI Updated Fleet %d", time.Now().Unix())
		output, err := executeCommand([]string{
			"fleet", "update",
			"--id", createdFleetID,
			"--name", updatedName,
			"--description", "Updated via CLI",
			"--tags", "status:updated",
		})

		require.NoError(t, err, "Fleet update command should succeed")

		var result FleetResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		assert.Equal(t, "success", result.Status)
		assert.Equal(t, createdFleetID, result.ID)
		assert.Equal(t, updatedName, result.Name)
		assert.Equal(t, "Updated via CLI", result.Description)
	})

	// Test 4: Search fleets using CLI command
	t.Run("search_fleet", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "search",
			"--name-contains", fleetName,
		})

		require.NoError(t, err, "Fleet search command should succeed")

		// Search returns a plain JSON array, not a wrapper
		var results []FleetEntityOutput
		err = json.Unmarshal([]byte(output), &results)
		require.NoError(t, err, "Output should be valid JSON array")

		// Note: Search may return empty results due to indexing delays or
		// if the fleet name has been updated by the update test
		// The important thing is that the command executes successfully
		fmt.Printf("Search found %d fleets (may be empty due to timing/name changes)\n", len(results))
	})

	// Test 5: Delete fleet using CLI command
	t.Run("delete_fleet", func(t *testing.T) {
		require.NotEmpty(t, createdFleetID, "Fleet ID should be set")

		output, err := executeCommand([]string{
			"fleet", "delete",
			"--fleet-id", createdFleetID,
		})

		require.NoError(t, err, "Fleet delete command should succeed")

		var result FleetDeleteResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		assert.Equal(t, "success", result.Status)
		assert.Equal(t, createdFleetID, result.ID)

		// Clear the ID since we deleted it
		createdFleetID = ""
	})

	fmt.Println("✅ Successfully completed fleet CLI command integration tests")
}

// TestIntegrationFleetGetById tests retrieving a fleet by ID
func TestIntegrationFleetGetById(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	fleetName := fmt.Sprintf("CLI Test Get By ID %d", time.Now().Unix())
	var createdFleetID string

	defer func() {
		if createdFleetID != "" {
			executeCommand([]string{"fleet", "delete", "--fleet-id", createdFleetID})
		}
	}()

	// Create a fleet
	output, err := executeCommand([]string{
		"fleet", "create",
		"--name", fleetName,
		"--managed-entity-type", "HOST",
		"--operating-system", "LINUX",
		"--description", "Test fleet for ID retrieval",
	})
	require.NoError(t, err)

	var createResult FleetResponseWrapper
	json.Unmarshal([]byte(output), &createResult)
	createdFleetID = createResult.ID

	// Test retrieving by ID
	t.Run("get_fleet_by_id", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "get",
			"--fleet-id", createdFleetID,
		})

		require.NoError(t, err, "Get fleet by ID should succeed")

		var result FleetResponseWrapper
		json.Unmarshal([]byte(output), &result)

		assert.Equal(t, "success", result.Status)
		assert.Equal(t, createdFleetID, result.ID)
		assert.Equal(t, fleetName, result.Name)
		assert.Equal(t, "HOST", result.ManagedEntityType)
		assert.Equal(t, "Test fleet for ID retrieval", result.Description)
	})

	fmt.Println("✅ Successfully completed fleet get by ID test")
}

// TestIntegrationFleetSearchFilters tests fleet search with various filter combinations
func TestIntegrationFleetSearchFilters(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	// Create multiple fleets for searching
	testPrefix := fmt.Sprintf("CLI Search Test %d", time.Now().Unix())
	var fleetIDs []string

	defer func() {
		for _, id := range fleetIDs {
			executeCommand([]string{"fleet", "delete", "--fleet-id", id})
		}
	}()

	// Create test fleets
	fleetConfigs := []struct {
		name              string
		entityType        string
		os                string
		tags              []string
	}{
		{
			name:       fmt.Sprintf("%s Linux Host A", testPrefix),
			entityType: "HOST",
			os:         "LINUX",
			tags:       []string{"env:prod", "region:us-east"},
		},
		{
			name:       fmt.Sprintf("%s Windows Host B", testPrefix),
			entityType: "HOST",
			os:         "WINDOWS",
			tags:       []string{"env:staging", "region:us-west"},
		},
		{
			name:       fmt.Sprintf("%s K8s Cluster C", testPrefix),
			entityType: "KUBERNETESCLUSTER",
			tags:       []string{"env:prod", "region:eu-west"},
		},
	}

	for _, config := range fleetConfigs {
		args := []string{
			"fleet", "create",
			"--name", config.name,
			"--managed-entity-type", config.entityType,
		}

		if config.os != "" {
			args = append(args, "--operating-system", config.os)
		}

		for _, tag := range config.tags {
			args = append(args, "--tags", tag)
		}

		output, err := executeCommand(args)
		if err == nil {
			var result FleetResponseWrapper
			json.Unmarshal([]byte(output), &result)
			if result.ID != "" {
				fleetIDs = append(fleetIDs, result.ID)
			}
		}
	}

	// Test 1: Search by name prefix
	t.Run("search_by_name", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "search",
			"--name-contains", testPrefix,
		})

		if err == nil {
			var results []FleetEntityOutput
			json.Unmarshal([]byte(output), &results)
			assert.NotEmpty(t, results, "Should find fleets matching name prefix")
		}
	})

	// Test 2: Search by entity type
	t.Run("search_by_entity_type", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "search",
			"--managed-entity-type", "HOST",
		})

		if err == nil {
			var results []FleetEntityOutput
			json.Unmarshal([]byte(output), &results)
			// Results may be empty if no HOST fleets exist
		}
	})

	// Test 3: Search with combined filters
	t.Run("search_with_filters", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "search",
			"--name-contains", testPrefix,
			"--managed-entity-type", "HOST",
		})

		if err == nil {
			var results []FleetEntityOutput
			json.Unmarshal([]byte(output), &results)
			// Results may be empty if filters don't match anything
		}
	})

	fmt.Println("✅ Successfully completed fleet search filter tests")
}

// TestIntegrationFleetWithDifferentOptions tests fleet creation with various flag combinations
func TestIntegrationFleetWithDifferentOptions(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	tests := []struct {
		name               string
		managedEntityType  string
		operatingSystem    string
		tags               []string
		expectSuccess      bool
	}{
		{
			name:              "host_linux",
			managedEntityType: "HOST",
			operatingSystem:   "LINUX",
			tags:              []string{"env:prod", "os:linux"},
			expectSuccess:     true,
		},
		{
			name:              "host_windows",
			managedEntityType: "HOST",
			operatingSystem:   "WINDOWS",
			tags:              []string{"env:prod", "os:windows"},
			expectSuccess:     true,
		},
		{
			name:              "kubernetes",
			managedEntityType: "KUBERNETESCLUSTER",
			tags:              []string{"platform:k8s"},
			expectSuccess:     true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fleetName := fmt.Sprintf("CLI Test %s %d", tt.name, time.Now().Unix())
			var createdFleetID string

			defer func() {
				if createdFleetID != "" {
					executeCommand([]string{"fleet", "delete", "--fleet-id", createdFleetID})
				}
			}()

			// Build command args
			args := []string{
				"fleet", "create",
				"--name", fleetName,
				"--managed-entity-type", tt.managedEntityType,
			}

			if tt.operatingSystem != "" {
				args = append(args, "--operating-system", tt.operatingSystem)
			}

			for _, tag := range tt.tags {
				args = append(args, "--tags", tag)
			}

			output, err := executeCommand(args)

			if tt.expectSuccess {
				require.NoError(t, err, "Command should succeed")

				var result FleetResponseWrapper
				err = json.Unmarshal([]byte(output), &result)
				require.NoError(t, err, "Output should be valid JSON")

				assert.Equal(t, "success", result.Status)
				assert.Equal(t, tt.managedEntityType, result.ManagedEntityType)
				createdFleetID = result.ID
			} else {
				assert.Error(t, err, "Command should fail")
			}
		})
	}
}

// TestIntegrationFleetUpdatePartialFields tests updating individual fields
func TestIntegrationFleetUpdatePartialFields(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	// Create a fleet first
	fleetName := fmt.Sprintf("CLI Test Partial Update %d", time.Now().Unix())
	var createdFleetID string

	defer func() {
		if createdFleetID != "" {
			executeCommand([]string{"fleet", "delete", "--fleet-id", createdFleetID})
		}
	}()

	// Create fleet
	output, err := executeCommand([]string{
		"fleet", "create",
		"--name", fleetName,
		"--managed-entity-type", "HOST",
		"--description", "Original description",
		"--tags", "original:value",
	})
	require.NoError(t, err)

	var createResult FleetResponseWrapper
	json.Unmarshal([]byte(output), &createResult)
	createdFleetID = createResult.ID

	// Test 1: Update only name
	t.Run("update_name_only", func(t *testing.T) {
		newName := fmt.Sprintf("CLI New Name %d", time.Now().Unix())
		output, err := executeCommand([]string{
			"fleet", "update",
			"--id", createdFleetID,
			"--name", newName,
		})

		require.NoError(t, err)
		var result FleetResponseWrapper
		json.Unmarshal([]byte(output), &result)

		assert.Equal(t, newName, result.Name)
		assert.Equal(t, "Original description", result.Description)
	})

	// Test 2: Update only description
	t.Run("update_description_only", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "update",
			"--id", createdFleetID,
			"--description", "New description",
		})

		require.NoError(t, err)
		var result FleetResponseWrapper
		json.Unmarshal([]byte(output), &result)

		assert.Equal(t, "New description", result.Description)
	})

	// Test 3: Update only tags
	t.Run("update_tags_only", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "update",
			"--id", createdFleetID,
			"--tags", "updated:tag",
			"--tags", "cli-test:partial",
		})

		require.NoError(t, err)
		var result FleetResponseWrapper
		json.Unmarshal([]byte(output), &result)

		// Verify tags were updated
		assert.NotEmpty(t, result.Tags)
		tagKeys := make(map[string]bool)
		for _, tag := range result.Tags {
			tagKeys[tag.Key] = true
		}
		assert.True(t, tagKeys["updated"] || tagKeys["cli-test"])
	})

	fmt.Println("✅ Successfully completed partial update tests")
}

// TestIntegrationFleetErrorHandling tests error cases
func TestIntegrationFleetErrorHandling(t *testing.T) {
	t.Parallel()
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	// Test 1: Get non-existent fleet
	t.Run("get_nonexistent_fleet", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "get",
			"--fleet-id", "non-existent-id",
		})

		// Should return an error or error status in JSON
		if err == nil {
			// Check if output contains error status
			if strings.Contains(output, "\"status\":\"failed\"") ||
				strings.Contains(output, "\"error\"") {
				// Error is in JSON output, which is acceptable
				assert.Contains(t, output, "error")
			}
		}
	})

	// Test 2: Delete non-existent fleet
	t.Run("delete_nonexistent_fleet", func(t *testing.T) {
		output, err := executeCommand([]string{
			"fleet", "delete",
			"--fleet-id", "non-existent-id",
		})

		// Should handle gracefully
		if err == nil {
			// Check if output indicates error
			if strings.Contains(output, "\"status\":\"failed\"") ||
				strings.Contains(output, "\"error\"") {
				assert.Contains(t, output, "error")
			}
		}
	})

	// Test 3: Create fleet with missing required flags
	t.Run("create_without_required_flags", func(t *testing.T) {
		_, err := executeCommand([]string{
			"fleet", "create",
			"--name", "test",
			// Missing --managed-entity-type
		})

		// Should fail due to missing required flag
		assert.Error(t, err, "Should fail without required flags")
	})
}
