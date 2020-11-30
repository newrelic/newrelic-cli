package install

type recipeFetcher interface {
	fetchRecommendations(*discoveryManifest) ([]recipeFile, error)
	fetchFilters() ([]recipeFilter, error)
}
