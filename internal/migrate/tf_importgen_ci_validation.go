package migrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Check for and resolve duplicate names
	resolvedData := resolveDuplicateNames(&dropRuleData)

	return resolvedData, nil
}

// resolveDuplicateNames checks for duplicate resource names and resolves them by adding random suffixes
func resolveDuplicateNames(dropRuleData *DropRuleInput) *DropRuleInput {
	nameMap := make(map[string]int)
	duplicates := make(map[string][]int)

	// First pass: identify all names and their occurrences
	for i, resource := range dropRuleData.DropRuleResourceIDs {
		if count, exists := nameMap[resource.Name]; exists {
			// This is a duplicate
			if count == 1 {
				// First time we see this as duplicate, add the original index
				duplicates[resource.Name] = []int{getFirstOccurrenceIndex(dropRuleData.DropRuleResourceIDs, resource.Name)}
			}
			duplicates[resource.Name] = append(duplicates[resource.Name], i)
			nameMap[resource.Name]++
		} else {
			nameMap[resource.Name] = 1
		}
	}

	// If no duplicates found, return original data
	if len(duplicates) == 0 {
		return dropRuleData
	}

	// Print warning about duplicates
	printDuplicateNamesWarning(duplicates)

	// Second pass: resolve duplicates by adding random suffixes
	resolvedNames := make(map[string]string) // original -> new name mapping

	for originalName, indices := range duplicates {
		log.Warnf("Resolving %d duplicate occurrences of resource name: %s", len(indices), originalName)

		for i, index := range indices {
			if i == 0 {
				// Keep the first occurrence as-is
				log.Infof("  - Keeping first occurrence: %s", originalName)
				continue
			}

			// Generate new name with random suffix for subsequent occurrences
			newName := generateUniqueResourceName(originalName, nameMap, resolvedNames)
			resolvedNames[fmt.Sprintf("%s_%d", originalName, index)] = newName
			dropRuleData.DropRuleResourceIDs[index].Name = newName
			nameMap[newName] = 1 // Mark new name as used

			log.Infof("  - Renamed occurrence %d: %s -> %s", i+1, originalName, newName)
		}
	}

	log.Infof("✅ Successfully resolved all duplicate resource names")
	return dropRuleData
}

// getFirstOccurrenceIndex finds the index of the first occurrence of a name
func getFirstOccurrenceIndex(resources []DropRuleResource, name string) int {
	for i, resource := range resources {
		if resource.Name == name {
			return i
		}
	}
	return -1
}

// generateUniqueResourceName creates a unique resource name by adding a random suffix
func generateUniqueResourceName(baseName string, existingNames map[string]int, resolvedNames map[string]string) string {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	for {
		suffix := generateRandomAlphabeticString(5)
		newName := fmt.Sprintf("%s_%s", baseName, suffix)

		// Check if this name is already used in original names or resolved names
		if _, exists := existingNames[newName]; !exists {
			alreadyResolved := false
			for _, resolvedName := range resolvedNames {
				if resolvedName == newName {
					alreadyResolved = true
					break
				}
			}
			if !alreadyResolved {
				return newName
			}
		}
	}
}

// generateRandomAlphabeticString generates a random string of specified length using only alphabetic characters
func generateRandomAlphabeticString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
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
