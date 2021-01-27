// +build unit

package recipes

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
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
		// Include a duplicate name
		{
			ID:   "MAo=",
			Name: "test",
			ProcessMatch: []string{
				"test",
			},
		},
		{
			ID:   "MAo=",
			Name: "othername",
			ProcessMatch: []string{
				"test",
			},
		},
	}

	c := newMockNerdGraphClient()
	c.respBody = wrapRecipes(r)

	s := NewServiceRecipeFetcher(c)

	recipes, err := s.FetchRecipes(context.Background(), &types.DiscoveryManifest{})
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	// The duplicate name should still be included when we perform the FetchRecipes() call.
	require.Equal(t, 3, len(recipes))
	require.True(t, reflect.DeepEqual(createRecipes(r), recipes))
}

func TestFetchRecommendations(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "MAo=",
			Name: "testing1",
			File: `
---
name: Test recipe file
description: test description
`,
		},
		{
			ID:   "non-zero",
			Name: "testing1",
			File: `
---
name: Test recipe file2
description: test description
`,
		},
		{
			ID:   "non-zero2",
			Name: "testing2",
			File: `
---
name: Test recipe file3
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
	// The duplicate name should be removed from the set when we fetch the reecommendations.
	require.Equal(t, 2, len(recipes))
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
