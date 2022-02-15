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
	bundlerDiscoveryManifest types.DiscoveryManifest
	bundlerRecipeCache       []types.OpenInstallationRecipe
	bundlerRepository        *RecipeRepository
	bundlerCtx               context.Context
	bundlerProcessEvaluator  = &mockDetector{}
	bundlerScriptedEvaluator = &mockDetector{}
	bundlerRecipeDetector    = newRecipeDetector(bundlerProcessEvaluator, bundlerScriptedEvaluator)
)

func TestBundlerShouldCreateCore(t *testing.T) {
	bundlerSetup()
	bundlerGivenRecipe("id1", types.InfraAgentRecipeName)
	bundlerGivenRecipe("id2", types.LoggingRecipeName)
	bundlerGivenRecipe("id3", types.GoldenRecipeName)
	bundlerGivenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestBundlerCoreShouldDetectAvailableStatus(t *testing.T) {
	bundlerSetup()
	bundlerGivenRecipe("id1", types.InfraAgentRecipeName)
	bundlerGivenRecipe("id2", types.LoggingRecipeName)
	bundlerGivenRecipe("id3", types.GoldenRecipeName)
	bundlerGivenRecipe("id4", "mysql")

	bundler := givenBundler()
	withBundlerRecipeDetector(bundler, execution.RecipeStatusTypes.AVAILABLE)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))

	for _, r := range coreBundle.BundleRecipes {
		lastStatusIndex := len(r.Statuses) - 1
		require.Equal(t, 2, len(r.Statuses))
		require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, r.Statuses[lastStatusIndex])
	}
}

func TestBundlerShouldIncludeDependencies(t *testing.T) {
	bundlerSetup()
	bundlerGivenRecipe("id1", types.InfraAgentRecipeName)
	bundlerGivenRecipe("id2", types.LoggingRecipeName)
	bundlerGivenRecipe("id3", "dep1")
	bundlerGivenRecipe("id4", "dep2")

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

func TestBundlerShouldCreateEmptyCore(t *testing.T) {
	bundlerSetup()
	bundlerGivenRecipe("id4", "mysql")

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

func bundlerSetup() {
	bundlerDiscoveryManifest = types.DiscoveryManifest{
		OS: "linux",
	}
	bundlerRecipeCache = []types.OpenInstallationRecipe{}
	bundlerRepository = NewRecipeRepository(bundlerRecipeLoader, &bundlerDiscoveryManifest)
}

func givenBundler() *Bundler {
	return newBundler(bundlerCtx, bundlerRepository, bundlerRecipeDetector)
}

func withBundlerRecipeDetector(bundler *Bundler, status execution.RecipeStatusType) {

	bundlerProcessEvaluator = &mockDetector{status}
	bundlerScriptedEvaluator = &mockDetector{status}
	bundlerRecipeDetector = newRecipeDetector(bundlerProcessEvaluator, bundlerScriptedEvaluator)
	bundler.RecipeDetector = bundlerRecipeDetector
}

func bundlerRecipeLoader() ([]types.OpenInstallationRecipe, error) {
	return bundlerRecipeCache, nil
}

func bundlerGivenRecipe(id string, name string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	t := types.OpenInstallationRecipeInstallTarget{
		Os: "linux",
	}
	r.InstallTargets = append(r.InstallTargets, t)
	r.Dependencies = []string{"dep1", "dep2", "dep3"}
	bundlerRecipeCache = append(bundlerRecipeCache, *r)
	return r
}

type mockDetector struct {
	status execution.RecipeStatusType
}

func (d mockDetector) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return d.status
}
