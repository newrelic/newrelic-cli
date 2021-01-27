package install

// nolint: maligned
type InstallerContext struct {
	AdvancedMode       bool
	AssumeYes          bool
	RecipeNames        []string
	RecipePaths        []string
	SkipDiscovery      bool
	SkipInfraInstall   bool
	SkipIntegrations   bool
	SkipLoggingInstall bool
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

// ShouldPrompt determines if the user should be prompted for input.
func (i *InstallerContext) ShouldPrompt() bool {
	if i.AdvancedMode {
		return true
	}

	return i.RecipesProvided() || i.AssumeYes
}
