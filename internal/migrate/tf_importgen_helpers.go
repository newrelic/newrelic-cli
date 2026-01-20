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
)

// CommandType represents the type of operation being performed
type CommandType int

const (
	CommandUpdate CommandType = iota
	CommandImport
	CommandDelist // Removed CommandDestroy since we use CommandDelist for safe state removal
)

// CommandContext holds shared context and configuration for commands
type CommandContext struct {
	ToolConfig    *ToolConfig
	WorkspacePath string
	SkipPrompts   bool
	CommandType   CommandType
	ResourceIDs   []string
}

// NewCommandContext creates and initializes a command context
func NewCommandContext(useTofu bool, workspacePath string, skipPrompts bool, cmdType CommandType, resourceIDs []string) (*CommandContext, error) {
	toolConfig := newToolConfigCI(useTofu)
	absPath, err := resolveWorkspacePathCI(workspacePath)
	if err != nil {
		return nil, err
	}

	return &CommandContext{
		ToolConfig:    toolConfig,
		WorkspacePath: absPath,
		SkipPrompts:   skipPrompts,
		CommandType:   cmdType,
		ResourceIDs:   resourceIDs,
	}, nil
}

// InitializeCommand performs common initialization steps
func (ctx *CommandContext) InitializeCommand() error {
	log.Infof("Using %s workspace: %s", ctx.ToolConfig.DisplayName, ctx.WorkspacePath)
	return checkTerraformPrerequisitesCI(ctx.ToolConfig)
}

