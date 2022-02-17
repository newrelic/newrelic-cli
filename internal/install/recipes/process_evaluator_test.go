//go:build unit
// +build unit

package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestProcessEvaluatorShouldGetAvailable(t *testing.T) {
	recipe := NewRecipeBuilder().Build()

	status := GivenProcessEvaluator().DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldGetAvailable_Matching(t *testing.T) {
	recipe := NewRecipeBuilder().ProcessMatch("abc").Build()
	processEvaluator := GivenProcessEvaluatorMatchedProcess()

	status := processEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldNotDetect_NoMatch(t *testing.T) {
	recipe := NewRecipeBuilder().ProcessMatch("abc").Build()
	processEvaluator := GivenProcessEvaluator()

	status := processEvaluator.DetectionStatus(ctx, recipe)

	require.Equal(t, execution.RecipeStatusTypes.NULL, status)
}

func AnyProcesses(ctx context.Context) []types.GenericProcess {
	return []types.GenericProcess{}
}

func GivenProcessEvaluatorMatchedProcess() *ProcessEvaluator {
	finder := NewMockProcessMatchFinder()
	p := &types.MatchedProcess{}
	finder.matchedProcesses = append(finder.matchedProcesses, *p)
	processEvaluator := newProcessEvaluator(finder, AnyProcesses, false)
	return processEvaluator
}

func GivenProcessEvaluator() *ProcessEvaluator {
	finder := NewMockProcessMatchFinder()
	finder.matchedProcesses = []types.MatchedProcess{}
	processEvaluator := newProcessEvaluator(finder, AnyProcesses, false)
	return processEvaluator
}
