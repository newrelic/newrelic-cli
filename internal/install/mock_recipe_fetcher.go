package install

import "context"

type mockRecipeFetcher struct {
	fetchRecommendationsFunc func(*discoveryManifest) ([]recipeFile, error)
	fetchFiltersFunc         func() ([]recipeFilter, error)
}

func newMockRecipeFetcher() *mockRecipeFetcher {
	f := mockRecipeFetcher{}
	f.fetchFiltersFunc = defaultFetchFiltersFunc
	f.fetchRecommendationsFunc = defaultFetchRecommendationsFunc

	return &f
}

func (f *mockRecipeFetcher) fetchRecommendations(ctx context.Context, manifest *discoveryManifest) ([]recipeFile, error) {
	return f.fetchRecommendationsFunc(manifest)
}

func (f *mockRecipeFetcher) fetchFilters(ctx context.Context) ([]recipeFilter, error) {
	return f.fetchFiltersFunc()
}

func defaultFetchFiltersFunc() ([]recipeFilter, error) {
	return []recipeFilter{}, nil
}

func defaultFetchRecommendationsFunc(manifest *discoveryManifest) ([]recipeFile, error) {
	return []recipeFile{}, nil
}
