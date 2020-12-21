package execution

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeExecutor is responsible for execution of the task steps defined in a
// recipe.
type RecipeExecutor interface {
	Execute(context.Context, types.DiscoveryManifest, types.Recipe) error
}
