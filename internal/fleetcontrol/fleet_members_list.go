package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// FleetMemberEntityWithoutTags is a filtered version of FleetControlFleetMemberEntityResult
// that excludes the Tags field, used when --include-tags is false
type FleetMemberEntityWithoutTags struct {
	ID       string                                   `json:"id"`
	Metadata fleetcontrol.FleetControlMetadata        `json:"metadata"`
	Name     string                                   `json:"name"`
	Scope    fleetcontrol.FleetControlScopedReference `json:"scope"`
	Type     string                                   `json:"type"`
}

// FleetMembersResultWithoutTags is the result structure without tags
type FleetMembersResultWithoutTags struct {
	Items      []FleetMemberEntityWithoutTags `json:"items"`
	NextCursor string                         `json:"nextCursor,omitempty"`
}

// handleFleetListMembers implements the 'list-members' command to retrieve entities in a fleet.
//
// This command fetches all managed entities that are members of a specific fleet.
// Results can optionally be filtered to show only entities in a specific ring.
// Tags can be excluded from the output (default) or included via the --show-tags flag.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Builds the filter input with fleet ID and optional ring name
// 3. Calls the New Relic API to retrieve the fleet members
// 4. Filters out tags if --show-tags is false (default behavior)
// 5. Returns the list of member entities with their details
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if listing members fails, nil on success
func handleFleetListMembers(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.ListMembers()

	// Build filter input with fleet ID and optional ring
	filter := &fleetcontrol.FleetControlFleetMembersFilterInput{
		FleetId: f.FleetID,
	}

	// Add ring filter if provided
	if f.Ring != "" {
		filter.Ring = f.Ring
	}

	// Call New Relic API to get fleet members
	// Empty cursor for first page - pagination can be added in the future if needed
	result, err := client.NRClient.FleetControl.GetFleetMembers("", filter)
	if err != nil {
		return PrintError(fmt.Errorf("failed to list fleet members: %w", err))
	}

	// Filter out tags if --show-tags is false (default)
	// Use a custom struct that completely omits the Tags field
	if !f.ShowTags {
		filteredResult := &FleetMembersResultWithoutTags{
			NextCursor: result.NextCursor,
			Items:      make([]FleetMemberEntityWithoutTags, len(result.Items)),
		}

		// Copy each item without the tags field
		for i, item := range result.Items {
			filteredResult.Items[i] = FleetMemberEntityWithoutTags{
				ID:       item.ID,
				Name:     item.Name,
				Type:     item.Type,
				Scope:    item.Scope,
				Metadata: item.Metadata,
			}
		}

		// Print the result without tags
		return PrintMemberSuccess(filteredResult)
	}

	// Print the result with tags included
	return PrintMemberSuccess(result)
}
