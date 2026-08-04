//go:build integration

package uninstall

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// newFixtureUninstaller wires a RecipeUninstaller against the real local-file
// recipe fetcher, host filter, and taskfile executor, using the YAML
// fixtures under test/recipes/. Detection and confirmation are faked so the
// test doesn't depend on the current host's process table or a terminal.
func newFixtureUninstaller(running bool) *RecipeUninstaller {
	fetcher := &recipes.LocalRecipeFetcher{Path: "../../test/recipes"}
	loaderFunc := func() ([]*types.OpenInstallationRecipe, error) {
		return fetcher.FetchRecipes(context.Background())
	}
	manifest := &types.DiscoveryManifest{OS: "linux"}

	return &RecipeUninstaller{
		Finder:       NewRepositoryRecipeFinder(loaderFunc, manifest),
		Detector:     &fakeProcessDetector{running: running},
		Confirmer:    &fakeConfirmer{answers: []bool{true, true}},
		TaskfileExec: execution.NewGoTaskRecipeExecutor(),
	}
}

func TestIntegration_TaskfileFixtureUninstallsCleanly(t *testing.T) {
	u := newFixtureUninstaller(true)

	res := u.Run(context.Background(), Options{RecipeName: "test-uninstall-taskfile", AssumeYes: true})

	require.Equal(t, StatusUninstalled, res.Status)
}

func TestIntegration_UnknownFixtureRecipeIsNotFound(t *testing.T) {
	u := newFixtureUninstaller(true)

	res := u.Run(context.Background(), Options{RecipeName: "does-not-exist-in-fixtures", AssumeYes: true})

	require.Equal(t, StatusFailed, res.Status)
	require.ErrorIs(t, res.Err, ErrRecipeNotFound)
}
