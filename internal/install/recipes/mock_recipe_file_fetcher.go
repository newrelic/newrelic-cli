package recipes

import (
	"net/url"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeFileFetcher struct {
	FetchRecipeFileFunc func(*url.URL) (*types.OpenInstallationRecipe, error)
	LoadRecipeFileFunc  func(string) (*types.OpenInstallationRecipe, error)
}

func NewMockRecipeFileFetcher() *MockRecipeFileFetcher {
	f := MockRecipeFileFetcher{}
	f.FetchRecipeFileFunc = defaultFetchRecipeFileFunc
	f.LoadRecipeFileFunc = defaultLoadRecipeFileFunc
	return &f
}

func (f *MockRecipeFileFetcher) FetchRecipeFile(url *url.URL) (*types.OpenInstallationRecipe, error) {
	return f.FetchRecipeFileFunc(url)
}

func (f *MockRecipeFileFetcher) LoadRecipeFile(filename string) (*types.OpenInstallationRecipe, error) {
	return f.LoadRecipeFileFunc(filename)
}

func defaultFetchRecipeFileFunc(recipeURL *url.URL) (*types.OpenInstallationRecipe, error) {
	return &types.OpenInstallationRecipe{}, nil
}

func defaultLoadRecipeFileFunc(filename string) (*types.OpenInstallationRecipe, error) {
	return &types.OpenInstallationRecipe{}, nil
}
