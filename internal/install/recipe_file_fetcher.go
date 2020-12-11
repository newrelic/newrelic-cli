package install

import (
	"net/url"
)

type recipeFileFetcher interface {
	fetchRecipeFile(recipeURL *url.URL) (*recipeFile, error)
	loadRecipeFile(filename string) (*recipeFile, error)
}
