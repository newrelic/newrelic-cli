package install

type installContext struct {
	specifyActions    bool
	interactiveMode   bool
	installLogging    bool
	installInfraAgent bool
	recipeNames       []string
	recipePaths       []string
}

func (i *installContext) ShouldInstallInfraAgent() bool {
	return !i.RecipePathsProvided() && (!i.specifyActions || i.installInfraAgent)
}

func (i *installContext) ShouldInstallLogging() bool {
	return !i.RecipePathsProvided() && (!i.specifyActions || i.installLogging)
}

func (i *installContext) RecipePathsProvided() bool {
	return len(i.recipePaths) > 0
}

func (i *installContext) RecipeNamesProvided() bool {
	return len(i.recipeNames) > 0
}
