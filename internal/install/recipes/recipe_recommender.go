package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeRecommender interface {
	Recommend(ctx context.Context, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error)
}
