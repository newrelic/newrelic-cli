package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeFetcher is responsible for retrieving recipe information.
type RecipeFetcher interface {
	FetchRecipes(context.Context) ([]*types.OpenInstallationRecipe, error)
	FetchLibraryVersion(context.Context) string
}
