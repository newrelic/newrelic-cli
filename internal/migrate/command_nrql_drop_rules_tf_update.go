package migrate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	workspacePath        string
	resourceIdentifiers  []string
	skipResponseToPrompt bool
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

  # Update drop rules in specific workspace
  newrelic migrate nrqldroprules tf-update --workspacePath /path/to/terraform

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
}

func runNRQLDropRulesTFUpdate() {
	// Resolve workspace path
	absWorkspacePath, err := filepath.Abs(workspacePath)
	if err != nil {
		log.Fatalf("Error resolving workspace path: %v", err)
	}

	log.Infof("Using Terraform workspace: %s", absWorkspacePath)

	// Try to run terraform state list
	dropRuleResources, err := getTerraformDropRuleResources(absWorkspacePath)
	if err != nil {
		log.Warnf("Failed to list Terraform state: %v", err)
		handleTerraformStateFailure()
		return
	}

	if len(dropRuleResources) > 0 {
		handleTerraformStateSuccess(absWorkspacePath, dropRuleResources)
	} else {
		handleTerraformStateFailure()
	}
}

func getTerraformDropRuleResources(workspacePath string) ([]string, error) {
	cmd := exec.Command("terraform", "state", "list")
	cmd.Dir = workspacePath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("terraform state list failed: %v", err)
	}

	var dropRuleResources []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "newrelic_nrql_drop_rule") {
			dropRuleResources = append(dropRuleResources, line)
		}
	}

	return dropRuleResources, nil
}

func handleTerraformStateSuccess(workspacePath string, resources []string) {
	log.Infof("Found %d NRQL drop rule resources in Terraform state", len(resources))

	// Check New Relic Terraform Provider version before proceeding
	log.Info("Checking New Relic Terraform Provider version...")
	providerVersion, err := getNewRelicProviderVersion(workspacePath)
	if err != nil {
		log.Warnf("Could not determine New Relic Terraform Provider version: %v", err)
		log.Info("Skipping provider version check. Note that New Relic Terraform Provider >= 3.63.0 is required for pipeline_cloud_rule_entity_id support.")
	} else {
		log.Infof("Detected New Relic Terraform Provider version: %s", providerVersion)

		// Parse and validate provider version using go-version
		if !isValidProviderVersion(providerVersion, "3.63.0") {
			log.Fatalf("Changes to add pipeline_cloud_rule_entity_id corresponding to drop rules would not be added to the state with New Relic Terraform Provider version %s. Provider version >= 3.63.0 is required. Please upgrade your provider.", providerVersion)
		}
		log.Info("New Relic Terraform Provider version check passed.")
	}

	// Generate target flags
	targetFlags := make([]string, len(resources))
	for i, resource := range resources {
		targetFlags[i] = fmt.Sprintf("-target=%s", resource)
	}
	targetString := strings.Join(targetFlags, " ")

	// Generate commands
	planCommand := fmt.Sprintf("terraform plan -refresh-only %s", targetString)
	applyCommand := fmt.Sprintf("terraform apply -refresh-only %s", targetString)

	// Print commands
	fmt.Printf("\nGenerated Terraform commands:\n")
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", applyCommand)

	// Ask user if they want to execute
	if skipResponseToPrompt {
		fmt.Println("\nSkipping execution due to --skipResponseToPrompt flag")
		return
	}

	if !promptForExecution() {
		fmt.Println("\nExecution halted. Please run the commands above manually in your Terraform workspace.")
		return
	}

	// Execute commands
	executeTerraformCommands(workspacePath, planCommand, applyCommand, resources)
}

