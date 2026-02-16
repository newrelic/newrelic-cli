package fleetcontrol

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

// handleFleetGetUnassignedEntities implements the 'get-unassigned' command to retrieve unassigned entities.
//
// This command searches for all entities that are available for fleet management but not yet assigned to any fleet.
// Unassigned entities are identified by having tags.nr.supervisor set but tags.nr.fleet not set.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Builds a query to search for entities with the unassigned criteria
// 3. Calls the New Relic Entities API to retrieve matching entities
// 4. Optionally filters by entity type if specified
// 5. Returns the list of unassigned entities
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if search fails, nil on success
func handleFleetGetUnassignedEntities(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values
	f := flags.GetUnassigned()

	// Build query for unassigned entities:
	// Entities with tags.nr.supervisor set but tags.nr.fleet not set (IS NULL)
	query := "tags.nr.fleet IS NULL AND tags.nr.supervisor IS NOT NULL"

	// Add entity type filter if provided
	if f.EntityType != "" {
		query = fmt.Sprintf("(%s) AND type = '%s'", query, f.EntityType)
	}

	// Call New Relic Entities API to search for unassigned entities
	results, err := client.NRClient.Entities.GetEntitySearchByQueryWithContext(
		context.Background(),
		entities.EntitySearchOptions{
			Limit: f.Limit,
		},
		query,
		[]entities.EntitySearchSortCriteria{},
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to search for unassigned entities: %w", err))
	}

	if results == nil || results.Results.Entities == nil || len(results.Results.Entities) == 0 {
		// Return error when no unassigned entities are found
		return PrintError(fmt.Errorf("no unassigned entities found"))
	}

	// Format results into simplified entity output
	var unassignedEntities []EntityOutput

	for _, entity := range results.Results.Entities {
		// Extract common entity fields
		entityOutput := EntityOutput{
			ID:     string(entity.GetGUID()),
			Name:   entity.GetName(),
			Type:   entity.GetType(),
			Domain: string(entity.GetDomain()),
		}

		// Add tags only if --include-tags flag is true
		if f.IncludeTags && entity.GetTags() != nil {
			entityOutput.Tags = entity.GetTags()
		}

		unassignedEntities = append(unassignedEntities, entityOutput)
	}

	// Output results directly without status wrapper (supports table format with --format=text)
	return output.Print(unassignedEntities)
}
