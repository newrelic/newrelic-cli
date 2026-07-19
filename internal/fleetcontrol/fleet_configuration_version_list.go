package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetGetConfigurationVersions implements the 'get-versions' command to retrieve all versions of a fleet configuration.
//
// This command retrieves the complete version history for a fleet configuration,
// including version numbers, blob IDs, entity GUIDs, and timestamps.
//
// Use cases:
//   - View all available versions of a configuration
//   - Check version history and timestamps
//   - Find specific version details for rollback or comparison
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Derives organization ID if not provided
// 3. Calls the New Relic API to retrieve all configuration versions
// 4. Returns a list of versions with their metadata
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if version retrieval fails, nil on success
func handleFleetGetConfigurationVersions(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.GetVersions()

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Call New Relic API to get all configuration versions
	result, err := client.NRClient.FleetControl.FleetControlGetConfigurationVersions(
		f.ConfigurationID,
		orgID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to get configuration versions: %w", err))
	}

	// Validate that versions were returned
	if result == nil || len(result.Versions) == 0 {
		return PrintError(fmt.Errorf("no version details found, please check the GUID of the configuration entity provided"))
	}

	// Print the versions list to stdout with status wrapper
	return PrintConfigurationSuccess(result)
}
