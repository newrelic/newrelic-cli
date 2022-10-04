package execution

import (
	"context"
	"io"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// RecipeExecutor is responsible for execution of the task steps defined in a recipe.
type RecipeExecutor interface {
	Execute(context.Context, types.OpenInstallationRecipe, types.RecipeVars) error
	ExecutePreInstall(context.Context, types.OpenInstallationRecipe, types.RecipeVars) error
	GetOutput() *OutputParser
	GetStdErr() *io.Writer
}
