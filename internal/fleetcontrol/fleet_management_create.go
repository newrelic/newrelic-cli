package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetCreate implements the 'create' command to create a new fleet entity.
//
// This command creates a fleet for managing collections of hosts or Kubernetes clusters.
// Fleets allow you to organize entities into rings for controlled deployment rollouts.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Maps validated strings to client library types using helper functions
// 3. Fetches organization ID (uses provided value or fetches from API)
// 4. Parses tags into the required format
// 5. Calls the New Relic API to create the fleet
//
// Note: Fleets are always scoped to the organization level.
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if fleet creation fails, nil on success
func handleFleetCreate(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.Create()

	// Map validated managed entity type to client library type
	// YAML validation has already confirmed this value is in allowed_values
	// This function will return an error if the mapping fails (should never happen)
	entityType, err := MapManagedEntityType(f.ManagedEntityType)
	if err != nil {
		return PrintError(err)
	}

	// Validate operating system requirements based on entity type
	// For HOST fleets, operating system must be specified
	// For KUBERNETESCLUSTER fleets, operating system should not be specified
	if entityType == fleetcontrol.FleetControlManagedEntityTypeTypes.HOST {
		if f.OperatingSystem == "" {
			return PrintError(fmt.Errorf("--operating-system is required when --managed-entity-type is HOST (allowed values: LINUX, WINDOWS)"))
		}
	} else if entityType == fleetcontrol.FleetControlManagedEntityTypeTypes.KUBERNETESCLUSTER {
		if f.OperatingSystem != "" {
			return PrintError(fmt.Errorf("--operating-system should not be specified for KUBERNETESCLUSTER fleets"))
		}
	}

	// Get organization ID (either from flag or fetch from API)
	// This avoids an unnecessary API call if the user already knows their org ID
	orgID := GetOrganizationID(f.OrganizationID)
	if orgID == "" {
		return PrintError(fmt.Errorf("failed to determine organization ID"))
	}

	// Parse tags from "key:value1,value2" format to API format
	tags, err := ParseTags(f.Tags)
	if err != nil {
		return PrintError(fmt.Errorf("invalid tags format: %w", err))
	}

	// Build the create input for the API
	// Fleets are always scoped to the organization level
	createInput := fleetcontrol.FleetControlFleetEntityCreateInput{
		Name:              f.Name,
		ManagedEntityType: entityType,
		Scope: fleetcontrol.FleetControlScopedReferenceInput{
			ID:   orgID,
			Type: fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION,
		},
	}

	// Add optional fields if provided
	if f.Description != "" {
		createInput.Description = f.Description
	}

	if f.Product != "" {
		createInput.Product = f.Product
	}

	// Add operating system if provided (only applicable for HOST fleets)
	// For KUBERNETESCLUSTER fleets, this should not be specified
	// Using a pointer ensures nil is sent instead of an empty object when not set
	if f.OperatingSystem != "" {
		osType, err := MapOperatingSystemType(f.OperatingSystem)
		if err != nil {
			return PrintError(err)
		}
		createInput.OperatingSystem = &fleetcontrol.FleetControlOperatingSystemCreateInput{
			Type: osType,
		}
	}

	if len(tags) > 0 {
		createInput.Tags = tags
	}

	// Call New Relic API to create the fleet
	result, err := client.NRClient.FleetControl.FleetControlCreateFleet(createInput)
	if err != nil {
		return PrintError(fmt.Errorf("failed to create fleet: %w", err))
	}

	// Filter and print only relevant fields from the fleet entity with status wrapper
	filteredOutput := FilterFleetEntity(result.Entity)
	return PrintFleetSuccess(filteredOutput)
}
