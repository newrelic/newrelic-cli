package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type MockRecipeFetcher struct {
	FetchRecipeErr                error
	FetchRecipesErr               error
	FetchRecommendationsErr       error
	FetchRecipeCallCount          int
	FetchRecipesCallCount         int
	FetchRecommendationsCallCount int
	FetchRecipeVals               []types.Recipe
	FetchRecipeVal                *types.Recipe
	FetchRecipesVal               []types.Recipe
	FetchRecommendationsVal       []types.Recipe
}

func NewMockRecipeFetcher() *MockRecipeFetcher {
	f := MockRecipeFetcher{}
	f.FetchRecipesVal = []types.Recipe{}
	f.FetchRecommendationsVal = []types.Recipe{}
	return &f
}

func (f *MockRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.Recipe, error) {
	f.FetchRecipeCallCount++

	if len(f.FetchRecipeVals) > 0 {
		i := utils.MinOf(f.FetchRecipeCallCount, len(f.FetchRecipeVals)) - 1
		return &f.FetchRecipeVals[i], f.FetchRecipesErr
	}

	return f.FetchRecipeVal, f.FetchRecipeErr
}

func (f *MockRecipeFetcher) FetchRecipes(ctx context.Context) ([]types.Recipe, error) {
	f.FetchRecipesCallCount++
	return f.FetchRecipesVal, f.FetchRecipesErr
}

func (f *MockRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {
	f.FetchRecommendationsCallCount++
	return f.FetchRecommendationsVal, f.FetchRecommendationsErr
}
