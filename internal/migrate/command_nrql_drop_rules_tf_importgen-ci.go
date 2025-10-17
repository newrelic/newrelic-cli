package migrate

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ciInputFile     string
	ciInputJSON     string
	ciWorkspacePath string
	useTofu         bool
	generateRemoved bool // Hidden flag for removed blocks
)

var cmdNRQLDropRulesTFImportGenCI = &cobra.Command{
	Use:   "tf-importgen-ci",
	Short: "Generate Terraform import configuration for CI environments",
	Long: `
            Generate Terraform import configuration for CI environments based on drop rules data.
            This command creates a complete Terraform workspace with import blocks and provider 
            configuration for Pipeline Cloud Rules migration in CI/CD pipelines.
    `,
	Example: `  # Generate import config from file
  newrelic migrate nrqldroprules tf-importgen-ci --file drop_rules.json

  # Generate import config from JSON string
  newrelic migrate nrqldroprules tf-importgen-ci --json '{"drop_rule_resource_ids":[...]}'

  # Generate import config in specific workspace
  newrelic migrate nrqldroprules tf-importgen-ci --file drop_rules.json --workspacePath /path/to/workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		runNRQLDropRulesTFImportGenCI()
	},
}

func init() {
	cmdNRQLDropRules.AddCommand(cmdNRQLDropRulesTFImportGenCI)

	cmdNRQLDropRulesTFImportGenCI.Flags().StringVar(&ciInputFile, "file", "", "JSON file containing drop rule resource IDs")
	cmdNRQLDropRulesTFImportGenCI.Flags().StringVar(&ciInputJSON, "json", "", "JSON string containing drop rule resource IDs")
	cmdNRQLDropRulesTFImportGenCI.Flags().StringVar(&ciWorkspacePath, "workspacePath", ".", "path to the Terraform workspace (defaults to current directory)")
	cmdNRQLDropRulesTFImportGenCI.Flags().BoolVar(&useTofu, "tofu", false, "use OpenTofu instead of Terraform")

	// Hidden flag for removed blocks functionality
	cmdNRQLDropRulesTFImportGenCI.Flags().BoolVarP(&generateRemoved, "removed", "r", false, "generate removed blocks configuration")
	_ = cmdNRQLDropRulesTFImportGenCI.Flags().MarkHidden("removed")
}

func runNRQLDropRulesTFImportGenCI() {
	// Print execution header
	printExecutionHeaderCI()

	// Validate input parameters
	if err := validateInputParametersCI(ciInputFile, ciInputJSON); err != nil {
		log.Fatal(err)
	}

	// Parse input data
	dropRuleData, err := parseDropRuleInputCI(ciInputFile, ciInputJSON)
	if err != nil {
		log.Fatalf("Failed to parse input data: %v", err)
	}

	if len(dropRuleData.DropRuleResourceIDs) == 0 {
		log.Fatal("No drop rule resource IDs found in input data")
	}

	// Check New Relic environment variables
	if err := checkNewRelicEnvironmentVariablesCI(); err != nil {
		log.Fatal(err)
	}

	// Validate account ID consistency
	validateAccountIDConsistencyCI(dropRuleData)

	// Create tool configuration
	toolConfig := newToolConfigCI(useTofu)

	// Resolve workspace path
	absWorkspacePath, err := resolveWorkspacePathCI(ciWorkspacePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Using %s workspace: %s", toolConfig.DisplayName, absWorkspacePath)

	// Check tool prerequisites
	if err := checkTerraformPrerequisitesCI(toolConfig); err != nil {
		log.Fatal(err)
	}

	// Validate workspace is empty or suitable for initialization
	if err := validateWorkspaceCI(absWorkspacePath); err != nil {
		log.Fatal(err)
	}

	// Generate configuration files
	if err := generateAllConfigFilesCI(toolConfig, absWorkspacePath, dropRuleData); err != nil {
		log.Fatalf("Failed to generate %s files: %v", toolConfig.DisplayName, err)
	}

	// Initialize workspace
	if err := initializeWorkspaceCI(toolConfig, absWorkspacePath); err != nil {
		log.Fatalf("Failed to initialize %s: %v", toolConfig.DisplayName, err)
	}

	// Generate configuration plan
	if err := generateConfigurationPlanCI(toolConfig, absWorkspacePath); err != nil {
		log.Fatalf("Failed to plan %s: %v", toolConfig.DisplayName, err)
	}

	// Generate removed blocks only if flag is set
	if generateRemoved {
		if err := generateRemovedBlocksCI(absWorkspacePath, dropRuleData); err != nil {
			log.Fatalf("Failed to generate removed blocks: %v", err)
		}
	}

	// Format configuration files
	if err := formatConfigurationCI(toolConfig, absWorkspacePath); err != nil {
		log.Warnf("Failed to format configuration: %v", err)
	}

	// Print final recommendations
	printCIRecommendationsCI(toolConfig, dropRuleData, absWorkspacePath, generateRemoved)
}
