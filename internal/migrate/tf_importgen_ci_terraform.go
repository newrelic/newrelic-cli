package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// newToolConfigCI creates a new ToolConfig based on the useTofu flag
func newToolConfigCI(useTofu bool) *ToolConfig {
	if useTofu {
		return &ToolConfig{
			UseTofu:     true,
			ToolName:    TofuToolName,
			DisplayName: OpenTofuDisplay,
		}
	}
	return &ToolConfig{
		UseTofu:     false,
		ToolName:    TerraformToolName,
		DisplayName: TerraformDisplay,
	}
}

// checkTerraformPrerequisitesCI validates that the required tool is installed and meets version requirements
func checkTerraformPrerequisitesCI(config *ToolConfig) error {
	// Check if tool exists
	if !isToolInstalledCI(config) {
		return fmt.Errorf("%s is not installed or not found in PATH. Please install %s and try again", config.DisplayName, config.DisplayName)
	}

	// Check tool version
	log.Infof("Checking %s version...", config.DisplayName)
	version, err := getToolVersionCI(config)
	if err != nil {
		return fmt.Errorf("could not determine %s version: %v", config.DisplayName, err)
	}

	log.Infof("Detected %s version: %s", config.DisplayName, version)

	if !isValidVersionCI(version, MinRequiredVersion) {
		return fmt.Errorf("this command requires %s version >= %s to generate import configuration. Your version: %s. Please update %s and try again",
			config.DisplayName, MinRequiredVersion, version, config.DisplayName)
	}

	log.Infof("%s version check passed.", config.DisplayName)
	return nil
}

// isToolInstalledCI checks if the specified tool is installed and available
func isToolInstalledCI(config *ToolConfig) bool {
	cmd := exec.Command(config.ToolName, "version")
	return cmd.Run() == nil
}

// getToolVersionCI retrieves the version of the specified tool
func getToolVersionCI(config *ToolConfig) (string, error) {
	cmd := exec.Command(config.ToolName, "version", "-json")
	output, _ := cmd.Output()

	// Parse JSON output
	var versionData map[string]interface{}
	if err := json.Unmarshal(output, &versionData); err != nil {
		return "", fmt.Errorf("failed to parse %s version JSON: %v", config.DisplayName, err)
	}

	// Both Terraform and OpenTofu use "terraform_version" field for compatibility
	if version, ok := versionData["terraform_version"].(string); ok {
		return version, nil
	}

	return "", fmt.Errorf("could not find %s version in output", config.DisplayName)
}

// isValidVersionCI checks if the current version meets the minimum requirement
func isValidVersionCI(currentVersion, minVersion string) bool {
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

		if currentVal > minVal {
			return true
		} else if currentVal < minVal {
			return false
		}
		// If equal, continue to next component
	}

	// If we've compared all components and they're equal, versions are equal
	return true
}

// initializeWorkspaceCI runs the init command in the specified workspace
func initializeWorkspaceCI(config *ToolConfig, workspacePath string) error {
	log.Infof("Initializing %s...", config.DisplayName)
	initCommand := fmt.Sprintf("%s init", config.ToolName)
	if err := executeCommandCI(workspacePath, initCommand); err != nil {
		return fmt.Errorf("%s init failed: %v", config.DisplayName, err)
	}
	return nil
}

// generateConfigurationPlanCI runs the plan command with config generation
func generateConfigurationPlanCI(config *ToolConfig, workspacePath string) error {
	log.Infof("Running `%s plan` to generate configuration...", config.ToolName)
	planCommand := fmt.Sprintf("%s plan -generate-config-out=%s", config.ToolName, PipelineCloudRulesFile)
	if err := executeCommandCI(workspacePath, planCommand); err != nil {
		return fmt.Errorf("%s plan failed: %v", config.DisplayName, err)
	}

	// Verify the file was generated
	pcrPath := filepath.Join(workspacePath, PipelineCloudRulesFile)
	if _, err := os.Stat(pcrPath); err != nil {
		return fmt.Errorf("expected %s was not generated: %v", PipelineCloudRulesFile, err)
	}

	log.Infof("Generated Pipeline Cloud Rules configuration: %s", pcrPath)
	return nil
}

// formatConfigurationCI runs the fmt command to format the configuration files
func formatConfigurationCI(config *ToolConfig, workspacePath string) error {
	log.Infof("Formatting %s configuration to finalize structure...", config.DisplayName)
	formatCommand := fmt.Sprintf("%s fmt", config.ToolName)
	if err := executeCommandCI(workspacePath, formatCommand); err != nil {
		return fmt.Errorf("%s fmt failed: %v", config.DisplayName, err)
	}

	log.Infof("%s configuration formatting completed successfully", config.DisplayName)
	return nil
}

// executeCommandCI executes a command in the specified workspace directory
func executeCommandCI(workspacePath, command string) error {
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
