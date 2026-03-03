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
//   - Agent configurations (type, version, and configuration versions)
//   - Tags
//
// The command supports two syntaxes (mutually exclusive):
// 1. New syntax: --agent "AgentType:Version:ConfigVersionID1,ConfigVersionID2" (supports multiple agents)
// 2. Legacy syntax: --configuration-version-ids (updates only configuration versions, not agent type/version)
//
// The command:
// 1. Validates flag values and mutual exclusivity
// 2. Parses agent specifications from --agent flag or legacy flags
// 3. Builds update input with only provided fields
// 4. Parses tags if provided
// 5. Calls the New Relic API to update the deployment
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if deployment update fails, nil on success
func handleFleetUpdateDeployment(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.UpdateDeployment()

	// Validate mutual exclusivity between --agent and legacy flags
	hasNewSyntax := len(f.Agent) > 0
	hasLegacySyntax := len(f.ConfigurationVersionIDs) > 0

	if hasNewSyntax && hasLegacySyntax {
		return PrintError(fmt.Errorf("cannot use --agent with --configuration-version-ids. Use either the new --agent syntax or the legacy flag, not both"))
	}

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

	// Handle agent configuration based on which syntax is used
	if hasNewSyntax {
		// Parse new --agent syntax: "AgentType:Version:ConfigVersionID1,ConfigVersionID2"
		var agents []fleetcontrol.FleetControlAgentInput
		for _, agentSpec := range f.Agent {
			agent, err := ParseAgentSpec(agentSpec)
			if err != nil {
				return PrintError(fmt.Errorf("invalid --agent format '%s': %w", agentSpec, err))
			}
			agents = append(agents, agent)
		}
		updateInput.Agents = agents
	} else if hasLegacySyntax {
		// Legacy syntax: Convert configuration version IDs only
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
		f.DeploymentID,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to update fleet deployment: %w", err))
	}

	// Print the updated deployment entity to stdout with status wrapper
	return PrintDeploymentSuccess(result.Entity)
}
