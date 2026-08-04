//go:build unit

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

	status := GivenProcessEvaluator().DetectionStatus(context.Background(), recipe, []string{})

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldGetAvailable_Matching(t *testing.T) {
	recipe := NewRecipeBuilder().ProcessMatch("abc").Build()
	processEvaluator := GivenProcessEvaluatorMatchedProcess()

	status := processEvaluator.DetectionStatus(context.Background(), recipe, []string{})

	require.Equal(t, execution.RecipeStatusTypes.AVAILABLE, status)
}

func TestProcessEvaluatorShouldNotDetect_NoMatch(t *testing.T) {
	recipe := NewRecipeBuilder().ProcessMatch("abc").Build()
	processEvaluator := GivenProcessEvaluator()

	status := processEvaluator.DetectionStatus(context.Background(), recipe, []string{})

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

func TestCountConcurrentNewRelicSubcommandProcesses(t *testing.T) {
	processes := []types.GenericProcess{
		NewMockProcess("newrelic install -n infra", "newrelic", 1),
		NewMockProcess("newrelic uninstall -n infra", "newrelic", 2),
		NewMockProcess("newrelic diagnose", "newrelic", 3),
		NewMockProcess("otherprocess install", "otherprocess", 4),
	}

	count := CountConcurrentNewRelicSubcommandProcesses(processes)

	require.Equal(t, 2, count, "should match the install and uninstall processes, but not diagnose or unrelated binaries")
}

func TestCountConcurrentNewRelicSubcommandProcesses_WindowsExeSuffix(t *testing.T) {
	processes := []types.GenericProcess{
		NewMockProcess("newrelic.exe uninstall -n infra", "newrelic.exe", 1),
	}

	count := CountConcurrentNewRelicSubcommandProcesses(processes)

	require.Equal(t, 1, count)
}
