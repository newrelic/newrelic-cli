package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

func TestLocalRecipeFetcher_FetchRecipe(t *testing.T) {

}

func TestLocalRecipeFetcher_FetchRecommendations(t *testing.T) {

}

func TestLocalRecipeFetcher_FetchRecipes(t *testing.T) {
	f := LocalRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
}
