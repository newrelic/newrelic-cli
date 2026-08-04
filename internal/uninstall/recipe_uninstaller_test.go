//go:build unit

package uninstall

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRun_NotDetectedConfirmationErrors(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u := &RecipeUninstaller{
		Finder:       &fakeRecipeFinder{recipe: recipe},
		Detector:     &fakeProcessDetector{running: false},
		Confirmer:    &fakeConfirmer{err: errors.New("boom")},
		TaskfileExec: &fakeTaskfileExecutor{},
	}

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusFailed, res.Status)
}

func TestRun_GeneralConfirmationErrors(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u := &RecipeUninstaller{
		Finder:       &fakeRecipeFinder{recipe: recipe},
		Detector:     &fakeProcessDetector{running: true},
		Confirmer:    &fakeConfirmer{err: errors.New("boom")},
		TaskfileExec: &fakeTaskfileExecutor{},
	}

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusFailed, res.Status)
}

func TestRun_RecipeNotFound(t *testing.T) {
	u, _, taskfileExec := newTestUninstaller(nil, ErrRecipeNotFound, false, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "unknown-recipe"})

	require.Equal(t, StatusFailed, res.Status)
	require.True(t, errors.Is(res.Err, ErrRecipeNotFound))
	require.Equal(t, 0, taskfileExec.calls)
}

func TestRun_RecipeUnsupportedOnHost(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "windows-only-recipe"}
	u, _, _ := newTestUninstaller(recipe, ErrRecipeUnsupportedOnHost, false, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "windows-only-recipe"})

	require.Equal(t, StatusUnsupported, res.Status)
}

func TestRun_AssumeYesSkipsAllConfirmations(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec := newTestUninstaller(recipe, nil, false /* not detected */, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Empty(t, confirmer.prompts, "assumeYes should skip every confirmation")
	require.Equal(t, 1, taskfileExec.calls)
}

func TestRun_DetectedAndConfirmed_RunsTaskfile(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec := newTestUninstaller(recipe, nil, true, []bool{true}, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Len(t, confirmer.prompts, 1, "detected host should only need the general confirmation")
	require.Equal(t, 1, taskfileExec.calls)
}

func TestRun_DetectedAndDeclined_Aborts(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, _, taskfileExec := newTestUninstaller(recipe, nil, true, []bool{false}, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusAborted, res.Status)
	require.Equal(t, 0, taskfileExec.calls)
}

func TestRun_NotDetectedAndDeclined_AbortsWithoutSecondPrompt(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec := newTestUninstaller(recipe, nil, false, []bool{false}, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusAborted, res.Status)
	require.Len(t, confirmer.prompts, 1, "declining the not-detected prompt should not trigger the general prompt too")
	require.Equal(t, 0, taskfileExec.calls)
}

func TestRun_NotDetectedButForced_SkipsNotDetectedPromptButStillAsksGeneralConfirm(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec := newTestUninstaller(recipe, nil, false, []bool{true}, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", Force: true})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Len(t, confirmer.prompts, 1)
	require.Equal(t, 1, taskfileExec.calls)
}

func TestRun_TaskfileExecutionFails(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	execErr := errors.New("boom")
	u, _, _ := newTestUninstaller(recipe, nil, true, nil, execErr)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusFailed, res.Status)
	require.ErrorIs(t, res.Err, execErr)
}

func TestRun_NoAutomatedUninstallWhenNoTaskfileDefined(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe"}
	u, _, taskfileExec := newTestUninstaller(recipe, nil, true, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusNoAutomatedUninstall, res.Status)
	require.Equal(t, 0, taskfileExec.calls)
}

func TestRun_NoAutomatedUninstallSkipsConfirmationsEntirely(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe"}
	u, confirmer, _ := newTestUninstaller(recipe, nil, true, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusNoAutomatedUninstall, res.Status)
	require.Empty(t, confirmer.prompts, "there's nothing to confirm when there's no automated uninstall for this recipe")
}

func TestRun_AssumeYesIsThreadedIntoRecipeVars(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, _, taskfileExec := newTestUninstaller(recipe, nil, true, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Equal(t, "true", taskfileExec.lastVars["assumeYes"])
}
