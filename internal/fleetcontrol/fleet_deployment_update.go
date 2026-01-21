package fleetcontrol

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// handleFleetUpdateDeployment implements the 'update-deployment' command to update an existing fleet deployment.
//
// This command updates a deployment's metadata and configuration, including:
//   - Name and description
//   - Configuration version list
//   - Tags
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Converts configuration version IDs if provided
// 3. Parses tags if provided
// 4. Builds update input with only provided fields
// 5. Calls the New Relic API to update the deployment
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if deployment update fails, nil on success
func handleFleetUpdateDeployment(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.UpdateDeployment()

	// Build the update input for the API
	updateInput := fleetcontrol.FleetControlFleetDeploymentUpdateInput{}

	// Add name if provided
	if f.Name != "" {
		updateInput.Name = f.Name
	}

	// Add description if provided
	if f.Description != "" {
		updateInput.Description = f.Description
	}

	// Convert configuration version IDs if provided
	if len(f.ConfigurationVersionIDs) > 0 {
		var configVersionList []fleetcontrol.FleetControlConfigurationVersionListInput
		for _, versionID := range f.ConfigurationVersionIDs {
			configVersionList = append(configVersionList, fleetcontrol.FleetControlConfigurationVersionListInput{
				ID: versionID,
			})
		}
		updateInput.ConfigurationVersionList = configVersionList
	}

	// Parse tags if provided
	if len(f.Tags) > 0 {
		tags, err := ParseTags(f.Tags)
		if err != nil {
			return PrintError(fmt.Errorf("invalid tags format: %w", err))
		}
		updateInput.Tags = tags
	}

	// Call New Relic API to update the deployment
	result, err := client.NRClient.FleetControl.FleetControlUpdateFleetDeployment(
		updateInput,
		f.ID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to update fleet deployment: %w", err))
	}

	// Print the updated deployment entity to stdout with status wrapper
	return PrintDeploymentSuccess(result.Entity)
}
