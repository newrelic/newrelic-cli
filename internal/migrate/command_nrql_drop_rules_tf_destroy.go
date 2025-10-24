package migrate

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	delistWorkspacePath        string
	delistResourceIdentifiers  []string
	delistSkipResponseToPrompt bool
	delistUseTofu              bool
)

var cmdNRQLDropRulesTFDestroy = &cobra.Command{
	Use:   "tf-delist",
	Short: "Delist NRQL drop rules from Terraform state (without destroying)",
	Long: `
			Safely remove NRQL drop rules from your Terraform state without destroying the actual
			resources. This command uses 'terraform state rm' to delist the resources, allowing
			you to stop managing them via Terraform while keeping the drop rules active in New Relic.
			
			⚠️  IMPORTANT: This command DOES NOT destroy the actual drop rules in New Relic.
			It only removes them from Terraform state management.
	`,
	Example: `  # Delist drop rules in current directory
  newrelic migrate nrqldroprules tf-delist

  # Delist drop rules in specific workspace with OpenTofu
  newrelic migrate nrqldroprules tf-delist --workspacePath /path/to/terraform --tofu

  # Delist specific resources without prompts
  newrelic migrate nrqldroprules tf-delist --resourceIdentifiers resource1,resource2 --skipResponseToPrompt`,
	Run: func(cmd *cobra.Command, args []string) {
		runNRQLDropRulesTFDelist()
	},
}

func init() {
	cmdNRQLDropRules.AddCommand(cmdNRQLDropRulesTFDestroy)

	cmdNRQLDropRulesTFDestroy.Flags().StringVar(&delistWorkspacePath, "workspacePath", ".", "path to the Terraform workspace")
	cmdNRQLDropRulesTFDestroy.Flags().StringSliceVar(&delistResourceIdentifiers, "resourceIdentifiers", []string{}, "list of resource identifiers for newrelic_nrql_drop_rule resources")
	cmdNRQLDropRulesTFDestroy.Flags().BoolVar(&delistSkipResponseToPrompt, "skipResponseToPrompt", false, "skip all user prompts (answers 'N' to all prompts)")
	cmdNRQLDropRulesTFDestroy.Flags().BoolVar(&delistUseTofu, "tofu", false, "use OpenTofu instead of Terraform")
}

func runNRQLDropRulesTFDelist() {
	// Create command context with new CommandDelist type
	ctx, err := NewCommandContext(delistUseTofu, delistWorkspacePath, delistSkipResponseToPrompt, CommandDelist, delistResourceIdentifiers)
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
		handleDelistStateFailure(ctx)
		return
	}

	if len(dropRuleResources) > 0 {
		handleDelistStateSuccess(ctx, dropRuleResources)
	} else {
		handleDelistStateFailure(ctx)
	}
}

func handleDelistStateSuccess(ctx *CommandContext, resources []string) {
	fmt.Printf("\n✅ Found %d NRQL drop rule resources in %s state\n", len(resources), ctx.ToolConfig.DisplayName)

	// Display important warning
	fmt.Println("\n" + strings.Repeat("⚠️", 20))
	fmt.Printf("🛡️  SAFE DELISTING MODE: Resources will be REMOVED FROM STATE ONLY\n")
	fmt.Printf("📋 The actual drop rules in New Relic will remain ACTIVE and UNCHANGED\n")
	fmt.Printf("🔄 This allows you to stop managing them via %s safely\n", ctx.ToolConfig.DisplayName)
	fmt.Println(strings.Repeat("⚠️", 20))

	// Check and cleanup import configuration file before proceeding
	checkAndCleanupImportConfig(ctx.WorkspacePath)

	// List found resources
	fmt.Println("\n📋 Resources to be delisted from state:")
	for i, resource := range resources {
		fmt.Printf("  %d. %s\n", i+1, resource)
	}

	// Generate and display state rm commands
	stateRmCommands := generateStateRmCommands(ctx.ToolConfig.ToolName, resources)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("📝 Generated %s state removal commands:\n", ctx.ToolConfig.DisplayName)
	fmt.Println(strings.Repeat("=", 60))

	for i, cmd := range stateRmCommands {
		fmt.Printf("%d. %s\n", i+1, cmd)
	}

	// Execute if user confirms
	fmt.Println("\n" + strings.Repeat("=", 60))
	if ctx.PromptForExecution() {
		if err := executeDelistFlow(ctx, stateRmCommands, resources); err != nil {
			fmt.Printf("\n❌ Delisting failed: %v\n", err)
			log.Error(err)
		} else {
			fmt.Println("\n🎉 NRQL drop rule delisting completed successfully!")
			printPostDelistInstructions(ctx.ToolConfig.DisplayName, resources)
		}
	} else {
		fmt.Printf("\n🛑 Execution halted by user.\n")
		fmt.Printf("Please run the commands above manually in your %s workspace.\n", ctx.ToolConfig.DisplayName)
		printPostDelistInstructions(ctx.ToolConfig.DisplayName, resources)
	}
}

func handleDelistStateFailure(ctx *CommandContext) {
	if len(ctx.ResourceIDs) == 0 {
		log.Fatalf("Unable to list %s state and no --resourceIdentifiers provided. Please specify resource identifiers for newrelic_nrql_drop_rule resources.", ctx.ToolConfig.DisplayName)
	}

	log.Infof("Using provided resource identifiers: %v", ctx.ResourceIDs)

	// Check and cleanup import configuration file before proceeding
	checkAndCleanupImportConfig(ctx.WorkspacePath)

	// Generate state rm commands for provided resources
	stateRmCommands := generateStateRmCommands(ctx.ToolConfig.ToolName, ctx.ResourceIDs)

	fmt.Printf("\nGenerated %s state removal commands for provided resources:\n", ctx.ToolConfig.DisplayName)
	for i, cmd := range stateRmCommands {
		fmt.Printf("%d. %s\n", i+1, cmd)
	}
	fmt.Printf("\nPlease run these commands in your %s workspace.\n", ctx.ToolConfig.DisplayName)
	printPostDelistInstructions(ctx.ToolConfig.DisplayName, ctx.ResourceIDs)
}
