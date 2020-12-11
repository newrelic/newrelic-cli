package install

import (
	"net/url"
)

type mockRecipeFileFetcher struct {
	fetchRecipeFileFunc func(*url.URL) (*recipeFile, error)
	loadRecipeFileFunc  func(string) (*recipeFile, error)
}

func newMockRecipeFileFetcher() *mockRecipeFileFetcher {
	f := mockRecipeFileFetcher{}
	f.fetchRecipeFileFunc = defaultFetchRecipeFileFunc
	f.loadRecipeFileFunc = defaultLoadRecipeFileFunc
	return &f
}

func (f *mockRecipeFileFetcher) fetchRecipeFile(url *url.URL) (*recipeFile, error) {
	return f.fetchRecipeFileFunc(url)
}

func (f *mockRecipeFileFetcher) loadRecipeFile(filename string) (*recipeFile, error) {
	return f.loadRecipeFileFunc(filename)
}

func defaultFetchRecipeFileFunc(recipeURL *url.URL) (*recipeFile, error) {
	return &recipeFile{}, nil
}

func defaultLoadRecipeFileFunc(filename string) (*recipeFile, error) {
	return &recipeFile{}, nil
}
