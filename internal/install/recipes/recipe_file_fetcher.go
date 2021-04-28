package recipes

import (
	"net/url"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeFileFetcher interface {
	FetchRecipeFile(recipeURL *url.URL) (*types.OpenInstallationRecipe, error)
	LoadRecipeFile(filename string) (*types.OpenInstallationRecipe, error)
}
