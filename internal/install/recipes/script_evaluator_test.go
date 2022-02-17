//go:build unit
// +build unit

package recipes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

func TestScriptEvaluatorShouldNotDetect(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	evaluator := GiveScriptEvaluatorError("something went wrong")
	status := evaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func TestScriptEvaluatorShouldDetect(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	evaluator := GiveScriptEvaluatorError("This is the specific message with exit status 132 special case")
	status := evaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.DETECTED, status)
}

func TestScriptEvaluatorShouldGetAvailable(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	status := GiveScriptEvaluatorSuccess().DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func GiveScriptEvaluatorSuccess() *ScriptEvaluator {
	recipeExecutor := execution.NewMockRecipeExecutor()
	return newScriptEvaluator(recipeExecutor)
}

func GiveScriptEvaluatorError(detail string) *ScriptEvaluator {
	recipeExecutor := execution.NewMockRecipeExecutor()
	recipeExecutor.ExecuteErr = fmt.Errorf(detail)
	return newScriptEvaluator(recipeExecutor)
}
