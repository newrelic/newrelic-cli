package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetDeleteDeployment implements the 'delete-deployment' command to delete a fleet deployment.
//
// This command deletes a fleet deployment.
// This operation cannot be undone and will remove the deployment.
// The deployment must not be actively in progress.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Calls the New Relic API to delete the deployment
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if deployment deletion fails, nil on success
func handleFleetDeleteDeployment(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.DeleteDeployment()

	// Call New Relic API to delete the deployment
	_, err := client.NRClient.FleetControl.FleetControlDeleteFleetDeployment(
		f.DeploymentID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to delete fleet deployment: %w", err))
	}

	// Print the deletion result to stdout with status wrapper
	return PrintDeploymentDeleteSuccess(f.DeploymentID)
}
