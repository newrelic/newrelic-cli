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
)

func Test_ShouldFindAll_Empty(t *testing.T) {

	repo := NewRecipeRepository(recipeLoader)

	recipes := repo.FindAll(discoveryManifest)

	require.Empty(t, recipes)
}

func Test_matchRecipeCriteria_Basic(t *testing.T) {
	m := types.DiscoveryManifest{
		Platform: "linux",
	}

	hostMap := getHostMap(m)
	actual := matchRecipeCriteriaWhenDefined(hostMap, "Platform", "linux")
	require.True(t, actual)
}

func Test_matchRecipeCriteria_EmptyString(t *testing.T) {
	m := types.DiscoveryManifest{}

	hostMap := getHostMap(m)
	actual := matchRecipeCriteriaWhenDefined(hostMap, "Platform", "")
	require.False(t, actual)
}

func Test_matchRecipeCriteria_KeyMissing(t *testing.T) {
	m := types.DiscoveryManifest{}

	hostMap := getHostMap(m)
	actual := matchRecipeCriteriaWhenDefined(hostMap, "KeyMissing", "")
	require.False(t, actual)
}

func recipeLoader() []types.OpenInstallationRecipe {
	return recipeCache
}
