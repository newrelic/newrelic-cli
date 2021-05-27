package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ConfigValidator interface {
	Validate(ctx context.Context) error
}

type RecipeRecommender interface {
	Recommend(ctx context.Context, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error)
}

type RecipeVarPreparer interface {
	Prepare(m types.DiscoveryManifest, r types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error)
}
