package install

// nolint: maligned
type InstallerContext struct {
	AssumeYes   bool
	RecipeNames []string
	RecipePaths []string
	// LocalRecipes is the path to a local recipe directory from which to load recipes.
	LocalRecipes       string
	SkipDiscovery      bool
	SkipIntegrations   bool
	SkipLoggingInstall bool
	SkipApm            bool
}

func (i *InstallerContext) ShouldRunDiscovery() bool {
	return !i.SkipDiscovery
}

func (i *InstallerContext) ShouldInstallInfraAgent() bool {
	return !i.RecipesProvided()
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
