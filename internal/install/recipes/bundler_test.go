package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"strings"

	"github.com/stretchr/testify/require"
)

type bundlerTest struct {
	discoveryManifest types.DiscoveryManifest
	recipeCache       []types.OpenInstallationRecipe
	recipeRepository  *RecipeRepository
	ctx               context.Context
	processEvaluator  *mockDetector
	scriptedEvaluator *mockDetector
	recipeDetector    *RecipeDetector
}

type mockDetector struct {
	status execution.RecipeStatusType
}

var (
	bundlerTestImpl = bundlerTest{
		processEvaluator:  &mockDetector{},
		scriptedEvaluator: &mockDetector{},
	}
)

func TestBundlerShouldCreateCoreBundle(t *testing.T) {
	setup()
	addRecipeToCache("id1", types.InfraAgentRecipeName)
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestBundlerShouldCreateAdditionalBundle(t *testing.T) {
	setup()
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()

	coreBundle := bundler.CreateAdditionalBundle()

	require.Equal(t, 1, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestBundlerCoreShouldDetectAvailableStatus(t *testing.T) {
	setup()
	addRecipeToCache("id1", types.InfraAgentRecipeName)
	addRecipeToCache("id2", types.LoggingRecipeName)
	addRecipeToCache("id3", types.GoldenRecipeName)
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()
	withRecipeStatusDetector(bundler, execution.RecipeStatusTypes.AVAILABLE)

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))

	for _, r := range coreBundle.BundleRecipes {
		lastStatusIndex := len(r.RecipeStatuses) - 1
		require.Equal(t, 2, len(r.RecipeStatuses))
		require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, r.RecipeStatuses[lastStatusIndex].Status)
	}
}

func TestBundlerShouldIncludeDependencies(t *testing.T) {
	setup()
	bundlerTestImpl.addRecipeWithDependenciesToCache("id1", types.InfraAgentRecipeName, []string{"dep1", "dep2"})
	bundlerTestImpl.addRecipeWithDependenciesToCache("id2", types.LoggingRecipeName, []string{"dep2"})
	addRecipeToCache("id3", "dep1")
	addRecipeToCache("id4", "dep2")
	bundler := createTestBundler()

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 2, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
	require.Nil(t, findRecipeByName(coreBundle, "dep2"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep1"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep2"))
	require.Nil(t, findDependencyByName(coreBundle.BundleRecipes[1], "dep1"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[1], "dep2"))
}

func TestBundlerShouldNotIncludeInvalidDependencies(t *testing.T) {
	setup()
	bundlerTestImpl.addRecipeWithDependenciesToCache("id1", types.InfraAgentRecipeName, []string{"dep1", "dep2", "dep3"})
	addRecipeToCache("id2", "dep1")
	addRecipeToCache("id3", "dep2")
	bundler := createTestBundler()

	coreBundle := bundler.CreateCoreBundle()

	require.Equal(t, 1, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "dep1"))
	require.Nil(t, findRecipeByName(coreBundle, "dep2"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep1"))
	require.NotNil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep2"))
	require.Nil(t, findDependencyByName(coreBundle.BundleRecipes[0], "dep3"))
}

func TestBundlerShouldCreateEmptyCore(t *testing.T) {
	setup()
	addRecipeToCache("id4", "mysql")
	bundler := createTestBundler()

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

func setup() {
	bundlerTestImpl.recipeDetector = newRecipeDetector(bundlerTestImpl.processEvaluator, bundlerTestImpl.scriptedEvaluator)
	bundlerTestImpl.discoveryManifest = types.DiscoveryManifest{
		OS: "linux",
	}
	bundlerTestImpl.recipeCache = []types.OpenInstallationRecipe{}
	bundlerTestImpl.recipeRepository = NewRecipeRepository(bundlerRecipeLoader, &bundlerTestImpl.discoveryManifest)
}

func createTestBundler() *Bundler {
	return newBundler(bundlerTestImpl.ctx, bundlerTestImpl.recipeRepository, bundlerTestImpl.recipeDetector)
}

func withRecipeStatusDetector(bundler *Bundler, status execution.RecipeStatusType) {
	bundlerTestImpl.processEvaluator = &mockDetector{status}
	bundlerTestImpl.scriptedEvaluator = &mockDetector{status}
	bundlerTestImpl.recipeDetector = newRecipeDetector(bundlerTestImpl.processEvaluator, bundlerTestImpl.scriptedEvaluator)
	bundler.RecipeDetector = bundlerTestImpl.recipeDetector
}

func bundlerRecipeLoader() ([]types.OpenInstallationRecipe, error) {
	return bundlerTestImpl.recipeCache, nil
}

func addRecipeToCache(id string, name string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	t := types.OpenInstallationRecipeInstallTarget{
		Os: "linux",
	}
	r.InstallTargets = append(r.InstallTargets, t)
	bundlerTestImpl.recipeCache = append(bundlerTestImpl.recipeCache, *r)
	return r
}

func (br *bundlerTest) addRecipeWithDependenciesToCache(id string, name string, dependencies []string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	t := types.OpenInstallationRecipeInstallTarget{
		Os: "linux",
	}
	r.InstallTargets = append(r.InstallTargets, t)

	if len(dependencies) > 0 {
		r.Dependencies = dependencies
	}

	br.recipeCache = append(bundlerTestImpl.recipeCache, *r)
	return r
}

func (d mockDetector) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	return d.status
}
