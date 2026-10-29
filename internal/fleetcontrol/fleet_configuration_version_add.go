package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetAddVersion implements the 'add-version' command to add a version to an existing configuration.
//
// This command adds a new version to an existing fleet configuration.
// It creates a new version by using the configuration GUID in the custom headers.
//
// The configuration must be provided via one of two mutually exclusive flags:
//   - --configuration-file-path: Path to a file (recommended for production)
//   - --configuration-content: Inline content (for testing/development/emergency only)
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Validates mutual exclusivity of configuration flags
// 3. Derives organization ID if not provided
// 4. Builds custom headers with the configuration GUID
// 5. Calls the New Relic API to add the version
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if version addition fails, nil on success
func handleFleetAddVersion(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	// Note: AddVersion() returns error because it handles file reading
	f, err := flags.AddVersion()
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

	// Build custom headers with the configuration GUID
	// This tells the API to add a version to the existing configuration
	customHeaders := map[string]interface{}{
		"x-newrelic-client-go-custom-headers": map[string]string{
			"Newrelic-Entity": fmt.Sprintf(
				`{"agentConfiguration": "%s"}`,
				f.ConfigurationID,
			),
		},
	}

	// Call New Relic API to add the version
	// Note: This uses the same endpoint as CreateConfiguration but with different headers
	result, err := client.NRClient.FleetControl.FleetControlCreateConfiguration(
		configBodyBytes,
		customHeaders,
		orgID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to add configuration version: %w", err))
	}

	// Print the created version to stdout with status wrapper
	return PrintConfigurationSuccess(result)
}
