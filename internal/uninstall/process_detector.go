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

// EvaluatorProcessDetector reuses the same live process-matching install uses
// to decide whether a recipe is running, since no local install manifest exists.
type EvaluatorProcessDetector struct {
	evaluator detectionStatusEvaluator
}

// NewEvaluatorProcessDetector returns a ProcessDetector backed by the given evaluator.
func NewEvaluatorProcessDetector(evaluator detectionStatusEvaluator) *EvaluatorProcessDetector {
	return &EvaluatorProcessDetector{evaluator: evaluator}
}

func (d *EvaluatorProcessDetector) IsRunning(ctx context.Context, r types.OpenInstallationRecipe) bool {
	return d.evaluator.DetectionStatus(ctx, &r, nil) == execution.RecipeStatusTypes.AVAILABLE
}
