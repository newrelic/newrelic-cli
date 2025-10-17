package migrate

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// generateProviderConfigCI creates the provider configuration file
func generateProviderConfigCI(config *ToolConfig, workspacePath string) error {
	providerPath := filepath.Join(workspacePath, ProviderConfigFile)
	if err := ioutil.WriteFile(providerPath, []byte(ProviderConfigTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", ProviderConfigFile, err)
	}

	log.Infof("Generated provider configuration for %s: %s", config.DisplayName, providerPath)
	return nil
}

// generateImportBlocksCI creates the import configuration file
func generateImportBlocksCI(workspacePath string, dropRuleData *DropRuleInput) error {
	var importBlocks []string

	for _, resource := range dropRuleData.DropRuleResourceIDs {
		importBlock := fmt.Sprintf(ImportBlockTemplate, resource.Name, resource.PipelineCloudRuleEntityID)
		importBlocks = append(importBlocks, importBlock)
	}

	importConfig := strings.Join(importBlocks, "\n\n")

	importsPath := filepath.Join(workspacePath, ImportConfigFile)
	if err := ioutil.WriteFile(importsPath, []byte(importConfig), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", ImportConfigFile, err)
	}

	log.Infof("Generated import configuration: %s", importsPath)
	return nil
}

// generateRemovedBlocksCI creates the removed blocks configuration file
func generateRemovedBlocksCI(workspacePath string, dropRuleData *DropRuleInput) error {
	var removedBlocks []string

	for _, resource := range dropRuleData.DropRuleResourceIDs {
		removedBlock := fmt.Sprintf(RemovedBlockTemplate, resource.Name)
		removedBlocks = append(removedBlocks, removedBlock)
	}

	removedConfig := strings.Join(removedBlocks, "\n\n")

	removalsPath := filepath.Join(workspacePath, RemovedBlocksFile)
	if err := ioutil.WriteFile(removalsPath, []byte(removedConfig), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", RemovedBlocksFile, err)
	}

	log.Infof("Generated removed blocks configuration: %s", removalsPath)
	return nil
}

// generateAllConfigFilesCI generates all required configuration files
func generateAllConfigFilesCI(config *ToolConfig, workspacePath string, dropRuleData *DropRuleInput) error {
	// Generate provider configuration
	if err := generateProviderConfigCI(config, workspacePath); err != nil {
		return fmt.Errorf("failed to generate provider config: %v", err)
	}

	// Generate import blocks
	if err := generateImportBlocksCI(workspacePath, dropRuleData); err != nil {
		return fmt.Errorf("failed to generate import blocks: %v", err)
	}

	return nil
}
