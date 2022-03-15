package recipes

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type DetectionStatusProvider interface {
	DetectionStatus(context.Context, *types.OpenInstallationRecipe) execution.RecipeStatusType
}

type RecipeDetector struct {
	processEvaluator DetectionStatusProvider
	scriptEvaluator  DetectionStatusProvider
	recipeEvaluated  map[string][]*DetectedStatusType // same recipe(ref) should only be evaluated one time
}

func NewRecipeDetector() *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: NewProcessEvaluator(),
		scriptEvaluator:  NewScriptEvaluator(),
		recipeEvaluated:  make(map[string][]*DetectedStatusType),
	}
}

func (dt *RecipeDetector) detectBundleRecipe(ctx context.Context, bundleRecipe *BundleRecipe) {

	// if already evaluated
	if s, exists := dt.recipeEvaluated[bundleRecipe.Recipe.Name]; exists {
		bundleRecipe.DetectedStatuses = s
		return
	}
	for i := 0; i < len(bundleRecipe.Dependencies); i++ {
		dependencyBundleRecipe := bundleRecipe.Dependencies[i]
		dt.detectBundleRecipe(ctx, dependencyBundleRecipe)
	}

	if bundleRecipe.AreAllDependenciesAvailable() {
		status, durationMs := dt.detectRecipe(ctx, bundleRecipe.Recipe)
		bundleRecipe.AddDetectionStatus(status, durationMs)
	} else {
		log.Debugf("Dependency not available for recipe %v", bundleRecipe.Recipe.Name)
	}
	dt.recipeEvaluated[bundleRecipe.Recipe.Name] = bundleRecipe.DetectedStatuses
}

func (dt *RecipeDetector) detectRecipe(ctx context.Context, recipe *types.OpenInstallationRecipe) (execution.RecipeStatusType, int64) {
	start := time.Now()
	status := dt.processEvaluator.DetectionStatus(ctx, recipe)
	durationMs := time.Since(start).Milliseconds()

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(ctx, recipe)
		durationMs = time.Since(start).Milliseconds()
		log.Debugf("ScriptEvaluation for recipe:%s completed in %dms with status:%s", recipe.Name, durationMs, status)
	}
	return status, durationMs
}
