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
// The command supports two syntaxes (mutually exclusive):
// 1. New syntax: --agent "AgentType:Version:ConfigVersionID1,ConfigVersionID2" (supports multiple agents)
// 2. Legacy syntax: --agent-type, --agent-version, --configuration-version-ids (single agent only)
//
// The command:
// 1. Validates flag values and mutual exclusivity
// 2. Parses agent specifications from --agent flag or legacy flags
// 3. Builds agent configuration structure with configuration version IDs
// 4. Derives organization scope if not explicitly provided
// 5. Parses tags into the required format
// 6. Calls the New Relic API to create the deployment
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// Returns:
//   - Error if deployment creation fails, nil on success
func handleFleetCreateDeployment(_ *cobra.Command, _ []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.CreateDeployment()

	// Validate mutual exclusivity between --agent and legacy flags
	hasNewSyntax := len(f.Agent) > 0
	hasLegacySyntax := f.AgentType != "" || f.AgentVersion != "" || len(f.ConfigurationVersionIDs) > 0

	if hasNewSyntax && hasLegacySyntax {
		return PrintError(fmt.Errorf("cannot use --agent with --agent-type, --agent-version, or --configuration-version-ids. Use either the new --agent syntax or the legacy flags, not both"))
	}

	if !hasNewSyntax && !hasLegacySyntax {
		return PrintError(fmt.Errorf("must specify agent configuration using either --agent flag or the legacy --agent-type, --agent-version, and --configuration-version-ids flags"))
	}

	// Build agent input based on which syntax is used
	var agents []fleetcontrol.FleetControlAgentInput

	if hasNewSyntax {
		// Parse new --agent syntax: "AgentType:Version:ConfigVersionID1,ConfigVersionID2"
		for _, agentSpec := range f.Agent {
			agent, err := ParseAgentSpec(agentSpec)
			if err != nil {
				return PrintError(fmt.Errorf("invalid --agent format '%s': %w", agentSpec, err))
			}
			agents = append(agents, agent)
		}
	} else {
		// Legacy syntax: validate all required fields are present
		if f.AgentType == "" {
			return PrintError(fmt.Errorf("--agent-type is required when not using --agent flag"))
		}
		if f.AgentVersion == "" {
			return PrintError(fmt.Errorf("--agent-version is required when not using --agent flag"))
		}
		if len(f.ConfigurationVersionIDs) == 0 {
			return PrintError(fmt.Errorf("--configuration-version-ids is required when not using --agent flag"))
		}

		// Convert configuration version IDs to the required format for agent input
		var configVersionList []fleetcontrol.FleetControlConfigurationVersionListInput
		for _, versionID := range f.ConfigurationVersionIDs {
			configVersionList = append(configVersionList, fleetcontrol.FleetControlConfigurationVersionListInput{
				ID: versionID,
			})
		}

		// Build single agent input from legacy flags
		agents = []fleetcontrol.FleetControlAgentInput{
			{
				AgentType:                f.AgentType,
				ConfigurationVersionList: configVersionList,
				Version:                  f.AgentVersion,
			},
		}
	}

	// Validate agent versions are compatible with the fleet type
	// This checks that "*" is only used with KUBERNETESCLUSTER fleets, not HOST fleets
	if err := ValidateAgentVersionsForFleet(f.FleetID, agents); err != nil {
		return PrintError(err)
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
		FleetId: f.FleetID,
		Name:    f.Name,
		Agents:  agents,
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
