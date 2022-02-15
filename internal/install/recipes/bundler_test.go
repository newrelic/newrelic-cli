package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"strings"

	"github.com/stretchr/testify/require"
)

var (
	bundler_discoveryManifest types.DiscoveryManifest
	bundler_recipeCache       []types.OpenInstallationRecipe
	bundler_repository        *RecipeRepository
	bundler_ctx               context.Context
	bundler_ProcessEvaluator  = &mockDetector{}
	bundler_ScriptedEvaluator = &mockDetector{}
	bundler_recipeDetector    = newRecipeDetector(bundler_ProcessEvaluator, bundler_ScriptedEvaluator)
)

func TestBundler_ShouldCreateCore(t *testing.T) {
	bundler_Setup()
	bundler_givenRecipe("id1", types.InfraAgentRecipeName)
	bundler_givenRecipe("id2", types.LoggingRecipeName)
	bundler_givenRecipe("id3", types.GoldenRecipeName)
	bundler_givenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestBundler_CoreShouldDetectAvailableStatus(t *testing.T) {
	bundler_Setup()
	bundler_givenRecipe("id1", types.InfraAgentRecipeName)
	bundler_givenRecipe("id2", types.LoggingRecipeName)
	bundler_givenRecipe("id3", types.GoldenRecipeName)
	bundler_givenRecipe("id4", "mysql")

	bundler := givenBundler()
	with_bundler_recipeDetector(bundler, execution.RecipeStatusTypes.AVAILABLE)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))

	for _, r := range coreBundle.BundleRecipes {
		lastStatusIndex := len(r.Statuses) - 1
		require.Equal(t, 2, len(r.Statuses))
		require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, r.Statuses[lastStatusIndex])
	}
}

func TestBundler_ShouldIncludeDependencies(t *testing.T) {
	bundler_Setup()
	bundler_givenRecipe("id1", types.InfraAgentRecipeName)
	bundler_givenRecipe("id2", types.LoggingRecipeName)
	bundler_givenRecipe("id3", "dep1")
	bundler_givenRecipe("id4", "dep2")

	bundler := givenBundler()

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
	require.Nil(t, findRecipeByName(coreBundle, "dep2"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep1"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep2"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[1], "dep1"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[1], "dep2"))
}

func TestBundler_ShouldCreateEmptyCore(t *testing.T) {
	bundler_Setup()
	bundler_givenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func findRecipeByName(bundle *Bundle, name string) *types.OpenInstallationRecipe {
	for _, r := range bundle.BundleRecipes {
		if strings.EqualFold(r.Recipe.Name, name) {
			return r.Recipe
		}
	}
	return nil
}

func findDependencyByName(recipe *BundleRecipe, name string) *types.OpenInstallationRecipe {
	for _, dep := range recipe.Dependencies {
		if strings.EqualFold(dep.Recipe.Name, name) {
			return dep.Recipe
		}
		found := findDependencyByName(dep, name)
		if found != nil {
			return found
		}
	}
	return nil
}

func bundler_Setup() {
	bundler_discoveryManifest = types.DiscoveryManifest{
		OS: "linux",
	}
	bundler_recipeCache = []types.OpenInstallationRecipe{}
	bundler_repository = NewRecipeRepository(bundler_recipeLoader, &bundler_discoveryManifest)
}

func givenBundler() *Bundler {
	return newBundler(bundler_ctx, bundler_repository, bundler_recipeDetector)
}

func with_bundler_recipeDetector(bundler *Bundler, status execution.RecipeStatusType) {

	bundler_ProcessEvaluator = &mockDetector{status}
	bundler_ScriptedEvaluator = &mockDetector{status}
	bundler_recipeDetector = newRecipeDetector(bundler_ProcessEvaluator, bundler_ScriptedEvaluator)
	bundler.RecipeDetector = bundler_recipeDetector
}

func bundler_recipeLoader() ([]types.OpenInstallationRecipe, error) {
	return bundler_recipeCache, nil
}

func bundler_givenRecipe(id string, name string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	t := types.OpenInstallationRecipeInstallTarget{
		Os: "linux",
	}
	r.InstallTargets = append(r.InstallTargets, t)
	r.Dependencies = []string{"dep1", "dep2", "dep3"}
	bundler_recipeCache = append(bundler_recipeCache, *r)
	return r
}

type mockDetector struct {
	status execution.RecipeStatusType
}

func (d mockDetector) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return d.status
}
