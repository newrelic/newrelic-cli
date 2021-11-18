package install

import (
	"context"
	"net/url"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ConfigValidator interface {
	Validate(ctx context.Context) error
}

// Discoverer is responsible for discovering information about the host system.
type Discoverer interface {
	Discover(context.Context) (*types.DiscoveryManifest, error)
}

// FileFilterer determines the existence of files on the underlying filesystem.
type FileFilterer interface {
	Filter(context.Context, []types.OpenInstallationRecipe) ([]types.OpenInstallationLogMatch, error)
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

type RecipeVarPreparer interface {
	Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error)
}

type RecipeRepository interface {
	FindAll(m types.DiscoveryManifest) []types.OpenInstallationRecipe
}
