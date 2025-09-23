package migrate

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	destroyWorkspacePath        string
	destroyResourceIdentifiers  []string
	destroySkipResponseToPrompt bool
)

var cmdNRQLDropRulesTFDestroy = &cobra.Command{
	Use:   "tf-destroy",
	Short: "Destroy NRQL drop rules in Terraform workspace",
	Long: `
			Destroy NRQL drop rules in your Terraform workspace by generating and executing
			terraform destroy commands. This command will attempt to run terraform commands
			in your workspace to destroy the drop rule resources.
	`,
	Example: `  # Destroy drop rules in current directory
  newrelic migrate nrqldroprules tf-destroy

  # Destroy drop rules in specific workspace
  newrelic migrate nrqldroprules tf-destroy --workspacePath /path/to/terraform

  # Destroy specific resources without prompts
  newrelic migrate nrqldroprules tf-destroy --resourceIdentifiers resource1,resource2 --skipResponseToPrompt`,
	Run: func(cmd *cobra.Command, args []string) {
		runNRQLDropRulesTFDestroy()
	},
}

func init() {
	cmdNRQLDropRules.AddCommand(cmdNRQLDropRulesTFDestroy)

	cmdNRQLDropRulesTFDestroy.Flags().StringVar(&destroyWorkspacePath, "workspacePath", ".", "path to the Terraform workspace")
	cmdNRQLDropRulesTFDestroy.Flags().StringSliceVar(&destroyResourceIdentifiers, "resourceIdentifiers", []string{}, "list of resource identifiers for newrelic_nrql_drop_rule resources")
	cmdNRQLDropRulesTFDestroy.Flags().BoolVar(&destroySkipResponseToPrompt, "skipResponseToPrompt", false, "skip all user prompts (answers 'N' to all prompts)")
}

func runNRQLDropRulesTFDestroy() {
	// Resolve workspace path
	absWorkspacePath, err := filepath.Abs(destroyWorkspacePath)
	if err != nil {
		log.Fatalf("Error resolving workspace path: %v", err)
	}

	log.Infof("Using Terraform workspace: %s", absWorkspacePath)

	// Try to run terraform state list
	dropRuleResources, err := getTerraformDropRuleResources(absWorkspacePath)
	if err != nil {
		log.Warnf("Failed to list Terraform state: %v", err)
		handleDestroyTerraformStateFailure()
		return
	}

	if len(dropRuleResources) > 0 {
		handleDestroyTerraformStateSuccess(absWorkspacePath, dropRuleResources)
	} else {
		handleDestroyTerraformStateFailure()
	}
}

func handleDestroyTerraformStateSuccess(workspacePath string, resources []string) {
	log.Infof("Found %d NRQL drop rule resources in Terraform state", len(resources))

	// Generate target flags
	targetFlags := make([]string, len(resources))
	for i, resource := range resources {
		targetFlags[i] = fmt.Sprintf("-target=%s", resource)
	}
	targetString := strings.Join(targetFlags, " ")

	// Generate commands
	planCommand := fmt.Sprintf("terraform plan -destroy %s", targetString)
	destroyCommand := fmt.Sprintf("terraform destroy %s", targetString)

	// Print commands
	fmt.Printf("\nGenerated Terraform commands:\n")
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", destroyCommand)

	// Ask user if they want to execute
	if destroySkipResponseToPrompt {
		fmt.Println("\nSkipping execution due to --skipResponseToPrompt flag")
		return
	}

	if !promptForDestroyExecution() {
		fmt.Println("\nExecution halted. Please run the commands above manually in your Terraform workspace.")
		return
	}

	// Execute commands
	executeDestroyTerraformCommands(workspacePath, planCommand, destroyCommand, resources)
}

func handleDestroyTerraformStateFailure() {
	if len(destroyResourceIdentifiers) == 0 {
		log.Fatal("Unable to list Terraform state and no --resourceIdentifiers provided. Please specify resource identifiers for newrelic_nrql_drop_rule resources.")
	}

	log.Infof("Using provided resource identifiers: %v", destroyResourceIdentifiers)

	// Generate target flags from provided identifiers
	targetFlags := make([]string, len(destroyResourceIdentifiers))
	for i, resource := range destroyResourceIdentifiers {
		targetFlags[i] = fmt.Sprintf("-target=%s", resource)
	}
	targetString := strings.Join(targetFlags, " ")

	// Generate and print commands
	planCommand := fmt.Sprintf("terraform plan -destroy %s", targetString)
	destroyCommand := fmt.Sprintf("terraform destroy %s", targetString)

	fmt.Printf("\nGenerated Terraform commands for provided resources:\n")
	fmt.Printf("1. %s\n", planCommand)
	fmt.Printf("2. %s\n", destroyCommand)
	fmt.Println("\nPlease run these commands in your appropriate Terraform workspace.")
}

func promptForDestroyExecution() bool {
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

func executeDestroyTerraformCommands(workspacePath, planCommand, destroyCommand string, resources []string) {
	// Execute plan command
	fmt.Println("\nExecuting terraform plan...")
	if err := executeDestroyTerraformCommand(workspacePath, planCommand); err != nil {
		log.Errorf("Terraform plan failed: %v", err)
		return
	}

	// Ask for confirmation before destroy
	fmt.Print("\nProceed with terraform destroy? (Y/N): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Errorf("Error reading input: %v", err)
		return
	}

	response = strings.TrimSpace(strings.ToUpper(response))
	if response != "Y" && response != "YES" {
		fmt.Println("Destroy cancelled.")
		return
	}

	// Execute destroy command
	fmt.Println("\nExecuting terraform destroy...")
	if err := executeDestroyTerraformCommand(workspacePath, destroyCommand); err != nil {
		log.Errorf("Terraform destroy failed: %v", err)
		return
	}

	// Validate destruction
	validateDropRuleDestruction(workspacePath, resources)
}

func executeDestroyTerraformCommand(workspacePath, command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = workspacePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func validateDropRuleDestruction(workspacePath string, resources []string) {
	fmt.Println("\nValidating drop rule destruction...")

	destroyedCount := 0
	for _, resource := range resources {
		if isDestroyed := checkResourceDestroyed(workspacePath, resource); isDestroyed {
			destroyedCount++
			fmt.Printf("✓ %s: Successfully destroyed\n", resource)
		} else {
			fmt.Printf("⚠ %s: Still exists in state\n", resource)
		}
	}

	if destroyedCount == len(resources) {
		fmt.Printf("\n✅ All %d NRQL drop rule resources have been successfully destroyed\n", destroyedCount)
	} else {
		fmt.Printf("\n⚠️ %d out of %d resources were destroyed. Please check the remaining resources manually.\n", destroyedCount, len(resources))
	}
}

func checkResourceDestroyed(workspacePath, resource string) bool {
	cmd := exec.Command("terraform", "state", "show", resource)
	cmd.Dir = workspacePath

	_, err := cmd.Output()
	if err != nil {
		// If the command fails, it likely means the resource no longer exists in state
		return true
	}

	// If the command succeeds, the resource still exists
	return false
}
