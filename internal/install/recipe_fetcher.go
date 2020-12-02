package install

import "context"

type recipeFetcher interface {
	fetchRecipe(context.Context, *discoveryManifest, string) (*recipe, error)
	fetchRecommendations(context.Context, *discoveryManifest) ([]recipe, error)
	fetchRecipes(context.Context) ([]recipe, error)
}
