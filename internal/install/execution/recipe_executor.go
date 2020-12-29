package execution

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeExecutor is responsible for execution of the task steps defined in a
// recipe.
type RecipeExecutor interface {
	Prepare(context.Context, types.DiscoveryManifest, types.Recipe) (types.RecipeVars, error)
	Execute(context.Context, types.DiscoveryManifest, types.Recipe, types.RecipeVars) error
}
