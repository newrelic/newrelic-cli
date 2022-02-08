package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeDetector struct {
	processEvaluator DetectionStatusProvider
	scriptEvaluator  DetectionStatusProvider
}

func newRecipeDetector(processEvaluator DetectionStatusProvider, scriptEvaluator DetectionStatusProvider) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: processEvaluator,
		scriptEvaluator:  scriptEvaluator,
	}
}

func NewRecipeDetector() *RecipeDetector {
	return newRecipeDetector(NewProcessEvaluator(), NewScriptEvaluator())
}

func (dt *RecipeDetector) DetectRecipes(ctx context.Context, recipes []types.OpenInstallationRecipe) map[*types.OpenInstallationRecipe]execution.RecipeStatusType {

	results := make(map[*types.OpenInstallationRecipe]execution.RecipeStatusType)

	for _, recipe := range recipes {

		status := dt.processEvaluator.DetectionStatus(ctx, &recipe)

		if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
			status = dt.scriptEvaluator.DetectionStatus(ctx, &recipe)
		}

		results[&recipe] = status
	}

	return results
}
