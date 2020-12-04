// +build unit

package install

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchFilters(t *testing.T) {
	r := []recipe{
		{
			ID:   "test",
			Name: "test",
			ProcessMatch: []string{
				"test",
			},
		},
	}

	c := newMockNerdGraphClient()
	c.respBody = wrapRecipes(r)

	s := newServiceRecipeFetcher(c)

	recipes, err := s.fetchRecipes(context.Background())
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
	require.True(t, reflect.DeepEqual(r, recipes))
}

func TestFetchRecommendations(t *testing.T) {
	r := []recipe{
		{
			ID: "test",
			File: `
---
name: Test recipe file
description: test description
`,
		},
	}

	c := newMockNerdGraphClient()
	c.respBody = wrapRecommendations(r)

	m := discoveryManifest{}

	s := newServiceRecipeFetcher(c)

	recipes, err := s.fetchRecommendations(context.Background(), &m)
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
}

func wrapRecipes(r []recipe) recipeSearchQueryResult {
	return recipeSearchQueryResult{
		Account: recipeSearchQueryAccount{
			OpenInstallation: recipeSearchQueryOpenInstallation{
				RecipeSearch: recipeSearchResult{
					Results: r,
				},
			},
		},
	}
}

func wrapRecommendations(r []recipe) recommendationsQueryResult {
	return recommendationsQueryResult{
		Account: recommendationsQueryAccount{
			OpenInstallation: recommendationsQueryOpenInstallation{
				Recommendations: recommendationsResult{
					Results: r,
				},
			},
		},
	}
}
