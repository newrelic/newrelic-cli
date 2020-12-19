package install

type InstallContext struct {
	SkipDiscovery      bool
	SkipLoggingInstall bool
	SkipInfraInstall   bool
	SkipIntegrations   bool
	RecipeNames        []string
	RecipePaths        []string
}

func (i *InstallContext) ShouldRunDiscovery() bool {
	return !i.SkipDiscovery
}

func (i *InstallContext) ShouldInstallInfraAgent() bool {
	return !i.RecipesProvided() && !i.SkipInfraInstall
}

func (i *InstallContext) ShouldInstallLogging() bool {
	return !i.RecipesProvided() && !i.SkipLoggingInstall
}

func (i *InstallContext) ShouldInstallIntegrations() bool {
	return i.RecipesProvided() || !i.SkipIntegrations
}

func (i *InstallContext) RecipePathsProvided() bool {
	return len(i.RecipePaths) > 0
}

func (i *InstallContext) RecipeNamesProvided() bool {
	return len(i.RecipeNames) > 0
}

func (i *InstallContext) RecipesProvided() bool {
	return i.RecipePathsProvided() || i.RecipeNamesProvided()
}
