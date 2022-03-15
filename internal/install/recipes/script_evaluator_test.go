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

func TestScriptEvaluatorShouldGetAvailable(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	status := GivenScriptEvaluator().DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestScriptEvaluatorShouldDetect(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	evaluator := GivenScriptEvaluatorError("This is the specific message with exit status 132 special case")
	status := evaluator.DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.DETECTED, status)
}

func TestScriptEvaluatorShouldNotDetect(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	evaluator := GivenScriptEvaluatorError("something went wrong")
	status := evaluator.DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func GivenScriptEvaluator() *ScriptEvaluator {
	recipeExecutor := execution.NewMockRecipeExecutor()
	return newScriptEvaluator(recipeExecutor)
}

func GivenScriptEvaluatorError(detail string) *ScriptEvaluator {
	recipeExecutor := execution.NewMockRecipeExecutor()
	recipeExecutor.ExecuteErr = fmt.Errorf(detail)
	return newScriptEvaluator(recipeExecutor)
}

func TestScriptEvaluatorShouldBeUnSupported(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	evaluator := GivenScriptEvaluatorError("This is the specific message with exit status 131 un-support case")
	status := evaluator.DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.UNSUPPORTED, status)
}
