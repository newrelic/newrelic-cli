package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// validateInputParametersCI validates that required input parameters are provided correctly
func validateInputParametersCI(inputFile, inputJSON string) error {
	if inputFile == "" && inputJSON == "" {
		return fmt.Errorf("either --file or --json must be provided")
	}

	if inputFile != "" && inputJSON != "" {
		return fmt.Errorf("cannot specify both --file and --json, please provide only one")
	}

	return nil
}

// parseDropRuleInputCI parses the input data from file or JSON string
func parseDropRuleInputCI(inputFile, inputJSON string) (*DropRuleInput, error) {
	var inputData []byte
	var err error

	if inputFile != "" {
		inputData, err = ioutil.ReadFile(inputFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", inputFile, err)
		}
	} else {
		inputData = []byte(inputJSON)
	}

	var dropRuleData DropRuleInput
	if err := json.Unmarshal(inputData, &dropRuleData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &dropRuleData, nil
}

// checkNewRelicEnvironmentVariablesCI validates that required New Relic environment variables are set
func checkNewRelicEnvironmentVariablesCI() error {
	missingVars := []string{}

	for _, envVar := range RequiredEnvVars {
		if value := os.Getenv(envVar); value == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables: %v. Please set these variables and try again", missingVars)
	}

	// Check optional region variable
	if region := os.Getenv(OptionalRegionEnvVar); region == "" {
		log.Info("NEW_RELIC_REGION not set, will default to 'US' region")
	}

	return nil
}

// validateAccountIDConsistencyCI checks if account IDs in the input match the environment variable
func validateAccountIDConsistencyCI(dropRuleData *DropRuleInput) {
	envAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if envAccountID == "" {
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
		printAccountMismatchWarningCI(envAccountID, mismatchedResources)
	}
}

// validateWorkspaceCI checks if the workspace directory is suitable for initialization
func validateWorkspaceCI(workspacePath string) error {
	// Check if directory exists, create if it doesn't
	if _, err := os.Stat(workspacePath); os.IsNotExist(err) {
		if err := os.MkdirAll(workspacePath, 0755); err != nil {
			return fmt.Errorf("failed to create workspace directory: %v", err)
		}
		log.Infof("Created workspace directory: %s", workspacePath)
		return nil
	}

	// Check for conflicting files
	entries, err := ioutil.ReadDir(workspacePath)
	if err != nil {
		return fmt.Errorf("failed to read workspace directory: %v", err)
	}

	conflictingFiles := []string{}
	for _, entry := range entries {
		name := entry.Name()
		for _, pattern := range ConflictingFilePatterns {
			if strings.HasSuffix(name, pattern) || name == pattern {
				conflictingFiles = append(conflictingFiles, name)
				break
			}
		}
	}

	if len(conflictingFiles) > 0 {
		return fmt.Errorf("workspace directory contains conflicting files that may interfere with the import process: %v. Please use an empty directory or clean up existing Terraform files", conflictingFiles)
	}

	return nil
}

// resolveWorkspacePathCI resolves the workspace path to an absolute path
func resolveWorkspacePathCI(workspacePath string) (string, error) {
	absPath, err := filepath.Abs(workspacePath)
	if err != nil {
		return "", fmt.Errorf("error resolving workspace path: %v", err)
	}
	return absPath, nil
}
