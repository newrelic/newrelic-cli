//go:build unit || integration

package uninstall

import (
	"context"

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
	calls    int
	err      error
	lastVars types.RecipeVars
}

func (f *fakeTaskfileExecutor) ExecuteTaskfile(ctx context.Context, taskfileYAML string, r types.OpenInstallationRecipe, vars types.RecipeVars) error {
	f.calls++
	f.lastVars = vars
	return f.err
}

func newTestUninstaller(recipe *types.OpenInstallationRecipe, findErr error, running bool, confirmAnswers []bool, taskfileErr error) (*RecipeUninstaller, *fakeConfirmer, *fakeTaskfileExecutor) {
	confirmer := &fakeConfirmer{answers: confirmAnswers}
	taskfileExec := &fakeTaskfileExecutor{err: taskfileErr}

	u := &RecipeUninstaller{
		Finder:       &fakeRecipeFinder{recipe: recipe, err: findErr},
		Detector:     &fakeProcessDetector{running: running},
		Confirmer:    confirmer,
		TaskfileExec: taskfileExec,
	}

	return u, confirmer, taskfileExec
}
