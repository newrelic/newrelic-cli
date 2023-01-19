package validation

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeValidator validates installation of a recipe.
type RecipeValidator interface {
	ValidateRecipe(context.Context, types.DiscoveryManifest, types.OpenInstallationRecipe, types.RecipeVars) (entityGUID string, err error)
}

type AgentValidator interface {
	Validate(ctx context.Context, url string) (string, error)
}
