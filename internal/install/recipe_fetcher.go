package install

import "context"

type recipeFetcher interface {
	fetchRecommendations(context.Context, *discoveryManifest) ([]recipeFile, error)
	fetchFilters(context.Context) ([]recipeFilter, error)
}
