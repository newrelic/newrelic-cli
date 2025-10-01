package migrate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	importWorkspacePath        string
	pipelineCloudRuleIDs       []string
	fileName                   string
	importSkipResponseToPrompt bool
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

  # Generate import config in specific workspace and save to file
  newrelic migrate nrqldroprules tf-importgen --workspacePath /path/to/terraform --fileName imports.tf

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
}

type TerraformState struct {
	Values struct {
		RootModule struct {
			Resources []struct {
				Address string                 `json:"address"`
				Values  map[string]interface{} `json:"values"`
			} `json:"resources"`
		} `json:"root_module"`
	} `json:"values"`
}

func runNRQLDropRulesTFImportGen() {
	// Resolve workspace path
	absWorkspacePath, err := filepath.Abs(importWorkspacePath)
	if err != nil {
		log.Fatalf("Error resolving workspace path: %v", err)
	}

	log.Infof("Using Terraform workspace: %s", absWorkspacePath)

	// Try to run terraform state list
	dropRuleResources, err := getTerraformDropRuleResources(absWorkspacePath)
	if err != nil {
		log.Warnf("Failed to list Terraform state: %v", err)
		handleImportStateFailure()
		return
	}

	if len(dropRuleResources) > 0 {
		handleImportStateSuccess(absWorkspacePath, dropRuleResources)
	} else {
		handleImportStateFailure()
	}
}

func handleImportStateSuccess(workspacePath string, resources []string) {
	log.Infof("Found %d NRQL drop rule resources in Terraform state", len(resources))

	// Validate that all resources have pipeline_cloud_rule_entity_id
	_, err := validateAndExtractPipelineRuleIDs(workspacePath, resources)
	if err != nil {
		log.Fatalf("Validation failed: %v", err)
	}

	// Check Terraform version before proceeding
	log.Info("Checking Terraform version...")
	terraformVersion, err := getTerraformVersion()
	if err != nil {
		log.Warnf("Could not determine Terraform version: %v", err)
		log.Info("Skipping Terraform version check. Note that Terraform >= 1.5 is required for generating import configuration.")
	} else {
		log.Infof("Detected Terraform version: %s", terraformVersion)

		// Parse and validate Terraform version
		if !isValidTerraformVersion(terraformVersion, "1.5") {
			log.Fatalf("This command requires Terraform version >= 1.5 to generate import configuration. Your version: %s", terraformVersion)
		}
		log.Info("Terraform version check passed.")
	}
	// Generate import configuration
	importConfig := generateImportConfigurationFromResources(resources)

	// Handle output based on fileName flag
	if fileName != "" {
		writeImportConfigToFile(workspacePath, importConfig)
	} else {
		fmt.Printf("\nGenerated import configuration:\n")
		fmt.Println(importConfig)
	}

	// Generate config file name for terraform command
	generateConfigOutFile := "generated_pipeline_rules.tf"
	importConfigFile := "import_config_pipeline_rules.tf"

	// Ask if user wants to write import config to the generated file (only if fileName not specified)
	if fileName == "" && !importSkipResponseToPrompt {
		if promptForWriteToFile() {
			writeImportConfigToGeneratedFile(workspacePath, importConfig, importConfigFile)
		}
	}

	// Generate and show terraform commands
	planCommand := fmt.Sprintf("terraform plan -generate-config-out=%s", generateConfigOutFile)
	applyCommand := "terraform apply"

	fmt.Printf("\nGenerated Terraform commands:\n")
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", applyCommand)

	// Ask user if they want to execute
	if importSkipResponseToPrompt {
		fmt.Println("\nSkipping execution due to --skipResponseToPrompt flag")
		return
	}

	if !promptForImportExecution() {
		fmt.Println("\nExecution halted. Please run the commands above manually in your Terraform workspace.")
		return
	}

	// Execute commands
	executeImportTerraformCommands(workspacePath, planCommand, applyCommand)
}

