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
	LocalRecipes       string
	SkipIntegrations   bool
	SkipLoggingInstall bool
	SkipApm            bool
	SkipInfraInstall   bool
}

func (i *InstallerContext) ShouldInstallInfraAgent() bool {
	return !i.RecipesProvided() && !i.SkipInfraInstall
}

func (i *InstallerContext) ShouldInstallLogging() bool {
	return !i.RecipesProvided() && !i.SkipLoggingInstall
}

func (i *InstallerContext) ShouldInstallIntegrations() bool {
	return i.RecipesProvided() || !i.SkipIntegrations
}

func (i *InstallerContext) ShouldInstallApm() bool {
	return i.RecipesProvided() || !i.SkipApm
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

	if i.SkipLoggingInstall {
		skipNames = append(skipNames, LoggingRecipeName)
	}

	return skipNames
}

func (i *InstallerContext) SkipTypes() []string {
	skipTypes := []string{}
	if i.SkipIntegrations {
		skipTypes = append(skipTypes, string(OpenInstallationTargetTypeTypes.HOST))
	}

	return skipTypes
}

func (i *InstallerContext) SkipKeywords() []string {
	skipKeywords := []string{}
	if i.SkipApm {
		skipKeywords = append(skipKeywords, ApmKeyword)
	}

	return skipKeywords
}
