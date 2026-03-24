package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetDeleteVersion implements the 'delete-version' command to delete a configuration version.
//
// This command deletes a specific version of a fleet configuration.
// This operation cannot be undone. The configuration entity itself will remain.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Derives organization ID if not provided
// 3. Calls the New Relic API to delete the version
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if version deletion fails, nil on success
func handleFleetDeleteVersion(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.DeleteVersion()

	// Get organization ID (provided or fetched from API)
	orgID := GetOrganizationID(f.OrganizationID)

	// Call New Relic API to delete the version
	err := client.NRClient.FleetControl.FleetControlDeleteConfigurationVersion(
		f.VersionID,
		orgID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to delete configuration version: %w", err))
	}

	// Print success message to stdout with status wrapper
	return PrintConfigurationDeleteSuccess(f.VersionID)
}
