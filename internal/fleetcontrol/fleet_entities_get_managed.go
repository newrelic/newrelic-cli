package fleetcontrol

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

// EntityOutput is a simplified representation of an entity
type EntityOutput struct {
	ID     string               `json:"id"`
	Name   string               `json:"name"`
	Type   string               `json:"type"`
	Domain string               `json:"domain"`
	Tags   []entities.EntityTag `json:"tags,omitempty"`
}

// handleFleetGetManagedEntities implements the 'get-managed' command to retrieve managed entities.
//
// This command searches for all entities that are currently managed by any fleet.
// Managed entities are identified by having both tags.nr.fleet and tags.nr.supervisor set.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Builds a query to search for entities with required tags
// 3. Calls the New Relic Entities API to retrieve matching entities
// 4. Optionally filters by entity type if specified
// 5. Returns the list of managed entities
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if search fails, nil on success
func handleFleetGetManagedEntities(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values
	f := flags.GetManaged()

	// Build query for managed entities:
	// Entities with both tags.nr.fleet and tags.nr.supervisor set (IS NOT NULL)
	query := "tags.nr.fleet IS NOT NULL AND tags.nr.supervisor IS NOT NULL"

	// Add entity type filter if provided
	if f.EntityType != "" {
		query = fmt.Sprintf("(%s) AND type = '%s'", query, f.EntityType)
	}

	// Call New Relic Entities API to search for managed entities
	results, err := client.NRClient.Entities.GetEntitySearchByQueryWithContext(
		context.Background(),
		entities.EntitySearchOptions{
			Limit: f.Limit,
		},
		query,
		[]entities.EntitySearchSortCriteria{},
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to search for managed entities: %w", err))
	}

	if results == nil || results.Results.Entities == nil || len(results.Results.Entities) == 0 {
		// Return error when no managed entities are found
		return PrintError(fmt.Errorf("no managed entities found"))
	}

	// Format results into simplified entity output
	var managedEntities []EntityOutput

	for _, entity := range results.Results.Entities {
		// Extract common entity fields
		entityOutput := EntityOutput{
			ID:     string(entity.GetGUID()),
			Name:   entity.GetName(),
			Type:   entity.GetType(),
			Domain: entity.GetDomain(),
		}

		// Add tags only if --include-tags flag is true
		if f.IncludeTags && entity.GetTags() != nil {
			entityOutput.Tags = entity.GetTags()
		}

		managedEntities = append(managedEntities, entityOutput)
	}

	// Output results directly without status wrapper (supports table format with --format=text)
	return output.Print(managedEntities)
}
