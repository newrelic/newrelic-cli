package recipes

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeDetection struct {
	Recipe     *types.OpenInstallationRecipe
	Status     execution.RecipeStatusType
	DurationMs int64
}

type DetectionStatusProvider interface {
	DetectionStatus(context.Context, *types.OpenInstallationRecipe) execution.RecipeStatusType
}

type RecipeDetector struct {
	processEvaluator DetectionStatusProvider
	scriptEvaluator  DetectionStatusProvider
}

func NewRecipeDetector() *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: NewProcessEvaluator(),
		scriptEvaluator:  NewScriptEvaluator(),
	}
}

func (dt *RecipeDetector) DetectRecipes(ctx context.Context, recipes []*types.OpenInstallationRecipe) map[string]*RecipeDetection {
	ds := make(map[string]*RecipeDetection)

	for _, r := range recipes {
		ds[r.Name] = dt.detectRecipe(ctx, r)
	}

	return ds
}

func (dt *RecipeDetector) detectRecipe(ctx context.Context, recipe *types.OpenInstallationRecipe) *RecipeDetection {
	start := time.Now()
	status := dt.processEvaluator.DetectionStatus(ctx, recipe)
	durationMs := time.Since(start).Milliseconds()

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(ctx, recipe)
		durationMs = time.Since(start).Milliseconds()
		log.Debugf("ScriptEvaluation for recipe:%s completed in %dms with status:%s", recipe.Name, durationMs, status)
	}

	return &RecipeDetection{
		recipe,
		status,
		durationMs,
	}
}