func handleImportStateFailure() {
	if len(pipelineCloudRuleIDs) == 0 {
		log.Fatal("Unable to list Terraform state and no --pipelineCloudRuleIDs provided. Please specify Pipeline Cloud Rule IDs to generate import configuration.")
	}

	log.Infof("Using provided Pipeline Cloud Rule IDs: %v", pipelineCloudRuleIDs)

	// Generate import configuration from provided IDs
	importConfig := generateImportConfiguration(pipelineCloudRuleIDs)

	// Handle output based on fileName flag
	if fileName != "" {
		writeImportConfigToFile(".", importConfig)
	} else {
		fmt.Printf("\nGenerated import configuration for provided Pipeline Cloud Rule IDs:\n")
		fmt.Println(importConfig)
	}

	// Generate config file name for terraform command
	generateConfigOutFile := "generated_pipeline_rules.tf"
	importConfigFile := "import_config_pipeline_rules.tf"

	// Ask if user wants to write import config to the generated file (only if fileName not specified)
	if fileName == "" && !importSkipResponseToPrompt {
		if promptForWriteToFile() {
			writeImportConfigToGeneratedFile(".", importConfig, importConfigFile)
		}
	}

	// Generate and show terraform commands
	planCommand := fmt.Sprintf("terraform plan -generate-config-out=%s", generateConfigOutFile)
	applyCommand := "terraform apply"

	fmt.Printf("\nGenerated Terraform commands:\n")
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", applyCommand)

	if !importSkipResponseToPrompt {
		if promptForImportExecution() {
			executeImportTerraformCommands(".", planCommand, applyCommand)
		} else {
			fmt.Println("\nExecution halted. Please run the commands above manually in your Terraform workspace.")
		}
	} else {
		fmt.Println("\nPlease run these commands in your appropriate Terraform workspace.")
	}
}

func validateAndExtractPipelineRuleIDs(workspacePath string, resources []string) ([]string, error) {
	var pipelineRuleIDs []string

	for _, resource := range resources {
		// Get detailed state information using terraform show -json
		cmd := exec.Command("terraform", "show", "-json")
		cmd.Dir = workspacePath

		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get terraform state JSON: %v", err)
		}

		var state TerraformState
		if err := json.Unmarshal(output, &state); err != nil {
			return nil, fmt.Errorf("failed to parse terraform state JSON: %v", err)
		}

		// Find the specific resource and extract pipeline_cloud_rule_entity_id
		found := false
		for _, res := range state.Values.RootModule.Resources {
			if res.Address == resource {
				found = true
				if entityID, ok := res.Values["pipeline_cloud_rule_entity_id"]; ok && entityID != nil && entityID != "" {
					pipelineRuleIDs = append(pipelineRuleIDs, fmt.Sprintf("%v", entityID))
				} else {
					return nil, fmt.Errorf("resource %s is missing pipeline_cloud_rule_entity_id. Please run tf-update first", resource)
				}
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("resource %s not found in state", resource)
		}
	}

	return pipelineRuleIDs, nil
}

func generateImportConfigurationFromResources(resources []string) string {
	var importBlocks []string

	for _, resource := range resources {
		// Extract the resource identifier (e.g., "foo" from "newrelic_nrql_drop_rule.foo")
		resourceParts := strings.Split(resource, ".")
		if len(resourceParts) < 2 {
			log.Warnf("Invalid resource format: %s", resource)
			continue
		}

		resourceIdentifier := strings.Join(resourceParts[1:], ".") // Handle cases like "module.something.resource.name"

		importBlock := fmt.Sprintf(`import {
  to = newrelic_pipeline_cloud_rule.%s
  id = %s.pipeline_cloud_rule_entity_id
}`, resourceIdentifier, resource)
		importBlocks = append(importBlocks, importBlock)
	}

	return strings.Join(importBlocks, "\n\n")
}

func generateImportConfiguration(pipelineRuleIDs []string) string {
	var importBlocks []string

	for i, ruleID := range pipelineRuleIDs {
		importBlock := fmt.Sprintf(`import {
  to = newrelic_pipeline_cloud_rule.pipeline_rule_%d
  id = "%s"
}`, i+1, ruleID)
		importBlocks = append(importBlocks, importBlock)
	}

	return strings.Join(importBlocks, "\n\n")
}

