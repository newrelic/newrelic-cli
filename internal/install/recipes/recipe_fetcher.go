package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeFetcher is responsible for retrieving recipe information.
type RecipeFetcher interface {
	FetchRecipe(context.Context, *types.DiscoveryManifest, string) (*types.OpenInstallationRecipe, error)
	FetchRecommendations(context.Context, *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error)
	FetchRecipes(context.Context, *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error)
}