func handleTerraformStateFailure() {
	if len(resourceIdentifiers) == 0 {
		log.Fatal("Unable to list Terraform state and no --resourceIdentifiers provided. Please specify resource identifiers for newrelic_nrql_drop_rule resources.")
	}

	log.Infof("Using provided resource identifiers: %v", resourceIdentifiers)

	// Generate target flags from provided identifiers
	targetFlags := make([]string, len(resourceIdentifiers))
	for i, resource := range resourceIdentifiers {
		targetFlags[i] = fmt.Sprintf("-target=%s", resource)
	}
	targetString := strings.Join(targetFlags, " ")

	// Generate and print commands
	planCommand := fmt.Sprintf("terraform plan -refresh-only %s", targetString)
	applyCommand := fmt.Sprintf("terraform apply -refresh-only %s", targetString)

	fmt.Printf("\nGenerated Terraform commands for provided resources:\n")
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", applyCommand)
	fmt.Println("\nPlease run these commands in your appropriate Terraform workspace.")
}

func promptForExecution() bool {
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

func executeTerraformCommands(workspacePath, planCommand, applyCommand string, resources []string) {
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

	// Validate updates
	validateDropRuleUpdates(workspacePath, resources)
}

func executeTerraformCommand(workspacePath, command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = workspacePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func validateDropRuleUpdates(workspacePath string, resources []string) {
	fmt.Println("\nValidating drop rule updates...")

	updatedCount := 0
	for _, resource := range resources {
		if hasEntityID := checkResourceHasEntityID(workspacePath, resource); hasEntityID {
			updatedCount++
			fmt.Printf("✓ %s: Updated with pipeline_cloud_rule_entity_id\n", resource)
		} else {
			fmt.Printf("⚠ %s: Missing pipeline_cloud_rule_entity_id\n", resource)
		}
	}

	if updatedCount == len(resources) {
		fmt.Printf("\n✅ All %d NRQL drop rule resources have been successfully updated with pipeline_cloud_rule_entity_id\n", updatedCount)
	} else {
		fmt.Printf("\n⚠️ %d out of %d resources were updated. Please check the remaining resources manually.\n", updatedCount, len(resources))
	}
}

func checkResourceHasEntityID(workspacePath, resource string) bool {
	cmd := exec.Command("terraform", "state", "show", resource)
	cmd.Dir = workspacePath

	output, err := cmd.Output()
	if err != nil {
		log.Warnf("Failed to show state for %s: %v", resource, err)
		return false
	}

	return strings.Contains(string(output), "pipeline_cloud_rule_entity_id")
}

// Helper functions for New Relic Terraform Provider version checking
func getNewRelicProviderVersion(workspacePath string) (string, error) {
	cmd := exec.Command("terraform", "version", "-json")
	cmd.Dir = workspacePath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get terraform version: %v", err)
	}

	// Parse JSON output
	var versionData map[string]interface{}
	if err := json.Unmarshal(output, &versionData); err != nil {
		return "", fmt.Errorf("failed to parse terraform version JSON: %v", err)
	}

	// Extract provider selections
	if providerSelections, ok := versionData["provider_selections"].(map[string]interface{}); ok {
		// Look for New Relic provider
		for providerKey, versionFound := range providerSelections {
			if strings.Contains(providerKey, "newrelic/newrelic") {
				if versionStr, ok := versionFound.(string); ok {
					return strings.TrimPrefix(versionStr, "v"), nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not find New Relic provider version in terraform version output")
}

func isValidProviderVersion(currentVersion, minVersion string) bool {
	// Clean version strings by removing "v" prefix if present
	currentVersion = strings.TrimPrefix(currentVersion, "v")
	minVersion = strings.TrimPrefix(minVersion, "v")

	// Parse versions using go-version
	current, err := version.NewVersion(currentVersion)
	if err != nil {
		log.Warnf("Failed to parse current provider version %s: %v", currentVersion, err)
		return false
	}

	minimum, err := version.NewVersion(minVersion)
	if err != nil {
		log.Warnf("Failed to parse minimum provider version %s: %v", minVersion, err)
		return false
	}

	// Check if current version is greater than or equal to minimum required version
	return current.GreaterThanOrEqual(minimum)
}