// GetDropRuleResources retrieves NRQL drop rule resources from state
func (ctx *CommandContext) GetDropRuleResources() ([]string, error) {
	cmd := exec.Command(ctx.ToolConfig.ToolName, "state", "list")
	cmd.Dir = ctx.WorkspacePath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s state list failed: %v", ctx.ToolConfig.DisplayName, err)
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

// GenerateTargetCommands creates terraform commands with target flags
func (ctx *CommandContext) GenerateTargetCommands(resources []string) (planCmd, actionCmd string) {
	targetFlags := make([]string, len(resources))
	for i, resource := range resources {
		targetFlags[i] = fmt.Sprintf("-target=%s", resource)
	}
	targetString := strings.Join(targetFlags, " ")

	switch ctx.CommandType {
	case CommandUpdate:
		planCmd = fmt.Sprintf("%s plan -refresh-only %s", ctx.ToolConfig.ToolName, targetString)
		actionCmd = fmt.Sprintf("%s apply -refresh-only %s", ctx.ToolConfig.ToolName, targetString)
	case CommandDelist:
		// For delist, we don't use plan/apply but state rm commands
		planCmd = "# No plan needed for state removal"
		actionCmd = fmt.Sprintf("%s state rm %s", ctx.ToolConfig.ToolName, strings.Join(resources, " "))
	}

	return planCmd, actionCmd
}

// PrintCommands displays generated commands to the user
func (ctx *CommandContext) PrintCommands(planCmd, actionCmd string) {
	fmt.Printf("\nGenerated %s commands:\n", ctx.ToolConfig.DisplayName)
	fmt.Printf("1. %s\n", planCmd)
	fmt.Printf("2. %s\n", actionCmd)
}

// PrintCommandsForResourceIDs displays commands for provided resource identifiers
func (ctx *CommandContext) PrintCommandsForResourceIDs() {
	planCmd, actionCmd := ctx.GenerateTargetCommands(ctx.ResourceIDs)
	fmt.Printf("\nGenerated %s commands for provided resources:\n", ctx.ToolConfig.DisplayName)
	fmt.Printf("1. %s\n", planCmd)
	fmt.Printf("2. %s\n", actionCmd)
	fmt.Printf("\nPlease run these commands in your appropriate %s workspace.\n", ctx.ToolConfig.DisplayName)
}

// PromptForExecution asks user if they want to execute commands
func (ctx *CommandContext) PromptForExecution() bool {
	if ctx.SkipPrompts {
		fmt.Println("\nSkipping execution due to --skipResponseToPrompt flag")
		return false
	}

	fmt.Print("\nWould you like this CLI to execute the commands above? (Y/N): ")
	return readUserInput()
}

// PromptForActionConfirmation asks for confirmation before destructive actions
func (ctx *CommandContext) PromptForActionConfirmation() bool {
	actionName := "apply"
	switch ctx.CommandType {
	case CommandDelist:
		actionName = "delist from state"
	}

	fmt.Printf("\nProceed with %s %s? (Y/N): ", ctx.ToolConfig.ToolName, actionName)
	return readUserInput()
}

// ExecuteCommand runs a command with proper logging and immediate output display
func (ctx *CommandContext) ExecuteCommand(command, actionName string) error {
	fmt.Printf("\n🔄 Executing %s %s...\n", ctx.ToolConfig.ToolName, actionName)
	fmt.Printf("Command: %s\n", command)
	fmt.Println(strings.Repeat("-", 50))

	if err := executeCommandWithOutput(ctx.WorkspacePath, command); err != nil {
		fmt.Printf("\n❌ %s %s failed!\n", ctx.ToolConfig.DisplayName, actionName)
		return fmt.Errorf("%s %s failed: %v", ctx.ToolConfig.DisplayName, actionName, err)
	}

	fmt.Printf("\n✅ %s %s completed successfully!\n", ctx.ToolConfig.DisplayName, actionName)
	return nil
}

// ExecuteStandardFlow runs the standard plan->confirm->action workflow with improved UX
func (ctx *CommandContext) ExecuteStandardFlow(planCmd, actionCmd string, resources []string) error {
	// For delist operations, skip plan and go directly to confirmation
	if ctx.CommandType == CommandDelist {
		return executeDelistFlow(ctx, generateStateRmCommands(ctx.ToolConfig.ToolName, resources), resources)
	}

	// Execute plan and show output immediately
	if err := ctx.ExecuteCommand(planCmd, "plan"); err != nil {
		return err
	}

	// Show plan completion and ask for apply confirmation
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("📋 Plan completed successfully! Review the changes above.\n")
	fmt.Println(strings.Repeat("=", 60))

	// Get confirmation for apply
	if !ctx.PromptForActionConfirmation() {
		fmt.Println("\n🛑 Action cancelled by user.")
		return nil
	}

	// Add auto-approve to apply command and show warning
	actionName := "apply"

	autoApproveCmd := actionCmd + " -auto-approve"

	fmt.Printf("\n⚠️  WARNING: Using -auto-approve flag to avoid interactive prompts.\n")
	fmt.Printf("The %s operation will proceed without additional confirmation.\n", actionName)
	fmt.Println(strings.Repeat("-", 60))

	// Execute action with auto-approve
	if err := ctx.ExecuteCommand(autoApproveCmd, actionName); err != nil {
		return err
	}

	// Validate results
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🔍 Validating %s results...\n", actionName)
	fmt.Println(strings.Repeat("=", 60))
	ctx.ValidateResults(resources)
	return nil
}

// ValidateResults performs post-action validation based on command type
func (ctx *CommandContext) ValidateResults(resources []string) {
	switch ctx.CommandType {
	case CommandUpdate:
		validateResourceUpdates(ctx.WorkspacePath, resources, ctx.ToolConfig)
	case CommandDelist:
		validateDelisting(ctx, resources)
	}
}

// CheckProviderVersion validates New Relic provider version for update operations
func (ctx *CommandContext) CheckProviderVersion() error {
	if ctx.CommandType != CommandUpdate {
		return nil
	}

	log.Infof("Checking New Relic %s Provider version...", ctx.ToolConfig.DisplayName)
	providerVersion, err := getNewRelicProviderVersion(ctx.WorkspacePath, ctx.ToolConfig)
	if err != nil {
		log.Warnf("Could not determine New Relic %s Provider version: %v", ctx.ToolConfig.DisplayName, err)
		log.Infof("Skipping provider version check. Note that New Relic %s Provider >= 3.68.0 is required for pipeline_cloud_rule_entity_id support.", ctx.ToolConfig.DisplayName)
		return nil
	}

	log.Infof("Detected New Relic %s Provider version: %s", ctx.ToolConfig.DisplayName, providerVersion)

	if !isValidProviderVersion(providerVersion, "3.68.0") {
		return fmt.Errorf("changes to add pipeline_cloud_rule_entity_id corresponding to drop rules would not be added to the state with New Relic %s Provider version %s. Provider version >= 3.68.0 is required. Please upgrade your provider", ctx.ToolConfig.DisplayName, providerVersion)
	}

	log.Infof("New Relic %s Provider version check passed.", ctx.ToolConfig.DisplayName)
	return nil
}

// Utility functions

func readUserInput() bool {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Errorf("Error reading input: %v", err)
		return false
	}

	response = strings.TrimSpace(strings.ToUpper(response))
	return response == "Y" || response == "YES"
}

