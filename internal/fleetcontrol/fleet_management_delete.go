package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
)

// handleFleetDelete implements the 'delete' command to remove one or more fleets.
//
// This command permanently deletes fleets. This operation cannot be undone.
// The entities managed by the fleet are not deleted, only the fleet itself.
//
// The command supports two modes:
// 1. Single deletion: Use --fleet-id to delete one fleet
// 2. Bulk deletion: Use --fleet-ids to delete multiple fleets
//
// The flags --fleet-id and --fleet-ids are mutually exclusive.
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if validation fails or deletion fails, nil on success
func handleFleetDelete(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values
	f := flags.Delete()

	// Validate that exactly one of --fleet-id or --fleet-ids is provided (mutually exclusive)
	if f.FleetID == "" && len(f.FleetIDs) == 0 {
		return PrintError(fmt.Errorf("one of --fleet-id or --fleet-ids must be provided"))
	}

	if f.FleetID != "" && len(f.FleetIDs) > 0 {
		return PrintError(fmt.Errorf("--fleet-id and --fleet-ids are mutually exclusive, use only one"))
	}

	// Handle single deletion with --fleet-id
	if f.FleetID != "" {
		return deleteSingleFleet(f.FleetID)
	}

	// Handle bulk deletion with --fleet-ids
	return deleteBulkFleets(f.FleetIDs)
}

// deleteSingleFleet deletes a single fleet by ID and returns the result.
func deleteSingleFleet(id string) error {
	result, err := client.NRClient.FleetControl.FleetControlDeleteFleet(id)
	if err != nil {
		return PrintError(fmt.Errorf("failed to delete fleet: %w", err))
	}

	return PrintDeleteSuccess(result.ID)
}

// deleteBulkFleets deletes multiple fleets and returns the results as a list.
// If --fleet-ids contains only one ID, suggests using --fleet-id instead.
func deleteBulkFleets(ids []string) error {
	// Validate that --fleet-ids has more than one ID
	if len(ids) == 1 {
		return PrintError(fmt.Errorf("--fleet-ids contains only one ID, use --fleet-id instead for single fleet deletion"))
	}

	// Collect results for each deletion
	var results []FleetDeleteResponseWrapper

	for _, id := range ids {
		result, err := client.NRClient.FleetControl.FleetControlDeleteFleet(id)
		if err != nil {
			// For bulk operations, we collect errors as failed responses rather than stopping
			results = append(results, FleetDeleteResponseWrapper{
				Status: "failed",
				Error:  fmt.Sprintf("failed to delete fleet: %v", err),
				ID:     id,
			})
		} else {
			results = append(results, FleetDeleteResponseWrapper{
				Status: "success",
				Error:  "",
				ID:     result.ID,
			})
		}
	}

	// Print all results as a list
	return PrintDeleteBulkSuccess(results)
}
