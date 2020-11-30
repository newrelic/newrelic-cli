package install

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

func (f *mockRecipeFetcher) fetchRecommendations(manifest *discoveryManifest) ([]recipeFile, error) {
	return f.fetchRecommendationsFunc(manifest)
}

func (f *mockRecipeFetcher) fetchFilters() ([]recipeFilter, error) {
	return f.fetchFiltersFunc()
}

func defaultFetchFiltersFunc() ([]recipeFilter, error) {
	return []recipeFilter{}, nil
}

func defaultFetchRecommendationsFunc(manifest *discoveryManifest) ([]recipeFile, error) {
	return []recipeFile{}, nil
}
