package migrate

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	importWorkspacePath        string
	pipelineCloudRuleIDs       []string
	fileName                   string
	importSkipResponseToPrompt bool
	importUseTofu              bool
)

var cmdNRQLDropRulesTFImportGen = &cobra.Command{
	Use:   "tf-importgen",
	Short: "Generate Terraform import configuration for Pipeline Cloud Rules",
	Long: `
			Generate Terraform import blocks for Pipeline Cloud Rules based on existing
			NRQL drop rules in your Terraform workspace. This command validates that
			drop rules have been updated with pipeline_cloud_rule_entity_id values
			and generates the necessary import configuration.
	`,
	Example: `  # Generate import config in current directory
  newrelic migrate nrqldroprules tf-importgen

  # Generate import config in specific workspace and save to file with OpenTofu
  newrelic migrate nrqldroprules tf-importgen --workspacePath /path/to/terraform --fileName imports.tf --tofu

  # Generate import config with specific Pipeline Cloud Rule IDs
  newrelic migrate nrqldroprules tf-importgen --pipelineCloudRuleIDs id1,id2 --skipResponseToPrompt`,
	Run: func(cmd *cobra.Command, args []string) {
		runNRQLDropRulesTFImportGen()
	},
}

func init() {
	cmdNRQLDropRules.AddCommand(cmdNRQLDropRulesTFImportGen)

	cmdNRQLDropRulesTFImportGen.Flags().StringVar(&importWorkspacePath, "workspacePath", ".", "path to the Terraform workspace")
	cmdNRQLDropRulesTFImportGen.Flags().StringSliceVar(&pipelineCloudRuleIDs, "pipelineCloudRuleIDs", []string{}, "list of Pipeline Cloud Rule IDs to generate import configuration with")
	cmdNRQLDropRulesTFImportGen.Flags().StringVar(&fileName, "fileName", "", "file name to write the import blocks to (prints to terminal if not specified)")
	cmdNRQLDropRulesTFImportGen.Flags().BoolVar(&importSkipResponseToPrompt, "skipResponseToPrompt", false, "skip all user prompts (answers 'N' to all prompts)")
	cmdNRQLDropRulesTFImportGen.Flags().BoolVar(&importUseTofu, "tofu", false, "use OpenTofu instead of Terraform")
}

func runNRQLDropRulesTFImportGen() {
	// Create command context
	ctx, err := NewCommandContext(importUseTofu, importWorkspacePath, importSkipResponseToPrompt, CommandImport, pipelineCloudRuleIDs)
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
		handleImportStateFailure(ctx)
		return
	}

	if len(dropRuleResources) > 0 {
		handleImportStateSuccess(ctx, dropRuleResources)
	} else {
		handleImportStateFailure(ctx)
	}
}

func handleImportStateSuccess(ctx *CommandContext, resources []string) {
	log.Infof("Found %d NRQL drop rule resources in %s state", len(resources), ctx.ToolConfig.DisplayName)

	// Validate resources have pipeline_cloud_rule_entity_id
	_, err := ValidateAndExtractPipelineRuleIDs(ctx.WorkspacePath, resources, ctx.ToolConfig)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	// Check tool version
	checkToolVersionForImport(ctx.ToolConfig)

	// Generate and handle import configuration
	importConfig := GenerateImportConfigFromResources(resources)
	handleImportOutput(ctx, importConfig)
	executeImportWorkflow(ctx, importConfig)
}

func handleImportStateFailure(ctx *CommandContext) {
	if len(pipelineCloudRuleIDs) == 0 {
		log.Fatalf("Unable to list %s state and no --pipelineCloudRuleIDs provided. Please specify Pipeline Cloud Rule IDs to generate import configuration.", ctx.ToolConfig.DisplayName)
	}

	log.Infof("Using provided Pipeline Cloud Rule IDs: %v", pipelineCloudRuleIDs)

	// Generate import configuration from provided IDs
	importConfig := GenerateImportConfigFromIDs(pipelineCloudRuleIDs)
	handleImportOutput(ctx, importConfig)
	executeImportWorkflow(ctx, importConfig)
}

