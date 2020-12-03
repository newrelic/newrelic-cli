package install

import "context"

type mockRecipeFetcher struct {
	fetchRecipeFunc          func(*discoveryManifest, string) (*recipe, error)
	fetchRecipesFunc         func() ([]recipe, error)
	fetchRecommendationsFunc func(*discoveryManifest) ([]recipe, error)
}

func newMockRecipeFetcher() *mockRecipeFetcher {
	f := mockRecipeFetcher{}
	f.fetchRecipeFunc = defaultFetchRecipeFunc
	f.fetchRecipesFunc = defaultFetchRecipesFunc
	f.fetchRecommendationsFunc = defaultFetchRecommendationsFunc

	return &f
}

func (f *mockRecipeFetcher) fetchRecipe(ctx context.Context, manifest *discoveryManifest, friendlyName string) (*recipe, error) {
	return f.fetchRecipeFunc(manifest, friendlyName)
}

func (f *mockRecipeFetcher) fetchRecipes(ctx context.Context) ([]recipe, error) {
	return f.fetchRecipesFunc()
}

func (f *mockRecipeFetcher) fetchRecommendations(ctx context.Context, manifest *discoveryManifest) ([]recipe, error) {
	return f.fetchRecommendationsFunc(manifest)
}

func defaultFetchRecipeFunc(manifest *discoveryManifest, friendlyName string) (*recipe, error) {
	return &recipe{}, nil
}

func defaultFetchRecommendationsFunc(manifest *discoveryManifest) ([]recipe, error) {
	return []recipe{}, nil
}

func defaultFetchRecipesFunc() ([]recipe, error) {
	return []recipe{}, nil
}
