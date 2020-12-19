// +build unit

package recipes

import (
	"context"
	"reflect"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

func TestFetchFilters(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "MAo=",
			Name: "test",
			ProcessMatch: []string{
				"test",
			},
		},
	}

	c := newMockNerdGraphClient()
	c.respBody = wrapRecipes(r)

	s := NewServiceRecipeFetcher(c)

	recipes, err := s.FetchRecipes(context.Background())
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
	require.True(t, reflect.DeepEqual(createRecipes(r), recipes))
}

func TestFetchRecommendations(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID: "MAo=",
			File: `
---
name: Test recipe file
description: test description
`,
		},
	}

	c := newMockNerdGraphClient()
	c.respBody = wrapRecommendations(r)

	m := types.DiscoveryManifest{}

	s := NewServiceRecipeFetcher(c)

	recipes, err := s.FetchRecommendations(context.Background(), &m)
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
}

func wrapRecipes(r []types.OpenInstallationRecipe) recipeSearchQueryResult {
	return recipeSearchQueryResult{
		Docs: recipeSearchQueryDocs{
			OpenInstallation: recipeSearchQueryOpenInstallation{
				RecipeSearch: recipeSearchResult{
					Results: r,
				},
			},
		},
	}
}

func wrapRecommendations(r []types.OpenInstallationRecipe) recommendationsQueryResult {
	return recommendationsQueryResult{
		Docs: recommendationsQueryDocs{
			OpenInstallation: recommendationsQueryOpenInstallation{
				Recommendations: recommendationsResult{
					Results: r,
				},
			},
		},
	}
}
