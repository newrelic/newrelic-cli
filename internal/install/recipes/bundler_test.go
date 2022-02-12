package recipes

import (
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"

	"strings"

	"github.com/stretchr/testify/require"
)

var (
	bundler_discoveryManifest types.DiscoveryManifest
	bundler_recipeCache       []types.OpenInstallationRecipe
	bundler_repository        *RecipeRepository
	bundler_detector          *RecipeDetector
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

func TestBundler_ShouldIncludeDependencies(t *testing.T) {
	bundler_Setup()
	bundler_givenRecipe("id1", types.InfraAgentRecipeName)
	bundler_givenRecipe("id2", types.LoggingRecipeName)
	bundler_givenRecipe("id3", "dep1")
	bundler_givenRecipe("id4", "dep2")

	bundler := givenBundler()
	coreBundle := bundler.CreateCoreBundle()

	t.Log(coreBundle)

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
	for _, r := range recipe.Dependencies {
		if strings.EqualFold(r.Recipe.Name, name) {
			return r.Recipe
		}
		return findDependencyByName(r, name)
	}
	return nil
}

func bundler_Setup() {
	bundler_discoveryManifest = types.DiscoveryManifest{
		OS: "linux",
	}
	bundler_recipeCache = []types.OpenInstallationRecipe{}
	bundler_repository = NewRecipeRepository(bundler_recipeLoader, &bundler_discoveryManifest)
	bundler_detector = NewRecipeDetector()
}

func givenBundler() *Bundler {
	return NewBundler(bundler_repository)
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
