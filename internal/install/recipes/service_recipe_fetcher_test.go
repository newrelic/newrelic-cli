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

	c := NewMockNerdGraphClient()
	c.RespBody = wrapRecipes(r)

	s := NewServiceRecipeFetcher(c)

	recipes, err := s.FetchRecipes(context.Background())
	require.NoError(t, err)
	require.NotNil(t, recipes)
	require.NotEmpty(t, recipes)
	// The duplicate name should still be included when we perform the FetchRecipes() call.
	require.Equal(t, 3, len(recipes))
	require.True(t, reflect.DeepEqual(r, recipes))
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
