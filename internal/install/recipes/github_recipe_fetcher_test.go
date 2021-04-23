// +build integration

package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestGithubRecipeFetcher_FetchRecommendations(t *testing.T) {
	f := GithubRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{OS: "linux"})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
}

func TestGithubRecipeFetcher_FetchRecipe(t *testing.T) {
	f := GithubRecipeFetcher{}
	recipe, err := f.FetchRecipe(context.Background(), &types.DiscoveryManifest{}, "not-a-real-recipe")
	require.ErrorIs(t, err, ErrRecipeNotFound)
	require.Nil(t, recipe)

	recipe, err = f.FetchRecipe(context.Background(), &types.DiscoveryManifest{OS: "linux"}, "infrastructure-agent-installer")
	require.NoError(t, err)
	require.Equal(t, "infrastructure-agent-installer", recipe.Name)
}

func TestGithubRecipeFetcher_FetchRecipes_EmptyManifest(t *testing.T) {
	f := GithubRecipeFetcher{}
	recipes, err := f.FetchRecipes(context.Background(), &types.DiscoveryManifest{})
	require.NoError(t, err)
	require.NotEmpty(t, recipes)
}

func TestGithubRecipeFetcher_cacheLatestRelease(t *testing.T) {
	err := cacheLatestRelease(context.Background())
	require.NoError(t, err)
}

func TestGithubRecipeFetcher_getLatestRelease(t *testing.T) {
	release, url, err := getLatestRelease(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, release)
	require.NotNil(t, url)
}
