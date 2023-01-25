package types

import (
	"os"
	"strings"
)

const (
	ApmKeyword                 = "Apm"
	DeployedByTagKey           = "nr_deployed_by"
	TagSeparator               = ":"
	BuiltinTags                = DeployedByTagKey + TagSeparator + "newrelic-cli"
	EnvInstallCustomAttributes = "INSTALL_CUSTOM_ATTRIBUTES"
)

// nolint: maligned
type InstallerContext struct {
	AssumeYes   bool
	RecipeNames []string
	RecipePaths []string
	// LocalRecipes is the path to a local recipe directory from which to load recipes.
	LocalRecipes string
	// Tags         []string
	deployedBy string
}

func (i *InstallerContext) RecipePathsProvided() bool {
	return len(i.RecipePaths) > 0
}

func (i *InstallerContext) RecipeNamesProvided() bool {
	return len(i.RecipeNames) > 0

}

func (i *InstallerContext) SetTags(tags []string) {
	csv := ""
	for _, value := range tags {
		parts := strings.Split(value, TagSeparator)
		if len(parts) == 2 {
			if parts[0] == DeployedByTagKey {
				i.deployedBy = parts[1]
			}
			if len(csv) > 0 {
				csv += ","
			}
			csv += value
		}
	}
	if !strings.Contains(csv, DeployedByTagKey) {
		if len(csv) > 0 {
			csv += ","
		}
		csv += BuiltinTags
	}
	os.Setenv(EnvInstallCustomAttributes, csv)
}

func (i *InstallerContext) GetDeployedBy() string {
	return i.deployedBy
}
