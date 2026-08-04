package uninstall

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// detectionStatusEvaluator matches recipes.ProcessEvaluator's method used here,
// so EvaluatorProcessDetector can be tested without a real process table.
type detectionStatusEvaluator interface {
	DetectionStatus(ctx context.Context, r *types.OpenInstallationRecipe, recipeNames []string) execution.RecipeStatusType
}

// EvaluatorProcessDetector reuses the same live process-matching and discovery-script
// evaluation install uses (recipes.RecipeDetector.detectRecipe), since no local install
// manifest exists to tell us whether a recipe is running.
type EvaluatorProcessDetector struct {
	processEvaluator detectionStatusEvaluator
	scriptEvaluator  detectionStatusEvaluator
}

// NewEvaluatorProcessDetector returns a ProcessDetector backed by the given process and script evaluators.
func NewEvaluatorProcessDetector(processEvaluator, scriptEvaluator detectionStatusEvaluator) *EvaluatorProcessDetector {
	return &EvaluatorProcessDetector{processEvaluator: processEvaluator, scriptEvaluator: scriptEvaluator}
}

// IsRunning reports whether the recipe looks like it's running on this host. A recipe with
// no processMatch has no signal to check, so it's treated as not-detected rather than as
// running (recipes.ProcessEvaluator.DetectionStatus would otherwise default to AVAILABLE).
func (d *EvaluatorProcessDetector) IsRunning(ctx context.Context, r types.OpenInstallationRecipe) bool {
	if len(r.ProcessMatch) == 0 {
		return false
	}

	status := d.processEvaluator.DetectionStatus(ctx, &r, nil)
	if status == execution.RecipeStatusTypes.AVAILABLE && r.PreInstall.RequireAtDiscovery != "" {
		status = d.scriptEvaluator.DetectionStatus(ctx, &r, nil)
	}
	return status == execution.RecipeStatusTypes.AVAILABLE
}
