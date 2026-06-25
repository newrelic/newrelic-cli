package fleetcontrol

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetSearch implements the 'search' command to search for fleet entities by name.
//
// This command searches for fleet entities and optionally filters them based on name criteria.
// Users can search using exact match (--name-equals), substring match (--name-contains),
// or retrieve all fleets (no flags). The flags are mutually exclusive.
//
// The command:
// 1. Validates that at most one search flag is provided (mutually exclusive)
// 2. Calls the New Relic API to search for fleet entities (type='FLEET')
// 3. Filters results based on the name criteria (if provided)
// 4. Returns filtered results in JSON or table format
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if search fails, nil on success
func handleFleetSearch(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values
	f := flags.Search()

	// Validate that both flags are not provided simultaneously (mutually exclusive)
	if f.NameEquals != "" && f.NameContains != "" {
		return PrintError(fmt.Errorf("--name-equals and --name-contains are mutually exclusive"))
	}

	// Call New Relic API to search for fleet entities
	// The backend only supports filtering by type, so we filter by name on the client side
	results, err := client.NRClient.FleetControl.GetEntitySearch("", "type='FLEET'")
	if err != nil {
		return PrintError(fmt.Errorf("failed to search fleets: %w", err))
	}

	if results == nil || len(results.Entities) == 0 {
		// Return empty array for JSON without wrapper (supports table format)
		return output.Print([]FleetEntityOutput{})
	}

	// Filter results based on name criteria
	var filteredFleets []FleetEntityOutput

	for _, entity := range results.Entities {
		// Type assert to fleet entity (pointer type)
		fleetEntity, ok := entity.(*fleetcontrol.EntityManagementFleetEntity)
		if !ok {
			// Skip non-fleet entities (shouldn't happen with type='FLEET' query)
			continue
		}

		// Apply name filter (if provided), otherwise include all entities
		matches := true // Default to true when no filters are provided
		if f.NameEquals != "" {
			matches = fleetEntity.Name == f.NameEquals
		} else if f.NameContains != "" {
			matches = strings.Contains(fleetEntity.Name, f.NameContains)
		}

		if matches {
			filteredFleet := FilterFleetEntityFromEntityManagement(*fleetEntity, f.ShowTags)
			filteredFleets = append(filteredFleets, *filteredFleet)
		}
	}

	// Output results directly without status wrapper (supports table format with --format=text)
	return output.Print(filteredFleets)
}
