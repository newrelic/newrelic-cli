// +build unit

package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

func TestEmbeddedRecipeFetcher_FetchRecommendations(t *testing.T) {
	f := EmbeddedRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{OS: "linux"})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)

}

func TestEmbeddedRecipeFetcher_FetchRecipe(t *testing.T) {
	f := EmbeddedRecipeFetcher{}
	recipe, err := f.FetchRecipe(context.Background(), &types.DiscoveryManifest{}, "not-a-real-recipe")
	require.ErrorIs(t, err, ErrRecipeNotFound)
	require.Nil(t, recipe)

	recipe, err = f.FetchRecipe(context.Background(), &types.DiscoveryManifest{OS: "linux"}, "infrastructure-agent-installer")
	require.NoError(t, err)
	require.Equal(t, "infrastructure-agent-installer", recipe.Name)
}

func TestEmbeddedRecipeFetcher_FetchRecipes_EmptyManifest(t *testing.T) {
	f := EmbeddedRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{})
	require.NoError(t, err)
	require.NotEmpty(t, recipes)
}

func TestEmbeddedRecipeFetcher_FetchRecipes_NonEmptyManifest(t *testing.T) {
	f := EmbeddedRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{OS: "linux"})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)

	recipes, err = f.FetchRecipes(context.Background(), &types.DiscoveryManifest{OS: "windows"})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
}
