// +build unit

package install

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchFilters(t *testing.T) {
	f := []recipeFilter{
		{
			ID: "test",
			Metadata: recipeFilterMetadata{
				Name: "test",
				ProcessMatch: []string{
					"test",
				},
			},
		},
	}

	c := newMockNerdGraphClient()
	c.respBody = wrapFilters(f)

	s := newServiceRecipeFetcher(c)

	filters, err := s.fetchFilters()
	require.NoError(t, err)
	require.NotNil(t, filters)
	require.NotEmpty(t, filters)
	require.Equal(t, 1, len(filters))
	require.True(t, reflect.DeepEqual(f, filters))
}

func TestFetchRecommendations(t *testing.T) {
	r := []recipe{
		{
			ID:       "test",
			Metadata: recipeMetadata{},
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

	recipes, err := s.fetchRecommendations(&m)
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	require.Equal(t, 1, len(recipes))
}

func wrapFilters(f []recipeFilter) recipeFilterQueryResult {
	return recipeFilterQueryResult{
		Account: recipeFilterQueryAccount{
			OpenInstallation: recipeFilterQueryOpenInstallation{
				RecipeSearch: recipeFilterResult{
					Results: f,
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
