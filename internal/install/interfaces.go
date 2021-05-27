package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ConfigValidator interface {
	ValidateConfig(ctx context.Context) error
}

type RecipeRecommender interface {
	Recommend(ctx context.Context, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error)
}