func getNewRelicProviderVersion(workspacePath string, config *ToolConfig) (string, error) {
	cmd := exec.Command(config.ToolName, "version", "-json")
	cmd.Dir = workspacePath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get %s version: %v", config.ToolName, err)
	}

	var versionData map[string]interface{}
	if err := json.Unmarshal(output, &versionData); err != nil {
		return "", fmt.Errorf("failed to parse %s version JSON: %v", config.ToolName, err)
	}

	if providerSelections, ok := versionData["provider_selections"].(map[string]interface{}); ok {
		for providerKey, versionFound := range providerSelections {
			if strings.Contains(providerKey, "newrelic/newrelic") {
				if versionStr, ok := versionFound.(string); ok {
					return strings.TrimPrefix(versionStr, "v"), nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not find New Relic provider version in %s version output", config.ToolName)
}

func isValidProviderVersion(currentVersion, minVersion string) bool {
	currentVersion = strings.TrimPrefix(currentVersion, "v")
	minVersion = strings.TrimPrefix(minVersion, "v")

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

	return current.GreaterThanOrEqual(minimum)
}

func validateResourceUpdates(workspacePath string, resources []string, config *ToolConfig) {
	fmt.Println("\nValidating drop rule updates...")

	updatedCount := 0
	for _, resource := range resources {
		if hasEntityID := checkResourceHasEntityID(workspacePath, resource, config); hasEntityID {
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

func checkResourceHasEntityID(workspacePath, resource string, config *ToolConfig) bool {
	cmd := exec.Command(config.ToolName, "state", "show", resource)
	cmd.Dir = workspacePath

	output, err := cmd.Output()
	if err != nil {
		log.Warnf("Failed to show state for %s: %v", resource, err)
		return false
	}

	return strings.Contains(string(output), "pipeline_cloud_rule_entity_id")
}

// Add helper functions for delisting functionality
func generateStateRmCommands(toolName string, resources []string) []string {
	var commands []string
	for _, resource := range resources {
		cmd := fmt.Sprintf("%s state rm %s", toolName, resource)
		commands = append(commands, cmd)
	}
	return commands
}

func executeDelistFlow(ctx *CommandContext, stateRmCommands []string, resources []string) error {
	fmt.Println("\n🔄 Starting state delisting process...")
	fmt.Println("This will remove resources from Terraform state without destroying them.")

	// Execute each state rm command
	for i, command := range stateRmCommands {
		fmt.Printf("\n[%d/%d] Delisting resource from state...\n", i+1, len(stateRmCommands))
		if err := ctx.ExecuteCommand(command, "state rm"); err != nil {
			return fmt.Errorf("failed to delist resource: %v", err)
		}
	}

	// Validate delisting
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("🔍 Validating delisting results...\n")
	fmt.Println(strings.Repeat("=", 60))
	validateDelisting(ctx, resources)

	return nil
}

func validateDelisting(ctx *CommandContext, originalResources []string) {
	fmt.Println("\nValidating resource delisting...")

	// Check remaining resources in state
	remainingResources, err := ctx.GetDropRuleResources()
	if err != nil {
		fmt.Printf("⚠️ Could not verify delisting: %v\n", err)
		return
	}

	delistedCount := len(originalResources) - len(remainingResources)

	if delistedCount == len(originalResources) {
		fmt.Printf("\n✅ All %d NRQL drop rule resources successfully delisted from state\n", delistedCount)
	} else {
		fmt.Printf("\n⚠️ %d out of %d resources were delisted. Some resources may still be in state.\n", delistedCount, len(originalResources))

		if len(remainingResources) > 0 {
			fmt.Println("\n🔍 Remaining drop rule resources in state:")
			for _, resource := range remainingResources {
				fmt.Printf("  - %s\n", resource)
			}
		}
	}
}

func printPostDelistInstructions(toolDisplayName string, resources []string) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("🔔 IMPORTANT POST-DELISTING INSTRUCTIONS\n")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("\n✅ Resources have been delisted from %s state management\n", toolDisplayName)
	fmt.Printf("🛡️  The actual drop rules remain ACTIVE in New Relic\n\n")

	fmt.Printf("📝 NEXT REQUIRED STEP - Comment out configurations:\n")
	fmt.Printf("   To prevent these resources from being recreated, you MUST comment out\n")
	fmt.Printf("   or remove the following resource configurations from your .tf files:\n\n")

	for _, resource := range resources {
		fmt.Printf("   📄 Resource: %s\n", resource)
	}

	fmt.Printf("\n💡 Example of what to do in your .tf files:\n")
	fmt.Printf("   # Comment out the entire resource block like this:\n")
	fmt.Printf("   # resource \"newrelic_nrql_drop_rule\" \"my_rule\" {\n")
	fmt.Printf("   #   account_id  = var.account_id\n")
	fmt.Printf("   #   name        = \"Drop rule name\"\n")
	fmt.Printf("   #   action      = \"drop_data\"\n")
	fmt.Printf("   #   nrql        = \"SELECT * FROM Log\"\n")
	fmt.Printf("   # }\n\n")

	fmt.Printf("⚠️  WARNING: If you don't comment out the configurations:\n")
	fmt.Printf("   - Next %s plan will show these resources as \"to be created\"\n", toolDisplayName)
	fmt.Printf("   - Next %s apply will attempt to recreate them\n", toolDisplayName)
	fmt.Printf("   - This may cause conflicts with existing drop rules\n\n")

	fmt.Printf("🔄 WORKFLOW SUMMARY:\n")
	fmt.Printf("   1. ✅ Resources delisted from %s state (completed)\n", toolDisplayName)
	fmt.Printf("   2. 📝 Comment out resource configurations in .tf files (YOU NEED TO DO THIS)\n")
	fmt.Printf("   3. 🔍 Run '%s plan' to verify no unwanted changes\n", strings.ToLower(toolDisplayName))
	fmt.Printf("   4. 🎯 Drop rules continue working normally in New Relic\n")

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("🎉 Delisting process completed! Don't forget to comment out configurations!\n")
	fmt.Println(strings.Repeat("=", 70))
}

// Import-specific helpers

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

func ValidateAndExtractPipelineRuleIDs(workspacePath string, resources []string, config *ToolConfig) ([]string, error) {
	var pipelineRuleIDs []string

	for _, resource := range resources {
		cmd := exec.Command(config.ToolName, "show", "-json")
		cmd.Dir = workspacePath

		output, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get %s state JSON: %v", config.DisplayName, err)
		}

		var state TerraformState
		if err := json.Unmarshal(output, &state); err != nil {
			return nil, fmt.Errorf("failed to parse %s state JSON: %v", config.DisplayName, err)
		}

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

func GenerateImportConfigFromResources(resources []string) string {
	var importBlocks []string

	for _, resource := range resources {
		resourceParts := strings.Split(resource, ".")
		if len(resourceParts) < 2 {
			log.Warnf("Invalid resource format: %s", resource)
			continue
		}

		resourceIdentifier := strings.Join(resourceParts[1:], ".")

		importBlock := fmt.Sprintf(`import {
  to = newrelic_pipeline_cloud_rule.%s
  id = %s.pipeline_cloud_rule_entity_id
}`, resourceIdentifier, resource)
		importBlocks = append(importBlocks, importBlock)
	}

	return strings.Join(importBlocks, "\n\n")
}

func GenerateImportConfigFromIDs(pipelineRuleIDs []string) string {
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

func WriteConfigToFile(workspacePath, config, fileName string) {
	filePath := filepath.Join(workspacePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		log.Errorf("Failed to create file %s: %v", filePath, err)
		fmt.Printf("\nFailed to write to file. Configuration:\n")
		fmt.Println(config)
		return
	}
	defer file.Close()

	_, err = file.WriteString(config)
	if err != nil {
		log.Errorf("Failed to write to file %s: %v", filePath, err)
		fmt.Printf("\nFailed to write to file. Configuration:\n")
		fmt.Println(config)
		return
	}

	fmt.Printf("\nConfiguration written to: %s\n", filePath)
}

// checkAndCleanupImportConfig checks for and attempts to delete the generated import configuration file
func checkAndCleanupImportConfig(workspacePath string) {
	importConfigPath := filepath.Join(workspacePath, "import_config_pipeline_rules.tf")

	// Check if the file exists
	if _, err := os.Stat(importConfigPath); os.IsNotExist(err) {
		// File doesn't exist, nothing to clean up
		return
	}

	fmt.Println("\n🧹 Checking for generated import configuration file...")
	fmt.Printf("Found: %s\n", importConfigPath)

	// Attempt to delete the file
	if err := os.Remove(importConfigPath); err != nil {
		fmt.Printf("\n⚠️  WARNING: Could not delete import configuration file\n")
		fmt.Printf("File: %s\n", importConfigPath)
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("📝 MANUAL ACTION REQUIRED: After running this command, please manually delete\n")
		fmt.Printf("   the import configuration file to clean up your workspace:\n")
		fmt.Printf("   rm %s\n", importConfigPath)
		fmt.Printf("   This file contains import{} blocks that are no longer needed after delisting.\n\n")
	} else {
		fmt.Printf("✅ Successfully cleaned up import configuration file: %s\n", importConfigPath)
	}
}

// executeCommandWithOutput executes a command and displays output in real-time
func executeCommandWithOutput(workspacePath, command string) error {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = workspacePath

	// Set up pipes for real-time output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command failed: %s", command)
	}

	return nil
}
