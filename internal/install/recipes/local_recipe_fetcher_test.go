package recipes

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

func TestLocalRecipeFetcher_FetchRecommendations(t *testing.T) {

}

func TestLocalRecipeFetcher_FetchRecipe(t *testing.T) {
	f := LocalRecipeFetcher{}
	recipe, err := f.FetchRecipe(context.Background(), &types.DiscoveryManifest{}, "not-a-real-recipe")
	require.ErrorIs(t, err, ErrRecipeNotFound)
	require.Nil(t, recipe)
}

func TestLocalRecipeFetcher_FetchRecipes_EmptyManifest(t *testing.T) {
	f := LocalRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{})
	require.NoError(t, err)
	require.Empty(t, recipes)
}

func TestLocalRecipeFetcher_FetchRecipes_NonEmptyManifest(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := LocalRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{OS: "linux"})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
}
