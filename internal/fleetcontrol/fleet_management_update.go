package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetUpdate implements the 'update' command to modify an existing fleet.
//
// This command allows updating the name, description, and tags of an existing fleet.
// Only the fields that are provided will be updated; others remain unchanged.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Validates that at least one field is provided for update
// 3. Builds an update input with only the provided fields
// 4. Calls the New Relic API to update the fleet
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if fleet update fails, nil on success
func handleFleetUpdate(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.Update()

	// Validate that at least one field besides ID is provided
	if f.Name == "" && f.Description == "" && len(f.Tags) == 0 {
		return PrintError(fmt.Errorf("at least one field must be provided for update: name, description, or tags"))
	}

	// Build update input - only include fields that were provided
	updateInput := fleetcontrol.FleetControlUpdateFleetEntityInput{}

	if f.Name != "" {
		updateInput.Name = f.Name
	}

	if f.Description != "" {
		updateInput.Description = f.Description
	}

	if len(f.Tags) > 0 {
		// Parse tags from "key:value1,value2" format
		tags, err := ParseTags(f.Tags)
		if err != nil {
			return PrintError(fmt.Errorf("invalid tags format: %w", err))
		}
		updateInput.Tags = tags
	}

	// Call New Relic API to update the fleet
	result, err := client.NRClient.FleetControl.FleetControlUpdateFleet(updateInput, f.ID)
	if err != nil {
		return PrintError(fmt.Errorf("failed to update fleet: %w", err))
	}

	// Filter and print only relevant fields from the fleet entity with status wrapper
	filteredOutput := FilterFleetEntity(result.Entity)
	return PrintFleetSuccess(filteredOutput)
}
