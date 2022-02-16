//go:build unit
// +build unit

package recipes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

var (
	recipeExecutor  *execution.MockRecipeExecutor = execution.NewMockRecipeExecutor()
	scriptEvaluator                               = newScriptEvaluator(recipeExecutor)
)

func TestScriptEvaluatorShouldNotDetect(t *testing.T) {
	GivenExecutorError("something went wrong")
	recipe := createRecipe("id1", "myrecipe")

	status := scriptEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func TestScriptEvaluatorShouldDetect(t *testing.T) {
	GivenExecutorError("This is the specific message with exit Status 132 special case")
	recipe := createRecipe("id1", "myrecipe")

	status := scriptEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.DETECTED, status)
}

func TestScriptEvaluatorShouldGetAvailable(t *testing.T) {
	GivenExecutorSuccess()
	recipe := createRecipe("id1", "myrecipe")

	status := scriptEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func GivenExecutorSuccess() {
	recipeExecutor.ExecuteErr = nil
}

func GivenExecutorError(detail string) {
	recipeExecutor.ExecuteErr = fmt.Errorf(detail)
}
