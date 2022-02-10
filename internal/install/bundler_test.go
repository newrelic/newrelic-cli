package install

import (
	"testing"

	recipes "github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"strings"

	"github.com/stretchr/testify/require"
)

var (
	discoveryManifest types.DiscoveryManifest
	recipeCache       []types.OpenInstallationRecipe
	repository        *recipes.RecipeRepository
)

func TestBundler_ShouldCreateCore(t *testing.T) {
	Setup()
	givenRecipe("id1", types.InfraAgentRecipeName)
	givenRecipe("id2", types.LoggingRecipeName)
	givenRecipe("id3", types.GoldenRecipeName)
	givenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.createCoreBundle()

	require.Equal(t, 3, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestBundler_ShouldIncludeDependencies(t *testing.T) {
	Setup()
	givenRecipe("id1", types.InfraAgentRecipeName)
	givenRecipe("id2", types.LoggingRecipeName)
	givenRecipe("id3", "dep1")
	givenRecipe("id4", "dep2")

	bundler := givenBundler()
	coreBundle := bundler.createCoreBundle()

	t.Log(coreBundle)

	require.Equal(t, 4, len(coreBundle.BundleRecipes))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, "dep1"))
	require.NotNil(t, findRecipeByName(coreBundle, "dep2"))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestBundler_ShouldCreateEmptyCore(t *testing.T) {
	Setup()
	givenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.createCoreBundle()

	require.Equal(t, 0, len(coreBundle.BundleRecipes))
}

func findRecipeByName(bundle *Bundle, name string) *types.OpenInstallationRecipe {
	for _, r := range bundle.BundleRecipes {
		if strings.EqualFold(r.recipe.Name, name) {
			return r.recipe
		}
	}
	return nil
}

func Setup() {
	discoveryManifest = types.DiscoveryManifest{
		OS: "linux",
	}
	recipeCache = []types.OpenInstallationRecipe{}
	repository = recipes.NewRecipeRepository(recipeLoader, &discoveryManifest)
}

func givenBundler() *Bundler {
	return NewBundler(repository)
}

func recipeLoader() ([]types.OpenInstallationRecipe, error) {
	return recipeCache, nil
}

func givenRecipe(id string, name string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	t := types.OpenInstallationRecipeInstallTarget{
		Os: "linux",
	}
	r.InstallTargets = append(r.InstallTargets, t)
	r.Dependencies = []string{"dep1", "dep2", "dep3"}
	recipeCache = append(recipeCache, *r)
	return r
}
