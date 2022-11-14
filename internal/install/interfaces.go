package install

import (
	"context"
	"net/url"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ConfigValidator interface {
	Validate(ctx context.Context) error
}

type Prompter interface {
	PromptYesNo(msg string) (bool, error)
	MultiSelect(msg string, options []string) ([]string, error)
}

type ProcessEvaluator interface {
	GetOrLoadProcesses(ctx context.Context) []types.GenericProcess
	DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe) execution.RecipeStatusType
}

type RecipeFileFetcher interface {
	FetchRecipeFile(recipeURL *url.URL) (*types.OpenInstallationRecipe, error)
	LoadRecipeFile(filename string) (*types.OpenInstallationRecipe, error)
}

type RecipeFilterRunner interface {
	RunFilterAll(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) []types.OpenInstallationRecipe
	EnsureDoesNotFilter(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) error
}

//
//// RecipeValidator validates installation of a recipe.
//type RecipeValidator interface {
//	ValidateRecipe(context.Context, types.DiscoveryManifest, types.OpenInstallationRecipe, types.RecipeVars) (entityGUID string, err error)
//}
//
//type AgentValidator interface {
//	Validate(ctx context.Context, url string) (string, error)
//}

type RecipeVarPreparer interface {
	Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error)
}

// RecipeInstaller wrapper responsible for performing recipe validation, installation, and reporting install status
// FIXME remove private methods from interface definition
type RecipeInstaller interface {
	promptIfNotLatestCLIVersion(ctx context.Context) error
	Install() error
	install(ctx context.Context) error
	//assertDiscoveryValid(ctx context.MockContext, m *types.DiscoveryManifest) error
	//discover(ctx context.MockContext) (*types.DiscoveryManifest, error)
	executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars, assumeYes bool) (string, error)
	validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars, assumeYes bool) (string, error)
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

type RecipeStatusDetector interface {
	GetDetectedRecipes() (recipes.RecipeDetectionResults, recipes.RecipeDetectionResults, error)
}
