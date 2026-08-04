//go:build unit

package uninstall

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestUninstallCommand(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "uninstall", Command.Name())

	testcobra.CheckCobraMetadata(t, Command)
	testcobra.CheckCobraRequiredFlags(t, Command, []string{"recipe"})
}

func TestEnsureSingleConcurrentUninstall_AllowsWhenOnlyItself(t *testing.T) {
	evaluator := recipes.NewMockProcessEvaluator()
	evaluator.WithProcesses([]types.GenericProcess{
		recipes.NewMockProcess("newrelic uninstall -n infra", "newrelic", 1),
	})

	err := ensureSingleConcurrentUninstall(context.Background(), evaluator)

	require.NoError(t, err)
}

func TestEnsureSingleConcurrentUninstall_FailsWhenAnotherIsRunning(t *testing.T) {
	evaluator := recipes.NewMockProcessEvaluator()
	evaluator.WithProcesses([]types.GenericProcess{
		recipes.NewMockProcess("newrelic uninstall -n infra", "newrelic", 1),
		recipes.NewMockProcess("newrelic install -n logging", "newrelic", 2),
	})

	err := ensureSingleConcurrentUninstall(context.Background(), evaluator)

	require.Error(t, err)
}

func TestRunE_ErrorsWhenAnotherInstallOrUninstallIsRunning(t *testing.T) {
	originalEvaluatorFactory := newProcessEvaluator
	originalRecipeName := recipeName
	defer func() {
		newProcessEvaluator = originalEvaluatorFactory
		recipeName = originalRecipeName
	}()

	evaluator := recipes.NewMockProcessEvaluator()
	evaluator.WithProcesses([]types.GenericProcess{
		recipes.NewMockProcess("newrelic uninstall -n infra", "newrelic", 1),
		recipes.NewMockProcess("newrelic install -n logging", "newrelic", 2),
	})
	newProcessEvaluator = func() recipes.ProcessEvaluatorInterface { return evaluator }
	recipeName = "infrastructure-agent-installer"

	err := Command.RunE(Command, []string{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "only 1 newrelic install/uninstall command can run at one time")
}

func TestReportResult(t *testing.T) {
	original := recipeName
	defer func() { recipeName = original }()
	recipeName = "test-recipe"

	require.NoError(t, reportResult(Result{Status: StatusUninstalled}))
	require.NoError(t, reportResult(Result{Status: StatusNoAutomatedUninstall}))
	require.NoError(t, reportResult(Result{Status: StatusAborted}))
	require.Error(t, reportResult(Result{Status: StatusUnsupported}))
	require.Error(t, reportResult(Result{Status: StatusFailed, Err: errors.New("boom"), Warnings: []error{errors.New("warn")}}))
}

func TestNewRecipeFetcher(t *testing.T) {
	originalLocal, originalPaths := localRecipes, recipePaths
	defer func() { localRecipes, recipePaths = originalLocal, originalPaths }()

	localRecipes = "/some/path"
	recipePaths = []string{"a.yml"}
	_, ok := newRecipeFetcher().(*recipes.LocalRecipeFetcher)
	require.True(t, ok, "localRecipes should take precedence")

	localRecipes = ""
	_, ok = newRecipeFetcher().(*recipes.RecipeFileFetcher)
	require.True(t, ok, "recipePaths should be used when localRecipes is unset")

	recipePaths = nil
	_, ok = newRecipeFetcher().(*recipes.EmbeddedRecipeFetcher)
	require.True(t, ok, "should fall back to the embedded fetcher")
}
