package uninstall

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type fakeRecipeFinder struct {
	recipe *types.OpenInstallationRecipe
	err    error
}

func (f *fakeRecipeFinder) FindRecipe(ctx context.Context, name string) (*types.OpenInstallationRecipe, error) {
	return f.recipe, f.err
}

type fakeProcessDetector struct {
	running bool
}

func (f *fakeProcessDetector) IsRunning(ctx context.Context, r types.OpenInstallationRecipe) bool {
	return f.running
}

type fakeConfirmer struct {
	answers []bool
	prompts []string
	err     error
}

func (f *fakeConfirmer) Confirm(prompt string) (bool, error) {
	f.prompts = append(f.prompts, prompt)
	if f.err != nil {
		return false, f.err
	}
	answer := f.answers[len(f.prompts)-1]
	return answer, nil
}

type fakeTaskfileExecutor struct {
	calls int
	err   error
}

func (f *fakeTaskfileExecutor) ExecuteTaskfile(ctx context.Context, taskfileYAML string, r types.OpenInstallationRecipe, vars types.RecipeVars) error {
	f.calls++
	return f.err
}

type fakeGenericExecutor struct {
	calls    int
	warnings []error
	fatal    error
}

func (f *fakeGenericExecutor) Uninstall(meta types.OpenInstallationUninstallMeta) ([]error, error) {
	f.calls++
	return f.warnings, f.fatal
}

func newTestUninstaller(recipe *types.OpenInstallationRecipe, findErr error, running bool, confirmAnswers []bool, taskfileErr error, genericWarnings []error, genericFatal error) (*RecipeUninstaller, *fakeConfirmer, *fakeTaskfileExecutor, *fakeGenericExecutor) {
	confirmer := &fakeConfirmer{answers: confirmAnswers}
	taskfileExec := &fakeTaskfileExecutor{err: taskfileErr}
	genericExec := &fakeGenericExecutor{warnings: genericWarnings, fatal: genericFatal}

	u := &RecipeUninstaller{
		Finder:       &fakeRecipeFinder{recipe: recipe, err: findErr},
		Detector:     &fakeProcessDetector{running: running},
		Confirmer:    confirmer,
		TaskfileExec: taskfileExec,
		Generic:      genericExec,
	}

	return u, confirmer, taskfileExec, genericExec
}

func TestRun_RecipeNotFound(t *testing.T) {
	u, _, taskfileExec, genericExec := newTestUninstaller(nil, ErrRecipeNotFound, false, nil, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "unknown-recipe"})

	require.Equal(t, StatusFailed, res.Status)
	require.True(t, errors.Is(res.Err, ErrRecipeNotFound))
	require.Equal(t, 0, taskfileExec.calls)
	require.Equal(t, 0, genericExec.calls)
}

func TestRun_RecipeUnsupportedOnHost(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "windows-only-recipe"}
	u, _, _, _ := newTestUninstaller(recipe, ErrRecipeUnsupportedOnHost, false, nil, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "windows-only-recipe"})

	require.Equal(t, StatusUnsupported, res.Status)
}

func TestRun_AssumeYesSkipsAllConfirmations(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec, _ := newTestUninstaller(recipe, nil, false /* not detected */, nil, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Empty(t, confirmer.prompts, "assumeYes should skip every confirmation")
	require.Equal(t, 1, taskfileExec.calls)
}

func TestRun_DetectedAndConfirmed_RunsTaskfileTier(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec, _ := newTestUninstaller(recipe, nil, true, []bool{true}, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Len(t, confirmer.prompts, 1, "detected host should only need the general confirmation")
	require.Equal(t, 1, taskfileExec.calls)
}

func TestRun_DetectedAndDeclined_Aborts(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, _, taskfileExec, _ := newTestUninstaller(recipe, nil, true, []bool{false}, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusAborted, res.Status)
	require.Equal(t, 0, taskfileExec.calls)
}

func TestRun_NotDetectedAndDeclined_AbortsWithoutSecondPrompt(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec, _ := newTestUninstaller(recipe, nil, false, []bool{false}, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe"})

	require.Equal(t, StatusAborted, res.Status)
	require.Len(t, confirmer.prompts, 1, "declining the not-detected prompt should not trigger the general prompt too")
	require.Equal(t, 0, taskfileExec.calls)
}

func TestRun_NotDetectedButForced_SkipsNotDetectedPromptButStillAsksGeneralConfirm(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	u, confirmer, taskfileExec, _ := newTestUninstaller(recipe, nil, false, []bool{true}, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", Force: true})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Len(t, confirmer.prompts, 1)
	require.Equal(t, 1, taskfileExec.calls)
}

func TestRun_TaskfileExecutionFails(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe", Uninstall: "default:\n  cmds:\n    - true\n"}
	execErr := errors.New("boom")
	u, _, _, _ := newTestUninstaller(recipe, nil, true, nil, execErr, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusFailed, res.Status)
	require.ErrorIs(t, res.Err, execErr)
}

func TestRun_NoAutomatedUninstallWhenNeitherTierDefined(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{Name: "test-recipe"}
	u, _, taskfileExec, genericExec := newTestUninstaller(recipe, nil, true, nil, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusNoAutomatedUninstall, res.Status)
	require.Equal(t, 0, taskfileExec.calls)
	require.Equal(t, 0, genericExec.calls)
}

func TestRun_GenericTierCleanRun(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{
		Name:          "test-recipe",
		UninstallMeta: types.OpenInstallationUninstallMeta{Packages: []string{"newrelic-infra"}},
	}
	u, _, _, genericExec := newTestUninstaller(recipe, nil, true, nil, nil, nil, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusUninstalled, res.Status)
	require.Equal(t, 1, genericExec.calls)
}

func TestRun_GenericTierWithWarningsIsReportedAsFailed(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{
		Name:          "test-recipe",
		UninstallMeta: types.OpenInstallationUninstallMeta{Packages: []string{"newrelic-infra"}},
	}
	warnings := []error{errors.New("package already removed")}
	u, _, _, _ := newTestUninstaller(recipe, nil, true, nil, nil, warnings, nil)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusFailed, res.Status)
	require.Equal(t, warnings, res.Warnings)
}

func TestRun_GenericTierFatalErrorIsFailed(t *testing.T) {
	recipe := &types.OpenInstallationRecipe{
		Name:          "test-recipe",
		UninstallMeta: types.OpenInstallationUninstallMeta{Packages: []string{"newrelic-infra"}},
	}
	fatal := ErrPermissionDenied
	u, _, _, _ := newTestUninstaller(recipe, nil, true, nil, nil, nil, fatal)

	res := u.Run(context.Background(), Options{RecipeName: "test-recipe", AssumeYes: true})

	require.Equal(t, StatusFailed, res.Status)
	require.ErrorIs(t, res.Err, ErrPermissionDenied)
}
