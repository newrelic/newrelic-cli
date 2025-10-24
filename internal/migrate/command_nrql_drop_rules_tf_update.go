package migrate

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	workspacePath        string
	resourceIdentifiers  []string
	skipResponseToPrompt bool
	updateUseTofu        bool
)

var cmdNRQLDropRulesTFUpdate = &cobra.Command{
	Use:   "tf-update",
	Short: "Update NRQL drop rules in Terraform workspace",
	Long: `
			Update NRQL drop rules in your Terraform workspace by refreshing the state
			to include pipeline_cloud_rule_entity_id values. This command will attempt to run
			terraform commands in your workspace to update the drop rule resources.
	`,
	Example: `  # Update drop rules in current directory
  newrelic migrate nrqldroprules tf-update

  # Update drop rules in specific workspace with OpenTofu
  newrelic migrate nrqldroprules tf-update --workspacePath /path/to/terraform --tofu

  # Update specific resources without prompts
  newrelic migrate nrqldroprules tf-update --resourceIdentifiers resource1,resource2 --skipResponseToPrompt`,
	Run: func(cmd *cobra.Command, args []string) {
		runNRQLDropRulesTFUpdate()
	},
}

func init() {
	cmdNRQLDropRules.AddCommand(cmdNRQLDropRulesTFUpdate)

	cmdNRQLDropRulesTFUpdate.Flags().StringVar(&workspacePath, "workspacePath", ".", "path to the Terraform workspace")
	cmdNRQLDropRulesTFUpdate.Flags().StringSliceVar(&resourceIdentifiers, "resourceIdentifiers", []string{}, "list of resource identifiers for newrelic_nrql_drop_rule resources")
	cmdNRQLDropRulesTFUpdate.Flags().BoolVar(&skipResponseToPrompt, "skipResponseToPrompt", false, "skip all user prompts (answers 'N' to all prompts)")
	cmdNRQLDropRulesTFUpdate.Flags().BoolVar(&updateUseTofu, "tofu", false, "use OpenTofu instead of Terraform")
}

func runNRQLDropRulesTFUpdate() {
	// Create command context
	ctx, err := NewCommandContext(updateUseTofu, workspacePath, skipResponseToPrompt, CommandUpdate, resourceIdentifiers)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize command
	if err := ctx.InitializeCommand(); err != nil {
		log.Fatal(err)
	}

	// Try to get drop rule resources from state
	dropRuleResources, err := ctx.GetDropRuleResources()
	if err != nil {
		log.Warnf("Failed to list %s state: %v", ctx.ToolConfig.DisplayName, err)
		handleUpdateStateFailure(ctx)
		return
	}

	if len(dropRuleResources) > 0 {
		handleUpdateStateSuccess(ctx, dropRuleResources)
	} else {
		handleUpdateStateFailure(ctx)
	}
}

func handleUpdateStateSuccess(ctx *CommandContext, resources []string) {
	fmt.Printf("\n✅ Found %d NRQL drop rule resources in %s state\n", len(resources), ctx.ToolConfig.DisplayName)

	// List found resources
	fmt.Println("\n📋 Resources to be updated:")
	for i, resource := range resources {
		fmt.Printf("  %d. %s\n", i+1, resource)
	}

	// Check provider version
	fmt.Println("\n🔍 Checking provider version requirements...")
	if err := ctx.CheckProviderVersion(); err != nil {
		log.Fatal(err)
	}

	// Generate and display commands
	planCmd, actionCmd := ctx.GenerateTargetCommands(resources)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("📝 Generated %s commands for resource updates:\n", ctx.ToolConfig.DisplayName)
	fmt.Println(strings.Repeat("=", 60))
	ctx.PrintCommands(planCmd, actionCmd)

	// Execute if user confirms
	fmt.Println("\n" + strings.Repeat("=", 60))
	if ctx.PromptForExecution() {
		if err := ctx.ExecuteStandardFlow(planCmd, actionCmd, resources); err != nil {
			fmt.Printf("\n❌ Execution failed: %v\n", err)
			log.Error(err)
		} else {
			fmt.Println("\n🎉 NRQL drop rule update completed successfully!")
			fmt.Println("Your resources now include pipeline_cloud_rule_entity_id values.")
		}
	} else {
		fmt.Printf("\n🛑 Execution halted by user.\n")
		fmt.Printf("Please run the commands above manually in your %s workspace:\n", ctx.ToolConfig.DisplayName)
		fmt.Printf("  1. %s\n", planCmd)
		fmt.Printf("  2. %s -auto-approve\n", actionCmd)
		fmt.Println("\nNote: The -auto-approve flag is recommended to avoid interactive prompts.")
	}
}

func handleUpdateStateFailure(ctx *CommandContext) {
	if len(ctx.ResourceIDs) == 0 {
		log.Fatalf("Unable to list %s state and no --resourceIdentifiers provided. Please specify resource identifiers for newrelic_nrql_drop_rule resources.", ctx.ToolConfig.DisplayName)
	}

	log.Infof("Using provided resource identifiers: %v", ctx.ResourceIDs)
	ctx.PrintCommandsForResourceIDs()
}
