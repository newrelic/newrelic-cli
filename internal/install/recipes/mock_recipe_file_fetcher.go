package recipes

import (
	"net/url"
)

type MockRecipeFileFetcher struct {
	FetchRecipeFileFunc func(*url.URL) (*RecipeFile, error)
	LoadRecipeFileFunc  func(string) (*RecipeFile, error)
}

func NewMockRecipeFileFetcher() *MockRecipeFileFetcher {
	f := MockRecipeFileFetcher{}
	f.FetchRecipeFileFunc = defaultFetchRecipeFileFunc
	f.LoadRecipeFileFunc = defaultLoadRecipeFileFunc
	return &f
}

func (f *MockRecipeFileFetcher) FetchRecipeFile(url *url.URL) (*RecipeFile, error) {
	return f.FetchRecipeFileFunc(url)
}

func (f *MockRecipeFileFetcher) LoadRecipeFile(filename string) (*RecipeFile, error) {
	return f.LoadRecipeFileFunc(filename)
}

func defaultFetchRecipeFileFunc(recipeURL *url.URL) (*RecipeFile, error) {
	return &RecipeFile{}, nil
}

func defaultLoadRecipeFileFunc(filename string) (*RecipeFile, error) {
	return &RecipeFile{}, nil
}
