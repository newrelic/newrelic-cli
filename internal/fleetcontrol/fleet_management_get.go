package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetGet implements the 'get' command to retrieve a fleet entity by ID.
//
// This command fetches detailed information about a specific fleet including its
// configuration, managed entity type, scope, tags, and metadata.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Calls the New Relic API to retrieve the entity by ID using GetEntity
// 3. Validates that the entity exists and is a fleet
// 4. Filters and returns the fleet details
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if fleet retrieval fails, nil on success
func handleFleetGet(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values
	f := flags.Get()

	// Use GetEntity to fetch the entity directly by ID
	entityInterface, err := client.NRClient.FleetControl.GetEntity(f.FleetID)
	if err != nil {
		return PrintError(fmt.Errorf("failed to get fleet: %w", err))
	}

	if entityInterface == nil {
		return PrintError(fmt.Errorf("fleet with ID '%s' not found", f.FleetID))
	}

	// Dereference the interface pointer and type assert to fleet entity
	fleetEntity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetEntity)
	if !ok {
		return PrintError(fmt.Errorf("entity '%s' is not a fleet (type: %T)", f.FleetID, *entityInterface))
	}

	// Convert to filtered output format with status wrapper
	// Always show tags for get command since it's fetching a specific entity
	filteredOutput := FilterFleetEntityFromEntityManagement(*fleetEntity, true)
	return PrintFleetSuccess(filteredOutput)
}
