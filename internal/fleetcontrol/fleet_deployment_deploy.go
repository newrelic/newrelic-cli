package fleetcontrol

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// spinner displays an animated spinner for the given duration
// to provide visual feedback during wait periods
func spinner(duration time.Duration, message string) {
	spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	done := time.After(duration)
	i := 0

	for {
		select {
		case <-done:
			// Clear the spinner line
			fmt.Print("\r\033[K")
			return
		case <-ticker.C:
			fmt.Printf("\r%s %s", spinChars[i%len(spinChars)], message)
			i++
		}
	}
}

// handleFleetDeploy implements the 'deploy' command to trigger a fleet deployment.
//
// This command triggers the actual deployment process for a fleet deployment,
// initiating the rollout of configurations across specified rings.
// After triggering, it polls the deployment status until completion.
//
// The command:
// 1. Validates flag values (done automatically by framework via YAML rules)
// 2. Constructs the deployment policy with the specified rings
// 3. Calls the New Relic API to trigger the deployment
// 4. Polls the deployment entity to track progress
// 5. Displays status updates while deployment is in progress
// 6. Shows final status when deployment completes (success or failure)
//
// Parameters:
//   - cmd: The cobra command being executed
//   - args: Command arguments (not used)
//   - flags: Validated flag values from YAML configuration
//
// # Returns error if deployment trigger fails, nil on success
//
// suppressing the linter here, since breaking this function down to address cyclomatic complexity concerns
// could affect the code-perception of the hierarchical structure of the fleet commands
//
//nolint:gocyclo
func handleFleetDeploy(cmd *cobra.Command, args []string, flags *FlagValues) error {
	// Get typed flag values - no hardcoded strings!
	f := flags.Deploy()

	// Build the deployment policy input
	policy := fleetcontrol.FleetControlFleetDeploymentPolicyInput{
		RingDeploymentPolicy: fleetcontrol.FleetControlRingDeploymentPolicyInput{
			RingsToDeploy: f.RingsToDeploy,
		},
	}

	// Call New Relic API to trigger the deployment
	result, err := client.NRClient.FleetControl.FleetControlDeploy(
		f.DeploymentID,
		policy,
	)
	if err != nil {
		return PrintError(fmt.Errorf("failed to trigger deployment: %w", err))
	}

	// Display initial message with rocket emoji
	fmt.Printf("🚀 Deployment triggered successfully (ID: %s)\n", result.ID)
	fmt.Printf("📊 Monitoring deployment progress...\n\n")

	// Poll the deployment entity to track progress
	deploymentID := result.ID
	pollInterval := 5 * time.Second
	maxPolls := 360 // 30 minutes at 5 second intervals

	for i := 0; i < maxPolls; i++ {
		// Fetch the deployment entity
		entityInterface, err := client.NRClient.FleetControl.GetEntity(deploymentID)
		if err != nil {
			return PrintError(fmt.Errorf("failed to fetch deployment status: %w", err))
		}

		if entityInterface == nil {
			return PrintError(fmt.Errorf("deployment entity '%s' not found", deploymentID))
		}

		// Type assert to deployment entity
		deploymentEntity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetDeploymentEntity)
		if !ok {
			return PrintError(fmt.Errorf("entity '%s' is not a deployment (type: %T)", deploymentID, *entityInterface))
		}

		// Check the deployment phase
		phase := deploymentEntity.Phase

		// Display status
		switch phase {
		case fleetcontrol.EntityManagementFleetDeploymentPhaseTypes.CREATED:
			fmt.Printf("⏳ [%s] Status: CREATED - Deployment is being prepared...\n", time.Now().Format("15:04:05"))
		case fleetcontrol.EntityManagementFleetDeploymentPhaseTypes.IN_PROGRESS:
			// Display ring deployment tracker info if available
			if len(deploymentEntity.RingsDeploymentTracker) > 0 {
				fmt.Printf("🔄 [%s] Status: IN_PROGRESS - Deploying across rings:\n", time.Now().Format("15:04:05"))
				for _, ring := range deploymentEntity.RingsDeploymentTracker {
					// Add status-specific emojis for each ring
					ringEmoji := "⏺️"
					if ring.Status == "COMPLETED" {
						ringEmoji = "✅"
					} else if ring.Status == "IN_PROGRESS" {
						ringEmoji = "🔄"
					} else if ring.Status == "FAILED" {
						ringEmoji = "❌"
					}
					fmt.Printf("  %s Ring '%s': %s\n", ringEmoji, ring.Name, ring.Status)
				}
			} else {
				fmt.Printf("🔄 [%s] Status: IN_PROGRESS - Deployment is in progress...\n", time.Now().Format("15:04:05"))
			}
		case fleetcontrol.EntityManagementFleetDeploymentPhaseTypes.COMPLETED:
			fmt.Printf("\n✅ Deployment COMPLETED successfully!\n")
			fmt.Printf("  📋 Deployment ID: %s\n", deploymentEntity.ID)
			fmt.Printf("  📝 Deployment Name: %s\n", deploymentEntity.Name)
			if len(deploymentEntity.RingsDeploymentTracker) > 0 {
				fmt.Printf("  🎯 Rings deployed:\n")
				for _, ring := range deploymentEntity.RingsDeploymentTracker {
					fmt.Printf("    ✅ %s: %s\n", ring.Name, ring.Status)
				}
			}
			return PrintDeploymentSuccess(deploymentEntity)
		case fleetcontrol.EntityManagementFleetDeploymentPhaseTypes.FAILED:
			fmt.Printf("\n❌ Deployment FAILED\n")
			fmt.Printf("  📋 Deployment ID: %s\n", deploymentEntity.ID)
			fmt.Printf("  📝 Deployment Name: %s\n", deploymentEntity.Name)
			if len(deploymentEntity.RingsDeploymentTracker) > 0 {
				fmt.Printf("  📊 Ring status:\n")
				for _, ring := range deploymentEntity.RingsDeploymentTracker {
					ringEmoji := "❌"
					if ring.Status == "COMPLETED" {
						ringEmoji = "✅"
					} else if ring.Status == "IN_PROGRESS" {
						ringEmoji = "🔄"
					}
					fmt.Printf("    %s %s: %s\n", ringEmoji, ring.Name, ring.Status)
				}
			}
			return PrintError(fmt.Errorf("deployment failed"))
		case fleetcontrol.EntityManagementFleetDeploymentPhaseTypes.INTERNAL_FAILURE:
			fmt.Printf("\n⚠️  Deployment encountered an INTERNAL_FAILURE\n")
			fmt.Printf("  📋 Deployment ID: %s\n", deploymentEntity.ID)
			fmt.Printf("  📝 Deployment Name: %s\n", deploymentEntity.Name)
			return PrintError(fmt.Errorf("deployment internal failure"))
		default:
			// Unknown phase - treat as complete and return
			fmt.Printf("\n🔔 Deployment reached phase: %s\n", phase)
			return PrintDeploymentSuccess(deploymentEntity)
		}

		// Show spinner while waiting before next poll
		spinner(pollInterval, "Checking deployment status...")
	}

	// Timeout reached
	return PrintError(fmt.Errorf("deployment status check timed out after %d polls", maxPolls))
}
