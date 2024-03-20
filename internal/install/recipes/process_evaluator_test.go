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

	status := GivenProcessEvaluator().DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldGetAvailable_Matching(t *testing.T) {
	recipe := NewRecipeBuilder().ProcessMatch("abc").Build()
	processEvaluator := GivenProcessEvaluatorMatchedProcess()

	status := processEvaluator.DetectionStatus(context.Background(), recipe)

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldNotDetect_NoMatch(t *testing.T) {
	recipe := NewRecipeBuilder().ProcessMatch("abc").Build()
	processEvaluator := GivenProcessEvaluator()

	status := processEvaluator.DetectionStatus(context.Background(), recipe)

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

func TestProcessEvaluatorShouldFailFindingNonExistingProcess(t *testing.T) {
	pe := NewProcessEvaluator()
	p := NewMockProcess("/bin/process-a", "process-a", 1234)
	pe.cachedProcess = append(pe.cachedProcess, p)

	found := pe.FindProcess("process-b")
	require.Equal(t, false, found)
}

func TestProcessEvaluatorShouldSucceedFindingExistingProcess(t *testing.T) {
	pe := NewProcessEvaluator()
	p := NewMockProcess("/bin/process-a", "process-a", 1234)
	pe.cachedProcess = append(pe.cachedProcess, p)

	found := pe.FindProcess("process-a")
	require.Equal(t, true, found)
}
