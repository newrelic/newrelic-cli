package types

const (
	ApmKeyword = "Apm"
)

// nolint: maligned
type InstallerContext struct {
	AssumeYes   bool
	RecipeNames []string
	RecipePaths []string
	// LocalRecipes is the path to a local recipe directory from which to load recipes.
	LocalRecipes     string
	SkipInfraInstall bool
}

func (i *InstallerContext) ShouldInstallInfraAgent() bool {
	return !i.RecipesProvided() && !i.SkipInfraInstall
}

func (i *InstallerContext) RecipePathsProvided() bool {
	return len(i.RecipePaths) > 0
}

func (i *InstallerContext) RecipeNamesProvided() bool {
	return len(i.RecipeNames) > 0
}

func (i *InstallerContext) RecipesProvided() bool {
	return i.RecipePathsProvided() || i.RecipeNamesProvided()
}

func (i *InstallerContext) SkipNames() []string {
	skipNames := []string{}
	if i.RecipesProvided() {
		// Skip infra only in Targeted
		if i.SkipInfraInstall {
			skipNames = append(skipNames, InfraAgentRecipeName)
		}
	}

	return skipNames
}
