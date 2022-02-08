//go:build unit
// +build unit

package recipes

import (
	"context"
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

var (
	executor  *execution.MockRecipeExecutor = execution.NewMockRecipeExecutor()
	evaluator                               = newScriptEvaluator(executor)
	recipe    *types.OpenInstallationRecipe = createRecipe("id1", "myrecipe")
	ctx       context.Context               = nil
)

func TestShouldNotDetect(t *testing.T) {
	GivenExecutorError("something went wrong")

	status := evaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func TestShouldDetect(t *testing.T) {
	GivenExecutorError("This is the specific message with exit status 132 special case")

	status := evaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.DETECTED, status)
}

func TestShouldGetAvailable(t *testing.T) {
	GivenExecutorSuccess()

	status := evaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func GivenExecutorSuccess() {
	executor.ExecuteErr = nil
}

func GivenExecutorError(detail string) {
	executor.ExecuteErr = fmt.Errorf(detail)
}
