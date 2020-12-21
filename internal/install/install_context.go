package install

type InstallerContext struct {
	SkipDiscovery      bool
	SkipLoggingInstall bool
	SkipInfraInstall   bool
	SkipIntegrations   bool
	RecipeNames        []string
	RecipePaths        []string
}

func (i *InstallerContext) ShouldRunDiscovery() bool {
	return !i.SkipDiscovery
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

func (i *InstallerContext) RecipePathsProvided() bool {
	return len(i.RecipePaths) > 0
}

func (i *InstallerContext) RecipeNamesProvided() bool {
	return len(i.RecipeNames) > 0
}

func (i *InstallerContext) RecipesProvided() bool {
	return i.RecipePathsProvided() || i.RecipeNamesProvided()
}
