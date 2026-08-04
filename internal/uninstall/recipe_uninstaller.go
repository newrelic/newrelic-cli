package uninstall

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// Status is the outcome of a single Run call.
type Status string

const (
	// StatusUninstalled means the recipe's automated removal completed cleanly.
	StatusUninstalled Status = "UNINSTALLED"
	// StatusFailed means something prevented a clean, complete removal.
	StatusFailed Status = "FAILED"
	// StatusUnsupported means the recipe exists but does not target this host.
	StatusUnsupported Status = "UNSUPPORTED"
	// StatusNoAutomatedUninstall means the recipe defines no Uninstall taskfile.
	StatusNoAutomatedUninstall Status = "NO_AUTOMATED_UNINSTALL"
	// StatusAborted means the user declined a confirmation prompt.
	StatusAborted Status = "ABORTED"
)

// ErrRecipeNotFound means no recipe with the given name exists.
var ErrRecipeNotFound = errors.New("recipe not found")

// ErrRecipeUnsupportedOnHost means the recipe exists but its InstallTargets don't match this host.
var ErrRecipeUnsupportedOnHost = errors.New("recipe not supported on this host")

// Result is the outcome of a Run call.
type Result struct {
	Status   Status
	Warnings []error
	Err      error
}

// RecipeFinder resolves a recipe name to a host-compatible recipe.
type RecipeFinder interface {
	FindRecipe(ctx context.Context, name string) (*types.OpenInstallationRecipe, error)
}

// ProcessDetector reports whether a recipe's process appears to be running on this host.
type ProcessDetector interface {
	IsRunning(ctx context.Context, r types.OpenInstallationRecipe) bool
}

// Confirmer asks the user a yes/no question.
type Confirmer interface {
	Confirm(prompt string) (bool, error)
}

// TaskfileExecutor runs an arbitrary go-task Taskfile for a recipe.
type TaskfileExecutor interface {
	ExecuteTaskfile(ctx context.Context, taskfileYAML string, r types.OpenInstallationRecipe, vars types.RecipeVars) error
}

// Options configures a single uninstall run.
type Options struct {
	RecipeName string
	AssumeYes  bool
	Force      bool
}

// RecipeUninstaller orchestrates removing a single named recipe.
type RecipeUninstaller struct {
	Finder       RecipeFinder
	Detector     ProcessDetector
	Confirmer    Confirmer
	TaskfileExec TaskfileExecutor
}

// Run resolves, confirms, and removes the named recipe.
func (u *RecipeUninstaller) Run(ctx context.Context, opts Options) Result {
	recipe, err := u.Finder.FindRecipe(ctx, opts.RecipeName)
	if err != nil {
		if errors.Is(err, ErrRecipeUnsupportedOnHost) {
			return Result{Status: StatusUnsupported, Err: err}
		}
		return Result{Status: StatusFailed, Err: err}
	}

	if recipe.Uninstall == "" {
		return Result{Status: StatusNoAutomatedUninstall}
	}

	detected := u.Detector.IsRunning(ctx, *recipe)

	if !detected && !opts.Force && !opts.AssumeYes {
		ok, err := u.Confirmer.Confirm(fmt.Sprintf(
			"%s does not appear to be running on this host. Uninstall it anyway?", recipe.Name))
		if err != nil {
			return Result{Status: StatusFailed, Err: err}
		}
		if !ok {
			return Result{Status: StatusAborted}
		}
	}

	if !opts.AssumeYes {
		ok, err := u.Confirmer.Confirm(fmt.Sprintf(
			"This will remove %s from this host. Continue?", recipe.Name))
		if err != nil {
			return Result{Status: StatusFailed, Err: err}
		}
		if !ok {
			return Result{Status: StatusAborted}
		}
	}

	vars := types.RecipeVars{"assumeYes": strconv.FormatBool(opts.AssumeYes)}
	if err := u.TaskfileExec.ExecuteTaskfile(ctx, recipe.Uninstall, *recipe, vars); err != nil {
		return Result{Status: StatusFailed, Err: err}
	}
	return Result{Status: StatusUninstalled}
}
