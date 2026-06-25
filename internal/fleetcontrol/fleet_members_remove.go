package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetRemoveMembers implements the 'remove-members' command to remove entities from a fleet ring.
//
// This command removes one or more managed entities from a specific ring within a fleet.
// This operation removes entities from fleet management but does not delete the entities themselves.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Builds the member input with the ring and entity IDs
// 3. Calls the New Relic API to remove the members
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if removing members fails, nil on success
func handleFleetRemoveMembers(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.RemoveMembers()

	// Build member input with ring name and entity IDs
	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{
			Ring:      f.Ring,
			EntityIds: f.EntityIDs,
		},
	}

	// Call New Relic API to remove members from the fleet
	result, err := client.NRClient.FleetControl.FleetControlRemoveFleetMembers(f.FleetID, members)
	if err != nil {
		return PrintError(fmt.Errorf("failed to remove fleet members: %w", err))
	}

	// Print the result to stdout with status wrapper
	return PrintMemberSuccess(result)
}
