package install

import (
	"context"
)

type mockRecipeFetcher struct {
	fetchRecipeErr                error
	fetchRecipesErr               error
	fetchRecommendationsErr       error
	fetchRecipeCallCount          int
	fetchRecipesCallCount         int
	fetchRecommendationsCallCount int
	fetchRecipeVals               []recipe
	fetchRecipeVal                *recipe
	fetchRecipesVal               []recipe
	fetchRecommendationsVal       []recipe
}

func newMockRecipeFetcher() *mockRecipeFetcher {
	f := mockRecipeFetcher{}
	f.fetchRecipesVal = []recipe{}
	f.fetchRecommendationsVal = []recipe{}
	return &f
}

func (f *mockRecipeFetcher) fetchRecipe(ctx context.Context, manifest *discoveryManifest, friendlyName string) (*recipe, error) {
	f.fetchRecipeCallCount++

	if len(f.fetchRecipeVals) > 0 {
		i := minOf(f.fetchRecipeCallCount, len(f.fetchRecipeVals)) - 1
		return &f.fetchRecipeVals[i], f.fetchRecipesErr
	}

	return f.fetchRecipeVal, f.fetchRecipeErr
}

func (f *mockRecipeFetcher) fetchRecipes(ctx context.Context) ([]recipe, error) {
	f.fetchRecipesCallCount++
	return f.fetchRecipesVal, f.fetchRecipesErr
}

func (f *mockRecipeFetcher) fetchRecommendations(ctx context.Context, manifest *discoveryManifest) ([]recipe, error) {
	f.fetchRecommendationsCallCount++
	return f.fetchRecommendationsVal, f.fetchRecommendationsErr
}

func minOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}

	return min
}
