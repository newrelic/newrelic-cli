package types

import (
	"os"
	"strings"
)

const (
	ApmKeyword                 = "Apm"
	BuiltinTags                = "nr_deployed_by:newrelic-cli"
	EnvInstallCustomAttributes = "INSTALL_CUSTOM_ATTRIBUTES"
)

// nolint: maligned
type InstallerContext struct {
	AssumeYes   bool
	RecipeNames []string
	RecipePaths []string
	// LocalRecipes is the path to a local recipe directory from which to load recipes.
	LocalRecipes string
	Tags         []string
}

func (i *InstallerContext) RecipePathsProvided() bool {
	return len(i.RecipePaths) > 0
}

func (i *InstallerContext) RecipeNamesProvided() bool {
	return len(i.RecipeNames) > 0

}
func (i *InstallerContext) SetTags(tags []string) {
	i.Tags = tags
	i.Tags = append([]string{BuiltinTags}, i.Tags...)
	os.Setenv(EnvInstallCustomAttributes, strings.Join(i.Tags, ","))
}
