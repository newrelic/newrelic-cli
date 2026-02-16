package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetCreateConfiguration implements the 'create-configuration' command to create a fleet configuration.
//
// This command creates a new fleet configuration with custom attributes.
// Configurations allow you to define settings and attributes for fleet management.
//
// The configuration body must be provided via one of two mutually exclusive flags:
//   - --configuration-file-path: Path to a file (recommended for production)
//   - --configuration-content: Inline content (for testing/development/emergency only)
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Validates mutual exclusivity of configuration flags
// 3. Derives organization ID if not provided
// 4. Builds custom headers with entity metadata
// 5. Calls the New Relic API to create the configuration
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if configuration creation fails, nil on success
func handleFleetCreateConfiguration(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	// Note: CreateConfiguration() returns error because it handles file reading
	f, err := flags.CreateConfiguration()
	if err != nil {
		return fmt.Errorf("failed to read flags: %w", err)
	}

	// Validate that exactly one of --configuration-file-path or --configuration-content is provided
	hasFilePath := f.ConfigurationFilePath != ""
	hasContent := f.ConfigurationContent != ""

	if !hasFilePath && !hasContent {
		return fmt.Errorf("one of --configuration-file-path or --configuration-content must be provided")
	}

	if hasFilePath && hasContent {
		return fmt.Errorf("--configuration-file-path and --configuration-content are mutually exclusive, use only one")
	}

	// Determine which configuration content to use
	configBody := f.ConfigurationFilePath
	if hasContent {
		configBody = f.ConfigurationContent
	}

	// Convert the configuration body to []byte to prevent JSON marshaling
	// This preserves newlines and formatting in the configuration file
	configBodyBytes := []byte(configBody)

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Build custom headers required by the API
	// These headers specify the entity name, agent type, and managed entity type
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": fmt.Sprintf(
				`{"name": "%s", "agentType": "%s", "managedEntityType": "%s"}`,
				f.Name,
				f.AgentType,
				f.ManagedEntityType,
			),
		},
	}

	// Call New Relic API to create the configuration
	result, err := client.NRClient.FleetControl.FleetControlCreateConfiguration(
		configBodyBytes,
		customHeaders,
		orgID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to create fleet configuration: %w", err))
	}

	// Print the created configuration to stdout with status wrapper
	return PrintConfigurationSuccess(result)
}
