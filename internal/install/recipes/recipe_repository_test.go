// build +unit

package recipes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	discoveryManifest types.DiscoveryManifest        = types.DiscoveryManifest{}
	recipeCache       []types.OpenInstallationRecipe = []types.OpenInstallationRecipe{}
	repository        *RecipeRepository              = NewRecipeRepository(recipeLoader)
)

func Test_ShouldFindAll_Empty(t *testing.T) {

	recipes, _ := repository.FindAll(discoveryManifest)

	require.Empty(t, recipes)
}

func Test_ShouldFindSingleRecipe(t *testing.T) {
	givenCachedRecipe("id1", "my-recipe")

	results, _ := repository.FindAll(discoveryManifest)

	require.Len(t, results, 1)
	require.Equal(t, results[0].ID, "id1")
}

func Test_matchRecipeCriteria_Basic(t *testing.T) {
	m := types.DiscoveryManifest{
		Platform: "linux",
	}

	hostMap := getHostMap(m)
	actual := matchRecipeCriteria(hostMap, "Platform", "linux")
	require.True(t, actual)
}

func Test_matchRecipeCriteria_EmptyString(t *testing.T) {
	m := types.DiscoveryManifest{}

	hostMap := getHostMap(m)
	actual := matchRecipeCriteria(hostMap, "Platform", "")
	require.True(t, actual)
}

func Test_matchRecipeCriteria_KeyMissing(t *testing.T) {
	m := types.DiscoveryManifest{}

	hostMap := getHostMap(m)
	actual := matchRecipeCriteria(hostMap, "KeyMissing", "xyz")
	require.False(t, actual)
}

func recipeLoader() ([]types.OpenInstallationRecipe, error) {
	return recipeCache, nil
}

func givenCachedRecipe(id string, name string) *types.OpenInstallationRecipe {
	r := createRecipe(id, name)
	recipeCache = append(recipeCache, *r)
	return r
}

func createRecipe(id string, name string) *types.OpenInstallationRecipe {
	r := &types.OpenInstallationRecipe{
		ID:   id,
		Name: name,
	}
	return r
}

// const givenCachedRecipe = function (id, name, os = null, platform = null, platformVersion = null, kernelArch = null) {
//     const recipe = createRecipe(id, name, os, platform, platformVersion, kernelArch);
//     service.addSingleRecipe(nodeCache, recipe);
//   }

//   const createRecipe = function (id, name, os = null, platform = null, platformVersion = null, kernelArch = null) {
//     const recipe = new Recipe();
//     recipe.id = id;
//     recipe.name = name;
//     if (os != null || platform != null) {
//       const target = createInstallTarget(os, platform, platformVersion, kernelArch);
//       recipe.installTargets = [target];
//     }
//     return recipe;
//   }
