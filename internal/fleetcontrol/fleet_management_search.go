package fleetcontrol

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// maxSearchPages bounds how many pages searchFleetEntities will fetch, as a safety net
// against a misbehaving API returning a cursor that never terminates.
const maxSearchPages = 1000

// searchFleetEntities fetches every page of type='FLEET' entities via GetEntitySearch,
// following NextCursor until the API reports no more pages, then applies the name filter
// against the full accumulated set. The API only supports filtering by type, so name
// filtering has to happen client-side - which only works correctly if every page has
// been fetched first.
func searchFleetEntities(f SearchFlags) ([]FleetEntityOutput, error) {
	var allEntities []fleetcontrol.EntityManagementEntityInterface

	var cursor *string
	for range maxSearchPages {
		results, err := client.NRClient.FleetControl.GetEntitySearch(cursor, "type='FLEET'")
		if err != nil {
			return nil, fmt.Errorf("failed to search fleets: %w", err)
		}
		if results == nil {
			break
		}

		allEntities = append(allEntities, results.Entities...)

		if results.NextCursor == "" {
			break
		}
		nextCursor := results.NextCursor
		cursor = &nextCursor
	}

	var filteredFleets []FleetEntityOutput
	for _, entity := range allEntities {
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

	return filteredFleets, nil
}

// handleFleetSearch implements the 'search' command to search for fleet entities by name.
//
// This command searches for fleet entities and optionally filters them based on name criteria.
// Users can search using exact match (--name-equals), substring match (--name-contains),
// or retrieve all fleets (no flags). The flags are mutually exclusive.
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if search fails, nil on success
func handleFleetSearch(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values
	f := flags.Search()

	// Validate that both flags are not provided simultaneously (mutually exclusive)
	if f.NameEquals != "" && f.NameContains != "" {
		return PrintError(fmt.Errorf("--name-equals and --name-contains are mutually exclusive"))
	}

	filteredFleets, err := searchFleetEntities(f)
	if err != nil {
		return PrintError(err)
	}

	if filteredFleets == nil {
		// Return empty array for JSON without wrapper (supports table format)
		filteredFleets = []FleetEntityOutput{}
	}

	// Output results directly without status wrapper (supports table format with --format=text)
	return output.Print(filteredFleets)
}
