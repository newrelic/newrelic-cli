package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetAddMembers implements the 'add-members' command to add entities to a fleet ring.
//
// This command adds one or more managed entities to a specific ring within a fleet.
// Rings are used to organize entities for controlled deployment rollouts (e.g., canary, production).
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Builds the member input with the ring and entity IDs
// 3. Calls the New Relic API to add the members
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if adding members fails, nil on success
func handleFleetAddMembers(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.AddMembers()

	// Build member input with ring name and entity IDs
	members := []fleetcontrol.FleetControlFleetMemberRingInput{
		{
			Ring:      f.Ring,
			EntityIds: f.EntityIDs,
		},
	}

	// Call New Relic API to add members to the fleet
	result, err := client.NRClient.FleetControl.FleetControlAddFleetMembers(f.FleetID, members)
	if err != nil {
		return PrintError(fmt.Errorf("failed to add fleet members: %w", err))
	}

	// Print the result to stdout with status wrapper
	return PrintMemberSuccess(result)
}
