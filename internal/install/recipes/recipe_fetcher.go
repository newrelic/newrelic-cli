package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeFetcher is responsible for retrieving recipe information.
type RecipeFetcher interface {
	FetchRecipe(context.Context, *types.DiscoveryManifest, string) (*types.Recipe, error)
	FetchRecommendations(context.Context, *types.DiscoveryManifest) ([]types.Recipe, error)
	FetchRecipes(context.Context, *types.DiscoveryManifest) ([]types.Recipe, error)
}
