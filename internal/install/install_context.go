package install

type installContext struct {
	skipDiscovery      bool
	skipLoggingInstall bool
	skipInfraInstall   bool
	skipIntegrations   bool
	recipeNames        []string
	recipePaths        []string
}

func (i *installContext) ShouldRunDiscovery() bool {
	return !i.skipDiscovery
}

func (i *installContext) ShouldInstallInfraAgent() bool {
	return !i.RecipesProvided() && !i.skipInfraInstall
}

func (i *installContext) ShouldInstallLogging() bool {
	return !i.RecipesProvided() && !i.skipLoggingInstall
}

func (i *installContext) ShouldInstallIntegrations() bool {
	return i.RecipesProvided() || !i.skipIntegrations
}

func (i *installContext) RecipePathsProvided() bool {
	return len(i.recipePaths) > 0
}

func (i *installContext) RecipeNamesProvided() bool {
	return len(i.recipeNames) > 0
}

func (i *installContext) RecipesProvided() bool {
	return i.RecipePathsProvided() || i.RecipeNamesProvided()
}
