package validation

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeValidator validates installation of a recipe.
type RecipeValidator interface {
	Validate(context.Context, types.DiscoveryManifest, types.Recipe) (entityGUID string, err error)
}
