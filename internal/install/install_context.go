package install

type installContext struct {
	specifyActions    bool
	interactiveMode   bool
	installLogging    bool
	installInfraAgent bool
	recipeNames       []string
	recipeFilenames   []string
}

func (i *installContext) ShouldInstallInfraAgent() bool {
	return !i.RecipeFilenamesProvided() && (!i.specifyActions || i.installInfraAgent)
}

func (i *installContext) ShouldInstallLogging() bool {
	return !i.RecipeFilenamesProvided() && (!i.specifyActions || i.installLogging)
}

func (i *installContext) RecipeFilenamesProvided() bool {
	return len(i.recipeFilenames) > 0
}

func (i *installContext) RecipeNamesProvided() bool {
	return len(i.recipeNames) > 0
}
