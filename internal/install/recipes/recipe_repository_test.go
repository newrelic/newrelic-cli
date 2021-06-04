// build +unit

package recipes

import (
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var (
	discoveryManifest types.DiscoveryManifest        = types.DiscoveryManifest{}
	recipeCache       []types.OpenInstallationRecipe = []types.OpenInstallationRecipe{}
)

func Test_ShouldFindAll_Empty(t *testing.T) {
	repo := NewRecipeRepository(recipeLoader)

	recipes, _ := repo.FindAll(discoveryManifest)

	require.Empty(t, recipes)
}

func Test_ShouldFindSingleRecipe(t *testing.T) {
	givenCachedRecipe("id1", "my-recipe")

	repo := NewRecipeRepository(recipeLoader)
	results, _ := repo.FindAll(discoveryManifest)

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
	log.Debugf("Test loading %d recipes", len(recipeCache))
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
