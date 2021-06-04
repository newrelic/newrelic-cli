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

func recipeLoader() []types.OpenInstallationRecipe {
	return recipeCache
}
