package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetGetConfiguration implements the 'get-configuration' command to retrieve a fleet configuration.
//
// This command retrieves a fleet configuration or a specific version of a configuration.
// You can:
//   - Get the latest version of a configuration
//   - Get a specific version by number
//   - Get a configuration version by its entity GUID
//   - Get entity metadata (e.g. configurationType) instead of the raw content, via --show-metadata
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Maps the validated mode string to the client library type
// 3. Derives organization ID if not provided
// 4. Calls the New Relic API to retrieve the configuration
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if configuration retrieval fails, nil on success
func handleFleetGetConfiguration(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.GetConfiguration()

	// --show-metadata reads entity attributes via the EntityManagement API instead of the raw
	// blob content, so --mode and --version (which only affect the Blob Service request) don't apply.
	if f.ShowMetadata {
		if flags.Has("mode") || flags.Has("version") {
			return PrintError(fmt.Errorf("--show-metadata is mutually exclusive with --mode and --version"))
		}
		return handleFleetGetConfigurationMetadata(f.ConfigurationID)
	}

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Map validated mode to client library type
	// YAML validation has already confirmed this value is in allowed_values
	mode, err := MapConfigurationMode(f.Mode)
	if err != nil {
		return PrintError(err)
	}

	// Call New Relic API to get the configuration
	result, err := client.NRClient.FleetControl.FleetControlGetConfiguration(
		f.ConfigurationID,
		orgID,
		mode,
		f.Version,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to get configuration: %w", err))
	}

	// Print the raw configuration directly to stdout
	// Bypass the output formatter entirely since this is raw YAML content
	// that should be displayed exactly as returned by the API
	fmt.Print(string(*result))
	return nil
}

// handleFleetGetConfigurationMetadata retrieves and prints entity metadata for a configuration
// or configuration version via the EntityManagement API. The Blob Service (used by the default
// mode of this command) only returns raw configuration content, never entity attributes like
// configurationType.
func handleFleetGetConfigurationMetadata(configurationID string) error {
	entityInterface, err := client.NRClient.FleetControl.GetEntity(configurationID)
	if err != nil {
		return PrintError(fmt.Errorf("failed to get configuration metadata: %w", err))
	}

	if entityInterface == nil {
		return PrintError(fmt.Errorf("configuration with ID '%s' not found", configurationID))
	}

	switch entity := (*entityInterface).(type) {
	case *fleetcontrol.EntityManagementAgentConfigurationEntity:
		return PrintConfigurationSuccess(FilterAgentConfigurationEntityFromEntityManagement(*entity))
	case *fleetcontrol.EntityManagementAgentConfigurationVersionEntity:
		return PrintConfigurationSuccess(FilterAgentConfigurationVersionEntityFromEntityManagement(*entity))
	default:
		return PrintError(fmt.Errorf("entity '%s' is not a fleet configuration or configuration version (type: %T)", configurationID, *entityInterface))
	}
}
