package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	ciInputFile     string
	ciInputJSON     string
	ciWorkspacePath string
	useTofu         bool
)

type DropRuleResource struct {
	Name                      string `json:"name"`
	ID                        string `json:"id"`
	PipelineCloudRuleEntityID string `json:"pipeline_cloud_rule_entity_id"`
}

type DropRuleInput struct {
	DropRuleResourceIDs []DropRuleResource `json:"drop_rule_resource_ids"`
}

var cmdNRQLDropRulesTFImportGenCI = &cobra.Command{
	Use:   "tf-importgen-ci",
	Short: "Generate Terraform import configuration for CI environments",
	Long: `
            Generate Terraform import configuration for CI environments based on drop rules data.
            This command creates a complete Terraform workspace with import blocks, provider 
            configuration, and removed blocks for Pipeline Cloud Rules migration in CI/CD pipelines.
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
}

func runNRQLDropRulesTFImportGenCI() {
	// Validate input parameters
	if ciInputFile == "" && ciInputJSON == "" {
		log.Fatal("Either --file or --json must be provided")
	}

	if ciInputFile != "" && ciInputJSON != "" {
		log.Fatal("Cannot specify both --file and --json, please provide only one")
	}

	// Parse input data
	dropRuleData, err := parseDropRuleInput()
	if err != nil {
		log.Fatalf("Failed to parse input data: %v", err)
	}

	if len(dropRuleData.DropRuleResourceIDs) == 0 {
		log.Fatal("No drop rule resource IDs found in input data")
	}

	// Check New Relic environment variables
	if err := checkNewRelicEnvironmentVariables(); err != nil {
		log.Fatal(err)
	}

	// Validate account ID consistency
	validateAccountIDConsistency(dropRuleData)

	// Resolve workspace path
	absWorkspacePath, err := filepath.Abs(ciWorkspacePath)
	if err != nil {
		log.Fatalf("Error resolving workspace path: %v", err)
	}

	log.Infof("Using %s workspace: %s", getTerraformToolDisplayName(), absWorkspacePath)

	// Check Terraform prerequisites
	if err := checkTerraformPrerequisites(); err != nil {
		log.Fatal(err)
	}

	// Validate workspace is empty or suitable for initialization
	if err := validateWorkspace(absWorkspacePath); err != nil {
		log.Fatal(err)
	}

	// Generate and write Terraform files
	if err := generateTerraformFiles(absWorkspacePath, dropRuleData); err != nil {
		log.Fatalf("Failed to generate %s files: %v", getTerraformToolDisplayName(), err)
	}

	// Initialize and plan Terraform
	if err := initializeAndPlanTerraform(absWorkspacePath); err != nil {
		log.Fatalf("Failed to initialize and plan %s: %v", getTerraformToolDisplayName(), err)
	}

	// Generate removed blocks
	if err := generateRemovedBlocks(absWorkspacePath, dropRuleData); err != nil {
		log.Fatalf("Failed to generate removed blocks: %v", err)
	}

	// Format configuration files
	if err := formatTerraformConfiguration(absWorkspacePath); err != nil {
		log.Warnf("Failed to format configuration: %v", err)
	}

	// Print final recommendations
	printCIRecommendations(dropRuleData)
}

func parseDropRuleInput() (*DropRuleInput, error) {
	var inputData []byte
	var err error

	if ciInputFile != "" {
		inputData, err = ioutil.ReadFile(ciInputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", ciInputFile, err)
		}
	} else {
		inputData = []byte(ciInputJSON)
	}

	var dropRuleData DropRuleInput
	if err := json.Unmarshal(inputData, &dropRuleData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &dropRuleData, nil
}

func checkTerraformPrerequisites() error {
	// toolName := getTerraformToolName()
	toolDisplayName := getTerraformToolDisplayName()

	// Check if tool exists
	if !isTerraformInstalled() {
		return fmt.Errorf("%s is not installed or not found in PATH. Please install %s and try again", toolDisplayName, toolDisplayName)
	}

	// Check tool version
	log.Infof("Checking %s version...", toolDisplayName)
	terraformVersion, err := getTerraformVersion()
	if err != nil {
		return fmt.Errorf("could not determine %s version: %v", toolDisplayName, err)
	}

	log.Infof("Detected %s version: %s", toolDisplayName, terraformVersion)

	if !isValidTerraformVersion(terraformVersion, "1.5") {
		return fmt.Errorf("this command requires %s version >= 1.5 to generate import configuration. Your version: %s. Please update %s and try again", toolDisplayName, terraformVersion, toolDisplayName)
	}

	log.Infof("%s version check passed.", toolDisplayName)
	return nil
}

func isTerraformInstalled() bool {
	toolName := getTerraformToolName()
	cmd := exec.Command(toolName, "version")
	return cmd.Run() == nil
}

func validateWorkspace(workspacePath string) error {
	// Check if directory exists, create if it doesn't
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		if err := os.MkdirAll(workspacePath, 0755); err != nil {
			return fmt.Errorf("failed to create workspace directory: %v", err)
		}
		log.Infof("Created workspace directory: %s", workspacePath)
		return nil
	}

	// Check if directory has existing Terraform state or configuration that might conflict
	entries, err := ioutil.ReadDir(workspacePath)
	if err != nil {
		return fmt.Errorf("failed to read workspace directory: %v", err)
	}

	conflictingFiles := []string{}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".tfstate") ||
			strings.HasSuffix(name, ".tfstate.backup") ||
			name == ".terraform" ||
			name == "provider.tf" ||
			name == "imports.tf" ||
			name == "removals.tf" ||
			name == "pcrs.tf" {
			conflictingFiles = append(conflictingFiles, name)
		}
	}

	if len(conflictingFiles) > 0 {
		return fmt.Errorf("workspace directory contains conflicting files that may interfere with the import process: %v. Please use an empty directory or clean up existing Terraform files", conflictingFiles)
	}

	return nil
}

func generateTerraformFiles(workspacePath string, dropRuleData *DropRuleInput) error {
	// Generate provider configuration
	if err := generateProviderConfig(workspacePath); err != nil {
		return fmt.Errorf("failed to generate provider config: %v", err)
	}

	// Generate import blocks
	if err := generateImportBlocks(workspacePath, dropRuleData); err != nil {
		return fmt.Errorf("failed to generate import blocks: %v", err)
	}

	return nil
}

func generateProviderConfig(workspacePath string) error {
	// toolName := getTerraformToolName()
	toolDisplayName := getTerraformToolDisplayName()

	providerConfig := `terraform {
  required_providers {
    newrelic = {
      source  = "newrelic/newrelic"
      version = "~> 3.0"
    }
  }
}

provider "newrelic" {
  # Configuration will be taken from environment variables:
  # NEW_RELIC_API_KEY
  # NEW_RELIC_ACCOUNT_ID
  # NEW_RELIC_REGION (optional, defaults to US)
}
`

	providerPath := filepath.Join(workspacePath, "provider.tf")
	if err := ioutil.WriteFile(providerPath, []byte(providerConfig), 0644); err != nil {
		return fmt.Errorf("failed to write provider.tf: %v", err)
	}

	log.Infof("Generated provider configuration for %s: %s", toolDisplayName, providerPath)
	return nil
}

func generateImportBlocks(workspacePath string, dropRuleData *DropRuleInput) error {
	var importBlocks []string

	for _, resource := range dropRuleData.DropRuleResourceIDs {
		importBlock := fmt.Sprintf(`import {
  to = newrelic_pipeline_cloud_rule.%s
  id = "%s"
}`, resource.Name, resource.PipelineCloudRuleEntityID)
		importBlocks = append(importBlocks, importBlock)
	}

	importConfig := strings.Join(importBlocks, "\n\n")

	importsPath := filepath.Join(workspacePath, "imports.tf")
	if err := ioutil.WriteFile(importsPath, []byte(importConfig), 0644); err != nil {
		return fmt.Errorf("failed to write imports.tf: %v", err)
	}

	log.Infof("Generated import configuration: %s", importsPath)
	return nil
}

func initializeAndPlanTerraform(workspacePath string) error {
	toolName := getTerraformToolName()
	toolDisplayName := getTerraformToolDisplayName()

	// Run init
	log.Infof("Initializing %s...", toolDisplayName)
	initCommand := fmt.Sprintf("%s init", toolName)
	if err := executeTerraformCommandVCI(workspacePath, initCommand); err != nil {
		return fmt.Errorf("%s init failed: %v", toolDisplayName, err)
	}

	// Run plan with config generation
	log.Infof("Running `%s plan` to generate configuration...", toolName)
	planCommand := fmt.Sprintf("%s plan -generate-config-out=pcrs.tf", toolName)
	if err := executeTerraformCommandVCI(workspacePath, planCommand); err != nil {
		return fmt.Errorf("%s plan failed: %v", toolDisplayName, err)
	}

	pcrPath := filepath.Join(workspacePath, "pcrs.tf")
	if _, err := os.Stat(pcrPath); err != nil {
		return fmt.Errorf("expected pcrs.tf was not generated: %v", err)
	}

	log.Infof("Generated Pipeline Cloud Rules configuration: %s", pcrPath)
	return nil
}

func generateRemovedBlocks(workspacePath string, dropRuleData *DropRuleInput) error {
	var removedBlocks []string

	for _, resource := range dropRuleData.DropRuleResourceIDs {
		removedBlock := fmt.Sprintf(`removed {
  from = newrelic_nrql_drop_rule.%s

  lifecycle {
    destroy = false
  }
}`, resource.Name)
		removedBlocks = append(removedBlocks, removedBlock)
	}

	removedConfig := strings.Join(removedBlocks, "\n\n")

	removalsPath := filepath.Join(workspacePath, "removals.tf")
	if err := ioutil.WriteFile(removalsPath, []byte(removedConfig), 0644); err != nil {
		return fmt.Errorf("failed to write removals.tf: %v", err)
	}

	log.Infof("Generated removed blocks configuration: %s", removalsPath)
	return nil
}

func formatTerraformConfiguration(workspacePath string) error {
	toolName := getTerraformToolName()
	toolDisplayName := getTerraformToolDisplayName()

	log.Infof("Formatting %s configuration to finalize structure...", toolDisplayName)
	formatCommand := fmt.Sprintf("%s fmt", toolName)
	if err := executeTerraformCommandVCI(workspacePath, formatCommand); err != nil {
		return fmt.Errorf("%s fmt failed: %v", toolDisplayName, err)
	}

	log.Infof("%s configuration formatting completed successfully", toolDisplayName)
	return nil
}

func printCIRecommendations(dropRuleData *DropRuleInput) {
	toolName := getTerraformToolName()
	toolDisplayName := getTerraformToolDisplayName()

	// Print attention-grabbing separator
	separator := strings.Repeat("=", 80)
	fmt.Printf("\n%s\n", separator)
	fmt.Println("IMPORTANT: CI/CD MIGRATION RECOMMENDATIONS")
	fmt.Printf("%s\n\n", separator)

	// Show loading animation while "preparing recommendations"
	showLoadingAnimation("Preparing migration recommendations", 2*time.Second)

	fmt.Println("✅ Local workspace setup completed successfully!")
	fmt.Printf("\nNext steps for your CI/CD pipeline migration (using %s):\n", toolDisplayName)
	fmt.Println()

	time.Sleep(1 * time.Second)

	// Step 1: Copy files to CI
	showStepHeader(fmt.Sprintf("1. COPY GENERATED FILES TO YOUR CI ENVIRONMENT (%s):", toolDisplayName))
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("   📁 Copy the following files to your CI %s workspace:\n", toolDisplayName)
	time.Sleep(300 * time.Millisecond)
	fmt.Println("      - pcrs.tf        (Pipeline Cloud Rules configuration)")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("      - imports.tf     (Import blocks)")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("      - removals.tf    (Removed blocks for drop rules)")
	fmt.Println()
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   ⚠️  REQUIREMENT: Ensure %s version >= 1.5 in your CI environment\n", toolDisplayName)
	fmt.Println("       for import block support.")
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Step 2: Comment out drop rules and add removals
	showStepHeader("2. PREPARE YOUR EXISTING CI CONFIGURATION:")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   📝 Comment out ALL existing NRQL drop rule resources in your CI configuration")
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   📝 Copy all content from removals.tf into your CI %s configuration\n", toolDisplayName)
	fmt.Println()
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   ⚠️  REQUIREMENT: Ensure %s version >= 1.7 in your CI environment\n", toolDisplayName)
	fmt.Println("       for removed block support.")
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Alternative approach
	showStepHeader("3. ALTERNATIVE: MANUAL STATE REMOVAL (if you prefer not to use removed blocks):")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   If you prefer to manually remove drop rules from state without using")
	time.Sleep(200 * time.Millisecond)
	fmt.Println("   removed blocks, run this command in your CI environment:")
	fmt.Println()

	time.Sleep(500 * time.Millisecond)

	// Generate state rm command
	var resourceNames []string
	for _, resource := range dropRuleData.DropRuleResourceIDs {
		resourceNames = append(resourceNames, fmt.Sprintf("newrelic_nrql_drop_rule.%s", resource.Name))
	}

	stateRmCommand := fmt.Sprintf("   %s state rm %s", toolName, strings.Join(resourceNames, " "))
	fmt.Println(stateRmCommand)
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	// Final steps
	showStepHeader("4. EXECUTE IN YOUR CI ENVIRONMENT:")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   After copying files and preparing configuration:")
	time.Sleep(300 * time.Millisecond)
	fmt.Printf("   📋 %s plan    (review the migration plan)\n", toolName)
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("   📋 %s apply   (execute the migration)\n", toolName)
	fmt.Println()

	time.Sleep(1500 * time.Millisecond)

	showStepHeader("5. VERIFICATION:")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("   ✅ Verify that Pipeline Cloud Rules are created successfully")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   ✅ Verify that old NRQL drop rules are removed from state (not destroyed)")
	time.Sleep(300 * time.Millisecond)
	fmt.Println("   ✅ Test that your log filtering continues to work as expected")
	fmt.Println()

	time.Sleep(1 * time.Second)

	fmt.Printf("%s\n", separator)
	showLoadingAnimation("Finalizing setup", 1*time.Second)
	fmt.Println("Migration workspace prepared successfully! 🎉")
	fmt.Printf("Generated files are ready in: %s\n", ciWorkspacePath)
	fmt.Printf("%s\n\n", separator)
}

func showStepHeader(step string) {
	fmt.Printf("\033[1;36m%s\033[0m\n", step) // Cyan bold for step headers
}

func showLoadingAnimation(message string, duration time.Duration) {
	fmt.Printf("%s", message)

	dots := []string{".", "..", "...", "...."}
	interval := duration / time.Duration(len(dots)*3) // Show each dot pattern 3 times

	for i := 0; i < len(dots)*3; i++ {
		fmt.Printf("\r%s%s   ", message, dots[i%len(dots)])
		time.Sleep(interval)
	}

	fmt.Printf("\r%s... ✅\n\n", message)
}

func checkNewRelicEnvironmentVariables() error {
	requiredVars := []string{"NEW_RELIC_ACCOUNT_ID", "NEW_RELIC_API_KEY"}
	missingVars := []string{}

	for _, envVar := range requiredVars {
		if value := os.Getenv(envVar); value == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v. Please set these variables and try again", missingVars)
	}

	// NEW_RELIC_REGION is optional, but log if it's not set
	if region := os.Getenv("NEW_RELIC_REGION"); region == "" {
		log.Info("NEW_RELIC_REGION not set, will default to 'US' region")
	}

	return nil
}

func validateAccountIDConsistency(dropRuleData *DropRuleInput) {
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if envAccountID == "" {
		// This should not happen since we check for it earlier, but just in case
		return
	}

	mismatchedResources := []string{}

	for _, resource := range dropRuleData.DropRuleResourceIDs {
		// Extract account ID from pipeline_cloud_rule_entity_id (format: "accountId:ruleId")
		parts := strings.Split(resource.ID, ":")
		if len(parts) < 2 {
			log.Warnf("Unexpected format for resource ID: %s", resource.ID)
			continue
		}

		ruleAccountID := parts[0]
		if ruleAccountID != envAccountID {
			mismatchedResources = append(mismatchedResources, fmt.Sprintf("%s (rule account: %s)", resource.Name, ruleAccountID))
		}
	}

	if len(mismatchedResources) > 0 {
		// Print bold warning
		fmt.Printf("\n\033[1;33m⚠️  WARNING: ACCOUNT ID MISMATCH DETECTED\033[0m\n")
		fmt.Printf("\033[1mEnvironment NEW_RELIC_ACCOUNT_ID: %s\033[0m\n", envAccountID)
		fmt.Printf("\033[1mThe following resources have different account IDs:\033[0m\n")
		for _, resource := range mismatchedResources {
			fmt.Printf("  - \033[1m%s\033[0m\n", resource)
		}
		fmt.Printf("\033[1;33mThis may cause import failures. Please verify your account configuration.\033[0m\n\n")
	}
}

func executeTerraformCommandVCI(workspacePath, command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = workspacePath

	// Capture both stdout and stderr for better error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command failed: %s\nOutput: %s", command, string(output))
	}

	// Log successful command output for transparency
	if len(output) > 0 {
		log.Debugf("Command output: %s", string(output))
	}

	return nil
}

func getTerraformToolName() string {
	if useTofu {
		return "tofu"
	}
	return "terraform"
}

func getTerraformToolDisplayName() string {
	if useTofu {
		return "OpenTofu"
	}
	return "Terraform"
}

func getTerraformVersion() (string, error) {
	toolName := getTerraformToolName()
	toolDisplayName := getTerraformToolDisplayName()
	cmd := exec.Command(toolName, "version", "-json")
	output, _ := cmd.Output()

	// Parse JSON output
	var versionData map[string]interface{}
	if err := json.Unmarshal(output, &versionData); err != nil {
		return "", fmt.Errorf("failed to parse %s version JSON: %v", toolDisplayName, err)
	}

	// Try both possible version field names
	versionField := "terraform_version"
	if useTofu {
		versionField = "terraform_version" // OpenTofu uses the same field name for compatibility
	}

	if version, ok := versionData[versionField].(string); ok {
		return version, nil
	}

	return "", fmt.Errorf("could not find %s version in output", toolDisplayName)
}
