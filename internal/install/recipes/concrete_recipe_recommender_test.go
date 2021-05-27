package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecommend_CustomScript_Success(t *testing.T) {
	rf := NewMockRecipeFetcher()
	rf.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name: "test-recipe",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				RequireAtDiscovery: "echo 1234",
			},
		},
	}
	pf := NewMockProcessFilterer()
	re := execution.NewShRecipeExecutor()

	m := &types.DiscoveryManifest{}

	r := NewConcreteRecipeRecommender(rf, pf, re)

	recipes, err := r.Recommend(context.Background(), m)
	require.NoError(t, err)
	require.Equal(t, 1, len(recipes))
}

func TestRecommend_CustomScript_Failure(t *testing.T) {
	rf := NewMockRecipeFetcher()
	rf.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name: "test-recipe",
			PreInstall: types.OpenInstallationPreInstallConfiguration{
				RequireAtDiscovery: "bogus command",
			},
		},
	}
	pf := NewMockProcessFilterer()
	re := execution.NewShRecipeExecutor()

	m := &types.DiscoveryManifest{}

	r := NewConcreteRecipeRecommender(rf, pf, re)

	recipes, err := r.Recommend(context.Background(), m)
	require.NoError(t, err)
	require.Equal(t, 0, len(recipes))
}
