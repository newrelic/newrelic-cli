package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeFetcher struct {
	FetchRecipesErr       error
	FetchRecipesCallCount int
	FetchRecipesVal       []types.OpenInstallationRecipe
}

func NewMockRecipeFetcher() *MockRecipeFetcher {
	f := MockRecipeFetcher{}
	f.FetchRecipesVal = []types.OpenInstallationRecipe{}
	return &f
}

func (f *MockRecipeFetcher) FetchRecipes(ctx context.Context) ([]types.OpenInstallationRecipe, error) {
	f.FetchRecipesCallCount++
	return f.FetchRecipesVal, f.FetchRecipesErr
}

func (f *MockRecipeFetcher) FetchLibraryVersion(ctx context.Context) string {
	return ""
}
