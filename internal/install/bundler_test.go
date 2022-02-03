//go:build unit
// +build unit

package install

import (
	"testing"

	recipes "github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/require"
	"strings"
)

func TestCreateBundles_ShouldCreateCoreBundle(t *testing.T) {
	Setup()
	givenRecipe("id1", types.InfraAgentRecipeName)
	givenRecipe("id2", types.LoggingRecipeName)
	givenRecipe("id3", types.GoldenRecipeName)
	givenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.createCoreBundle()

	require.Equal(t, 3, len(coreBundle))
	require.NotNil(t, findRecipeByName(coreBundle, types.InfraAgentRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.LoggingRecipeName))
	require.NotNil(t, findRecipeByName(coreBundle, types.GoldenRecipeName))
	require.Nil(t, findRecipeByName(coreBundle, "mysql"))
}

func TestCreateBundles_ShouldCreateEmptyCoreBundle(t *testing.T) {
	Setup()
	givenRecipe("id4", "mysql")

	bundler := givenBundler()
	coreBundle := bundler.createCoreBundle()

	require.Equal(t, 0, len(coreBundle))
}

func findRecipeByName(recipes []types.OpenInstallationRecipe, name string) *types.OpenInstallationRecipe {
	for _, r := range recipes {
		if strings.EqualFold(r.Name, name) {
			return &r
		}
	}
	return nil
}

var (
	discoveryManifest types.DiscoveryManifest
	recipeCache       []types.OpenInstallationRecipe
	repository        *recipes.RecipeRepository
)

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
	recipeCache = append(recipeCache, *r)
	return r
}
