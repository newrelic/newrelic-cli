package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetCreateDeployment implements the 'create-deployment' command to create a new fleet deployment.
//
// This command creates a deployment for managing configuration rollouts across fleet rings.
// A deployment tracks the rollout progress and status across different rings in a phased manner.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Converts configuration version IDs from string slice to proper input format
// 3. Derives organization scope if not explicitly provided
// 4. Parses tags into the required format
// 5. Calls the New Relic API to create the deployment
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if deployment creation fails, nil on success
func handleFleetCreateDeployment(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.CreateDeployment()

	// Convert configuration version IDs to the required format
	var configVersionList []fleetcontrol.FleetControlConfigurationVersionListInput
	for _, versionID := range f.ConfigurationVersionIDs {
		configVersionList = append(configVersionList, fleetcontrol.FleetControlConfigurationVersionListInput{
			ID: versionID,
		})
	}

	// Determine scope for the deployment (defaults to organization scope)
	var scopeID string
	var scopeType fleetcontrol.FleetControlEntityScope

	// Fetch organization details from API
	org, err := client.NRClient.Organization.GetOrganization()
	if err != nil {
		return fmt.Errorf("failed to get organization: %w", err)
	}
	scopeID = org.ID
	scopeType = fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION

	// Parse tags from "key:value1,value2" format to API format
	tags, err := ParseTags(f.Tags)
	if err != nil {
		return PrintError(fmt.Errorf("invalid tags format: %w", err))
	}

	// Build the create input for the API
	createInput := fleetcontrol.FleetControlFleetDeploymentCreateInput{
		FleetId:                  f.FleetID,
		Name:                     f.Name,
		ConfigurationVersionList: configVersionList,
		Scope: fleetcontrol.FleetControlScopedReferenceInput{
			ID:   scopeID,
			Type: scopeType,
		},
	}

	// Add optional fields if provided
	if f.Description != "" {
		createInput.Description = f.Description
	}

	if len(tags) > 0 {
		createInput.Tags = tags
	}

	// Call New Relic API to create the deployment
	result, err := client.NRClient.FleetControl.FleetControlCreateFleetDeployment(createInput)
	if err != nil {
		return PrintError(fmt.Errorf("failed to create fleet deployment: %w", err))
	}

	// Print the created deployment entity to stdout with status wrapper
	return PrintDeploymentSuccess(result.Entity)
}
