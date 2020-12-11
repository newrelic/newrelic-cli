// +build unit

package install

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchFilters(t *testing.T) {
	r := []OpenInstallationRecipe{
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

	s := newServiceRecipeFetcher(c)

	recipes, err := s.fetchRecipes(context.Background())
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
	require.True(t, reflect.DeepEqual(createRecipes(r), recipes))
}

func TestFetchRecommendations(t *testing.T) {
	r := []OpenInstallationRecipe{
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

	m := discoveryManifest{}

	s := newServiceRecipeFetcher(c)

	recipes, err := s.fetchRecommendations(context.Background(), &m)
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
}

func wrapRecipes(r []OpenInstallationRecipe) recipeSearchQueryResult {
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

func wrapRecommendations(r []OpenInstallationRecipe) recommendationsQueryResult {
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
