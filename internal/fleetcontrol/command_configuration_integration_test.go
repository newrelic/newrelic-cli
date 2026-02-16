//go:build integration
// +build integration

package fleetcontrol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock "github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// TestIntegrationConfigurationLifecycle tests the complete lifecycle of configuration CLI commands
// This test follows the same pattern as TestIntegrationFleetCreateUpdateDelete
func TestIntegrationConfigurationLifecycle(t *testing.T) {
	// Note: Not using t.Parallel() because tests share global Command state and stdout capture
	_, err := mock.GetTestAccountID()
	if err != nil {
		t.Skipf("Skipping integration test: %s", err)
	}

	configName := fmt.Sprintf("CLI Test Config %d", time.Now().Unix())
	var createdConfigID string

	defer func() {
		if createdConfigID != "" {
			executeCommand([]string{
				"configuration", "delete",
				"--configuration-id", createdConfigID,
			})
			fmt.Printf("Cleaned up configuration: %s\n", createdConfigID)
		}
	}()

	// Test 1: Create configuration using inline content (default case - no org ID)
	t.Run("create_configuration_inline", func(t *testing.T) {
		configContent := `{"log": {"level": "info"}, "metrics": {"enabled": true}}`

		output, err := executeCommand([]string{
			"configuration", "create",
			"--name", configName,
			"--agent-type", "NRInfra",
			"--managed-entity-type", "HOST",
			"--configuration-content", configContent,
		})

		require.NoError(t, err, "Configuration create command should succeed")

		var result ConfigurationResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		assert.Equal(t, "success", result.Status)
		assert.NotNil(t, result.Result)

		// Extract configuration ID from result
		if resultMap, ok := result.Result.(map[string]interface{}); ok {
			if entityGuid, ok := resultMap["entityGuid"].(string); ok {
				createdConfigID = entityGuid
			} else if configEntityGUID, ok := resultMap["configurationEntityGUID"].(string); ok {
				createdConfigID = configEntityGUID
			}
		}

		require.NotEmpty(t, createdConfigID, "Should have created configuration ID")
		fmt.Printf("Created configuration with ID: %s\n", createdConfigID)
	})

	// Test 2: Get configuration (default ConfigEntity mode)
	t.Run("get_configuration", func(t *testing.T) {
		require.NotEmpty(t, createdConfigID, "Configuration ID should be set")

		output, err := executeCommand([]string{
			"configuration", "get",
			"--configuration-id", createdConfigID,
		})

		require.NoError(t, err, "Configuration get command should succeed")
		assert.NotEmpty(t, output, "Should return configuration content")
	})

	// Test 3: Add a new version
	t.Run("add_version", func(t *testing.T) {
		require.NotEmpty(t, createdConfigID)

		output, err := executeCommand([]string{
			"configuration", "versions", "add",
			"--configuration-id", createdConfigID,
			"--configuration-content", `{"version": "2.0"}`,
		})

		require.NoError(t, err, "Add version command should succeed")

		var result ConfigurationResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err)

		assert.Equal(t, "success", result.Status)
	})

	// Test 4: List versions and extract version IDs
	var versionIDs []string
	t.Run("list_versions", func(t *testing.T) {
		require.NotEmpty(t, createdConfigID)

		output, err := executeCommand([]string{
			"configuration", "versions", "list",
			"--configuration-id", createdConfigID,
		})

		require.NoError(t, err, "List versions command should succeed")

		var result ConfigurationResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err)

		assert.Equal(t, "success", result.Status)
		assert.NotNil(t, result.Result)

		// Extract version IDs from the result
		if resultMap, ok := result.Result.(map[string]interface{}); ok {
			if versions, ok := resultMap["versions"].([]interface{}); ok {
				for _, v := range versions {
					if versionMap, ok := v.(map[string]interface{}); ok {
						// The field name is "entity_guid" not "id"
						if versionID, ok := versionMap["entity_guid"].(string); ok {
							versionIDs = append(versionIDs, versionID)
						}
					}
				}
			}
		}

		fmt.Printf("Found %d version(s) for configuration\n", len(versionIDs))
	})

	// Test 5: Delete extra versions (keep only one version for configuration deletion to succeed)
	t.Run("delete_extra_versions", func(t *testing.T) {
		if len(versionIDs) <= 1 {
			t.Skip("Only one version exists, no need to delete extra versions")
			return
		}

		// Delete all versions except the first one (keep one version)
		versionsToDelete := versionIDs[1:]
		fmt.Printf("Deleting %d extra version(s) to leave one version remaining\n", len(versionsToDelete))

		for i, versionID := range versionsToDelete {
			output, err := executeCommand([]string{
				"configuration", "versions", "delete",
				"--version-id", versionID,
			})

			require.NoError(t, err, "Delete version %d command should succeed", i+1)

			var result ConfigurationDeleteResponseWrapper
			err = json.Unmarshal([]byte(output), &result)
			require.NoError(t, err)

			assert.Equal(t, "success", result.Status)
			fmt.Printf("Deleted version %d of %d: %s\n", i+1, len(versionsToDelete), versionID)
		}
	})

	// Test 6: Delete configuration (should succeed now with only one version remaining)
	t.Run("delete_configuration", func(t *testing.T) {
		require.NotEmpty(t, createdConfigID, "Configuration ID should be set")

		output, err := executeCommand([]string{
			"configuration", "delete",
			"--configuration-id", createdConfigID,
		})

		require.NoError(t, err, "Configuration delete command should succeed")

		var result ConfigurationDeleteResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		assert.Equal(t, "success", result.Status)
		assert.Equal(t, createdConfigID, result.ID)

		createdConfigID = ""
		fmt.Println("✅ Successfully deleted configuration")
	})

	// Test 7: Create configuration from file
	t.Run("create_from_file", func(t *testing.T) {
		// Create a temporary config file
		configContent := `{
			"log": {
				"level": "debug"
			},
			"metrics": {
				"enabled": true,
				"interval": 60
			}
		}`

		tmpFile, err := ioutil.TempFile("", "test-config-*.json")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(configContent)
		require.NoError(t, err)
		tmpFile.Close()

		fileConfigName := fmt.Sprintf("CLI Test Config File %d", time.Now().Unix())

		output, err := executeCommand([]string{
			"configuration", "create",
			"--name", fileConfigName,
			"--agent-type", "NRInfra",
			"--managed-entity-type", "HOST",
			"--configuration-file-path", tmpFile.Name(),
		})

		require.NoError(t, err, "Configuration create from file should succeed")

		var result ConfigurationResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err, "Output should be valid JSON")

		assert.Equal(t, "success", result.Status)

		// Clean up the file-based configuration
		if resultMap, ok := result.Result.(map[string]interface{}); ok {
			var fileConfigID string
			if entityGuid, ok := resultMap["entityGuid"].(string); ok {
				fileConfigID = entityGuid
			} else if configEntityGUID, ok := resultMap["configurationEntityGUID"].(string); ok {
				fileConfigID = configEntityGUID
			}

			if fileConfigID != "" {
				executeCommand([]string{
					"configuration", "delete",
					"--configuration-id", fileConfigID,
				})
			}
		}

		fmt.Println("✅ Successfully tested configuration creation from file")
	})

	// Test 8: Create configuration with complex nested content
	t.Run("create_complex_content", func(t *testing.T) {
		complexConfig := `{
			"log": {
				"level": "debug",
				"file": "/var/log/newrelic/agent.log",
				"max_size": 100
			},
			"metrics": {
				"enabled": true,
				"interval": 60,
				"filters": {
					"include": ["cpu.*", "memory.*"],
					"exclude": ["temp.*"]
				}
			},
			"integrations": [
				{
					"name": "nginx",
					"enabled": true,
					"config": {
						"status_url": "http://localhost/status"
					}
				},
				{
					"name": "redis",
					"enabled": false
				}
			]
		}`

		complexConfigName := fmt.Sprintf("CLI Test Complex Config %d", time.Now().Unix())

		output, err := executeCommand([]string{
			"configuration", "create",
			"--name", complexConfigName,
			"--agent-type", "NRInfra",
			"--managed-entity-type", "HOST",
			"--configuration-content", complexConfig,
		})

		require.NoError(t, err, "Should create configuration with complex content")

		var result ConfigurationResponseWrapper
		err = json.Unmarshal([]byte(output), &result)
		require.NoError(t, err)

		assert.Equal(t, "success", result.Status)

		var complexConfigID string
		if resultMap, ok := result.Result.(map[string]interface{}); ok {
			if entityGuid, ok := resultMap["entityGuid"].(string); ok {
				complexConfigID = entityGuid
			} else if configEntityGUID, ok := resultMap["configurationEntityGUID"].(string); ok {
				complexConfigID = configEntityGUID
			}
		}

		require.NotEmpty(t, complexConfigID, "Should have created configuration ID")

		// Verify we can retrieve the complex configuration
		output, err = executeCommand([]string{
			"configuration", "get",
			"--configuration-id", complexConfigID,
		})

		require.NoError(t, err, "Should retrieve complex configuration")
		assert.NotEmpty(t, output)

		// Clean up
		executeCommand([]string{
			"configuration", "delete",
			"--configuration-id", complexConfigID,
		})

		fmt.Println("✅ Successfully completed complex configuration content test")
	})

	// Test 9: Test different agent types
	t.Run("test_different_agent_types", func(t *testing.T) {
		agentTypes := []string{"NRInfra", "NRDOT"}

		for _, agentType := range agentTypes {
			agentType := agentType
			t.Run(fmt.Sprintf("agent_type_%s", agentType), func(t *testing.T) {
				agentConfigName := fmt.Sprintf("CLI Test %s %d", agentType, time.Now().Unix())

				output, err := executeCommand([]string{
					"configuration", "create",
					"--name", agentConfigName,
					"--agent-type", agentType,
					"--managed-entity-type", "HOST",
					"--configuration-content", `{"test": true}`,
				})

				require.NoError(t, err, "Create config for agent type %s should succeed", agentType)

				var result ConfigurationResponseWrapper
				json.Unmarshal([]byte(output), &result)

				assert.Equal(t, "success", result.Status)

				var agentConfigID string
				if resultMap, ok := result.Result.(map[string]interface{}); ok {
					if entityGuid, ok := resultMap["entityGuid"].(string); ok {
						agentConfigID = entityGuid
					} else if configEntityGUID, ok := resultMap["configurationEntityGUID"].(string); ok {
						agentConfigID = configEntityGUID
					}
				}

				// Clean up
				if agentConfigID != "" {
					executeCommand([]string{
						"configuration", "delete",
						"--configuration-id", agentConfigID,
					})
				}

				fmt.Printf("✅ Successfully tested agent type: %s\n", agentType)
			})
		}
	})

	fmt.Println("✅ Successfully completed configuration CLI command lifecycle tests")
}
