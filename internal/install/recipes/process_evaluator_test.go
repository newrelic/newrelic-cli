//go:build unit
// +build unit

package recipes

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
)

var (
	finder           *MockProcessMatchFinder = NewMockProcessMatchFinder()
	processEvaluator *ProcessEvaluator       = newProcessEvaluator(finder, AnyProcesses)
)

func TestProcessEvaluatorShouldGetAvailable(t *testing.T) {
	GivenExecutorSuccess()
	recipe := createRecipe("id1", "myrecipe")

	status := processEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldGetAvailable_Matching(t *testing.T) {
	GivenExecutorSuccess()
	recipe := createRecipe("id1", "myrecipe")
	recipe.ProcessMatch = append(recipe.ProcessMatch, "abc")
	GivenMatchedProcess()

	status := processEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldNotDetect_NoMatch(t *testing.T) {
	GivenExecutorSuccess()
	recipe := createRecipe("id1", "myrecipe")
	recipe.ProcessMatch = append(recipe.ProcessMatch, "abc")
	GivenNoMatchedProcess()

	status := processEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func AnyProcesses(ctx context.Context) []types.GenericProcess {
	return []types.GenericProcess{}
}

func GivenMatchedProcess() {
	p := &types.MatchedProcess{}
	finder.matchedProcesses = append(finder.matchedProcesses, *p)
}

func GivenNoMatchedProcess() {
	finder.matchedProcesses = []types.MatchedProcess{}
}
