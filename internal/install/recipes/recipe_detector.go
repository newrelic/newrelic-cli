package recipes

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeDetectionResult struct {
	Recipe     *types.OpenInstallationRecipe
	Status     execution.RecipeStatusType
	DurationMs int64
}

type DetectionStatusProvider interface {
	DetectionStatus(context.Context, *types.OpenInstallationRecipe) execution.RecipeStatusType
}

type RecipeDetector struct {
	processEvaluator      DetectionStatusProvider
	scriptEvaluator       DetectionStatusProvider
	context               context.Context
	repo                  Finder
	recipeDetectionResult map[string]*RecipeDetectionResult
	availableRecipes      map[string]*RecipeDetectionResult
	unavaliableRecipes    map[string]*RecipeDetectionResult
}

func NewRecipeDetector(contex context.Context, repo *RecipeRepository) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator:      NewProcessEvaluator(),
		scriptEvaluator:       NewScriptEvaluator(),
		context:               contex,
		repo:                  repo,
		recipeDetectionResult: make(map[string]*RecipeDetectionResult),
		availableRecipes:      make(map[string]*RecipeDetectionResult),
		unavaliableRecipes:    make(map[string]*RecipeDetectionResult),
	}
}

func (dt *RecipeDetector) GetDetectedRecipes() (map[string]*RecipeDetectionResult,
	map[string]*RecipeDetectionResult,
	error) {

	if len(dt.availableRecipes) != 0 || len(dt.unavaliableRecipes) != 0 {
		return dt.availableRecipes, dt.unavaliableRecipes, nil
	}

	recipes, err := dt.repo.FindAll()
	if err != nil {
		return nil, nil, err
	}

	for _, r := range recipes {
		dr := dt.detectRecipe(r)
		dt.recipeDetectionResult[dr.Recipe.Name] = dr

		if dr.Status == execution.RecipeStatusTypes.AVAILABLE {
			dt.availableRecipes[dr.Recipe.Name] = dr
		} else {
			dt.unavaliableRecipes[dr.Recipe.Name] = dr
		}
	}
	return dt.availableRecipes, dt.unavaliableRecipes, nil
}

func (dt *RecipeDetector) detectRecipe(recipe *types.OpenInstallationRecipe) *RecipeDetectionResult {
	start := time.Now()
	status := dt.processEvaluator.DetectionStatus(dt.context, recipe)
	durationMs := time.Since(start).Milliseconds()

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(dt.context, recipe)
		durationMs = time.Since(start).Milliseconds()
		log.Debugf("ScriptEvaluation for recipe:%s completed in %dms with status:%s", recipe.Name, durationMs, status)
	}

	return &RecipeDetectionResult{
		recipe,
		status,
		durationMs,
	}
}
