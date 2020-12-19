package recipes

import (
	"net/url"
)

type RecipeFileFetcher interface {
	FetchRecipeFile(recipeURL *url.URL) (*RecipeFile, error)
	LoadRecipeFile(filename string) (*RecipeFile, error)
}