func checkToolVersionForImport(config *ToolConfig) {
	log.Infof("Checking %s version...", config.DisplayName)
	toolVersion, err := getToolVersionCI(config)
	if err != nil {
		log.Warnf("Could not determine %s version: %v", config.DisplayName, err)
		log.Infof("Skipping %s version check. Note that %s >= 1.5 is required for generating import configuration.", config.DisplayName, config.DisplayName)
		return
	}

	log.Infof("Detected %s version: %s", config.DisplayName, toolVersion)

	if !isValidVersionCI(toolVersion, MinRequiredVersion) {
		log.Fatalf("This command requires %s version >= %s to generate import configuration. Your version: %s", config.DisplayName, MinRequiredVersion, toolVersion)
	}
	log.Infof("%s version check passed.", config.DisplayName)
}

func handleImportOutput(ctx *CommandContext, importConfig string) {
	if fileName != "" {
		WriteConfigToFile(ctx.WorkspacePath, importConfig, fileName)
	} else {
		fmt.Printf("\nGenerated import configuration:\n")
		fmt.Println(importConfig)
	}
}

func executeImportWorkflow(ctx *CommandContext, importConfig string) {
	// Handle writing to generated file
	if fileName == "" && !ctx.SkipPrompts {
		fmt.Print("\nWould you like to write the import configuration to the generated config file? (Y/N): ")
		if readUserInput() {
			WriteConfigToFile(ctx.WorkspacePath, importConfig, "import_config_pipeline_rules.tf")
		}
	}

	// Generate and show commands
	planCommand := fmt.Sprintf("%s plan -generate-config-out=generated_pipeline_rules.tf", ctx.ToolConfig.ToolName)
	applyCommand := fmt.Sprintf("%s apply", ctx.ToolConfig.ToolName)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("📝 Generated %s commands for Pipeline Cloud Rule import:\n", ctx.ToolConfig.DisplayName)
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", applyCommand)

	// Execute if user confirms
	fmt.Println("\n" + strings.Repeat("=", 60))
	if ctx.PromptForExecution() {
		if err := executeImportCommandsStandardFlow(ctx, planCommand, applyCommand); err != nil {
			fmt.Printf("\n❌ Execution failed: %v\n", err)
			log.Error(err)
		} else {
			fmt.Println("\n🎉 Pipeline Cloud Rule import completed successfully!")
			fmt.Println("Your Pipeline Cloud Rules have been imported into Terraform state.")
		}
	} else {
		fmt.Printf("\n🛑 Execution halted by user.\n")
		fmt.Printf("Please run the commands above manually in your %s workspace:\n", ctx.ToolConfig.DisplayName)
		fmt.Printf("  1. %s\n", planCommand)
		fmt.Printf("  2. %s -auto-approve\n", applyCommand)
		fmt.Println("\nNote: The -auto-approve flag is recommended to avoid interactive prompts.")
	}
}

func executeImportCommandsStandardFlow(ctx *CommandContext, planCmd, applyCmd string) error {
	// Execute plan and show output immediately
	if err := ctx.ExecuteCommand(planCmd, "plan"); err != nil {
		return err
	}

	// Show plan completion and ask for apply confirmation
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("📋 Plan completed successfully! Review the import plan above.\n")
	fmt.Println(strings.Repeat("=", 60))

	// Get confirmation for apply
	if !ctx.PromptForActionConfirmation() {
		fmt.Println("\n🛑 Import cancelled by user.")
		return nil
	}

	// Add auto-approve to apply command and show warning
	autoApproveCmd := applyCmd + " -auto-approve"

	fmt.Printf("\n⚠️  WARNING: Using -auto-approve flag to avoid interactive prompts.\n")
	fmt.Printf("The import operation will proceed without additional confirmation.\n")
	fmt.Println(strings.Repeat("-", 60))

	// Execute apply with auto-approve
	if err := ctx.ExecuteCommand(autoApproveCmd, "apply"); err != nil {
		return err
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🔍 Import operation completed successfully!\n")
	fmt.Println(strings.Repeat("=", 60))

	return nil
}
