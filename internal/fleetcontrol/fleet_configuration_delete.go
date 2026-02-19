package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetDeleteConfiguration implements the 'delete-configuration' command to delete a fleet configuration.
//
// This command deletes a fleet configuration and all its versions.
// This operation cannot be undone and will remove all versions of the configuration.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Derives organization ID if not provided
// 3. Calls the New Relic API to delete the configuration
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if configuration deletion fails, nil on success
func handleFleetDeleteConfiguration(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.DeleteConfiguration()

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Call New Relic API to delete the configuration
	_, err := client.NRClient.FleetControl.FleetControlDeleteConfiguration(
		f.ConfigurationID,
		orgID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to delete configuration: %w", err))
	}

	// Print the deletion result to stdout with status wrapper
	return PrintConfigurationDeleteSuccess(f.ConfigurationID)
}