func writeImportConfigToFile(workspacePath, importConfig string) {
	filePath := filepath.Join(workspacePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		log.Errorf("Failed to create file %s: %v", filePath, err)
		fmt.Printf("\nFailed to write to file. Import configuration:\n")
		fmt.Println(importConfig)
		return
	}
	defer file.Close()

	_, err = file.WriteString(importConfig)
	if err != nil {
		log.Errorf("Failed to write to file %s: %v", filePath, err)
		fmt.Printf("\nFailed to write to file. Import configuration:\n")
		fmt.Println(importConfig)
		return
	}

	fmt.Printf("\nImport configuration written to: %s\n", filePath)
}

func promptForImportExecution() bool {
	fmt.Print("\nWould you like this CLI to execute the commands above? (Y/N): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Errorf("Error reading input: %v", err)
		return false
	}

	response = strings.TrimSpace(strings.ToUpper(response))
	return response == "Y" || response == "YES"
}

func promptForWriteToFile() bool {
	fmt.Print("\nWould you like to write the import configuration to the generated config file? (Y/N): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Errorf("Error reading input: %v", err)
		return false
	}

	response = strings.TrimSpace(strings.ToUpper(response))
	return response == "Y" || response == "YES"
}

func executeImportTerraformCommands(workspacePath, planCommand, applyCommand string) {
	// Execute plan command
	fmt.Println("\nExecuting terraform plan...")
	if err := executeTerraformCommand(workspacePath, planCommand); err != nil {
		log.Errorf("Terraform plan failed: %v", err)
		return
	}

	// Ask for confirmation before apply
	fmt.Print("\nProceed with terraform apply? (Y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Errorf("Error reading input: %v", err)
		return
	}

	response = strings.TrimSpace(strings.ToUpper(response))
	if response != "Y" && response != "YES" {
		fmt.Println("Apply cancelled.")
		return
	}

	// Execute apply command
	fmt.Println("\nExecuting terraform apply...")
	if err := executeTerraformCommand(workspacePath, applyCommand); err != nil {
		log.Errorf("Terraform apply failed: %v", err)
		return
	}

	fmt.Println("\n✅ Import configuration and terraform commands executed successfully!")
}

func writeImportConfigToGeneratedFile(workspacePath, importConfig, importConfigFile string) {
	filePath := filepath.Join(workspacePath, importConfigFile)

	file, err := os.Create(filePath)
	if err != nil {
		log.Errorf("Failed to create file %s: %v", filePath, err)
		fmt.Printf("\nFailed to write import configuration to file.\n")
		return
	}
	defer file.Close()

	_, err = file.WriteString(importConfig)
	if err != nil {
		log.Errorf("Failed to write to file %s: %v", filePath, err)
		fmt.Printf("\nFailed to write import configuration to file.\n")
		return
	}

	fmt.Printf("\nImport configuration written to: %s\n", filePath)
}

func isValidTerraformVersion(currentVersion, minVersion string) bool {
	// Split versions by dots and remove any "v" prefix
	currentVersion = strings.TrimPrefix(currentVersion, "v")
	minVersion = strings.TrimPrefix(minVersion, "v")

	currentParts := strings.Split(currentVersion, ".")
	minParts := strings.Split(minVersion, ".")

	// Ensure we have at least two parts for comparison
	for i := 0; i < 2; i++ {
		var currentVal, minVal int

		// Get current value, default to 0 if not enough parts
		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentVal)
		}

		// Get min value, default to 0 if not enough parts
		if i < len(minParts) {
			fmt.Sscanf(minParts[i], "%d", &minVal)
		}

		if currentVal >= minVal {
			// Current version component is greater, so overall version is higher
			if i < 1 {
				continue
			}
			return true
		} else if currentVal < minVal {
			// Current version component is less, so overall version is lower
			return false
		}

		// If equal, continue to next component
	}

	// If we've compared all components and they're equal, versions are equal
	return true

}
