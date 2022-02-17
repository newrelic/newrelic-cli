//go:build unit
// +build unit

package recipes

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

var (
	recipeExecutor *execution.MockRecipeExecutor = execution.NewMockRecipeExecutor()
)

func TestScriptEvaluatorShouldNotDetect(t *testing.T) {
	GivenExecutorError("something went wrong")
	recipe := createRecipe("id1", "myrecipe")

	status := GiveScriptEvaluator().DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func TestScriptEvaluatorShouldDetect(t *testing.T) {
	GivenExecutorError("This is the specific message with exit status 132 special case")
	recipe := createRecipe("id1", "myrecipe")

	status := GiveScriptEvaluator().DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.DETECTED, status)
}

func TestScriptEvaluatorShouldGetAvailable(t *testing.T) {
	GivenExecutorSuccess()
	recipe := createRecipe("id1", "myrecipe")

	status := GiveScriptEvaluator().DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func GiveScriptEvaluator() *ScriptEvaluator {
	return newScriptEvaluator(recipeExecutor)
}

func GivenExecutorSuccess() {
	recipeExecutor.ExecuteErr = nil
}

func GivenExecutorError(detail string) {
	recipeExecutor.ExecuteErr = fmt.Errorf(detail)
}
