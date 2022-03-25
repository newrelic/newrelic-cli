package install

import (
	"context"
	"net/url"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ConfigValidator interface {
	Validate(ctx context.Context) error
}

// Discoverer is responsible for discovering information about the host system.
type Discoverer interface {
	Discover(context.Context) (*types.DiscoveryManifest, error)
}

type Prompter interface {
	PromptYesNo(msg string) (bool, error)
	MultiSelect(msg string, options []string) ([]string, error)
}

type RecipeFileFetcher interface {
	FetchRecipeFile(recipeURL *url.URL) (*types.OpenInstallationRecipe, error)
	LoadRecipeFile(filename string) (*types.OpenInstallationRecipe, error)
}

type RecipeFilterRunner interface {
	RunFilterAll(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) []types.OpenInstallationRecipe
	EnsureDoesNotFilter(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) error
}

// RecipeValidator validates installation of a recipe.
type RecipeValidator interface {
	ValidateRecipe(context.Context, types.DiscoveryManifest, types.OpenInstallationRecipe, types.RecipeVars) (entityGUID string, err error)
}

type AgentValidator interface {
	Validate(ctx context.Context, url string) (string, error)
}

type RecipeVarPreparer interface {
	Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error)
}

type RecipeRepository interface {
	FindAll(m types.DiscoveryManifest) []types.OpenInstallationRecipe
}

// RecipeInstaller wrapper responsible for performing recipe validation, installation, and reporting install status
type RecipeInstaller interface {
	promptIfNotLatestCLIVersion(ctx context.Context) error
	Install() error
	install(ctx context.Context) error
	assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error
	discover(ctx context.Context) (*types.DiscoveryManifest, error)
	executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars) (string, error)
	validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars) (string, error)
	executeAndValidateWithProgress(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, assumeYes bool) (string, error)
}

type RecipeBundler interface {
	CreateCoreBundle() *recipes.Bundle
	CreateAdditionalTargetedBundle(names []string) *recipes.Bundle
	CreateAdditionalGuidedBundle() *recipes.Bundle
}
type RecipeBundleInstaller interface {
	InstallStopOnError(bundle *recipes.Bundle, assumeYes bool) error
	InstallContinueOnError(bundle *recipes.Bundle, assumeYes bool)
	InstalledRecipesCount() int
}
