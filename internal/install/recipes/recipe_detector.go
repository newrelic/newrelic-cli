package recipes

import (
	"context"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type DetectionStatusProvider interface {
	DetectionStatus(context.Context, *types.OpenInstallationRecipe, []string) execution.RecipeStatusType
}

type RecipeDetectionResult struct {
	Recipe     *types.OpenInstallationRecipe
	Status     execution.RecipeStatusType
	DurationMs int64
}

type RecipeDetectionResults []*RecipeDetectionResult

func (rd *RecipeDetectionResults) GetRecipeDetection(name string) (*RecipeDetectionResult, bool) {
	for _, d := range *rd {
		if d.Recipe.Name == name {
			return d, true
		}
	}
	return nil, false
}

func (rd RecipeDetectionResults) Len() int {
	return len(rd)
}

func (rd RecipeDetectionResults) Swap(i, j int) {
	rd[i], rd[j] = rd[j], rd[i]
}

func (rd RecipeDetectionResults) Less(i, j int) bool {
	return rd[i].Recipe.Name < rd[j].Recipe.Name
}

type RecipeDetector struct {
	processEvaluator DetectionStatusProvider
	scriptEvaluator  DetectionStatusProvider
	context          context.Context
	repo             Finder
	installerContext *types.InstallerContext
}

func NewRecipeDetector(contex context.Context, repo *RecipeRepository, peval ProcessEvaluatorInterface, ic *types.InstallerContext) *RecipeDetector {
	return &RecipeDetector{
		processEvaluator: peval,
		scriptEvaluator:  NewScriptEvaluator(),
		context:          contex,
		repo:             repo,
		installerContext: ic,
	}
}

func (dt *RecipeDetector) GetDetectedRecipes() (RecipeDetectionResults, RecipeDetectionResults, error) {
	availableRecipes := RecipeDetectionResults{}
	unavailableRecipes := RecipeDetectionResults{}
	recipes, err := dt.repo.FindAll()
	if err != nil {
		return nil, nil, err
	}
	for _, r := range recipes {
		dr := dt.detectRecipe(r)
		if r.Name == "logs-integration" {
			log.Debug("logs-integration output: ", strings.Contains(r.Name, "logs-integration"))
			log.Debug("Status of It:", r.InstallTargets)
		}

		if dr.Status == execution.RecipeStatusTypes.AVAILABLE {
			availableRecipes = append(availableRecipes, dr)
		} else {
			unavailableRecipes = append(unavailableRecipes, dr)
		}
	}
	sort.Sort(availableRecipes)

	return availableRecipes, unavailableRecipes, nil
}

func (dt *RecipeDetector) shouldDiscover(recipe *types.OpenInstallationRecipe) bool {
	isTargeted := dt.installerContext.IsRecipeTargeted(recipe.Name)
	if len(recipe.PreInstall.DiscoveryMode) == 1 &&
		(recipe.PreInstall.DiscoveryMode[0] == types.OpenInstallationDiscoveryModeTypes.TARGETED) {
		return isTargeted
	}

	return true
}

func (dt *RecipeDetector) detectRecipe(recipe *types.OpenInstallationRecipe) *RecipeDetectionResult {
	start := time.Now()

	if recipe.Name == "logs-integration" {
		log.Debug("log integration installation status:")
	}

	if !dt.shouldDiscover(recipe) {
		durationMs := time.Since(start).Milliseconds()
		return &RecipeDetectionResult{
			recipe,
			execution.RecipeStatusTypes.NULL,
			durationMs,
		}
	}

	recipeNames := dt.installerContext.RecipeNames

	status := dt.processEvaluator.DetectionStatus(dt.context, recipe, recipeNames)
	durationMs := time.Since(start).Milliseconds()

	if status == execution.RecipeStatusTypes.AVAILABLE && recipe.PreInstall.RequireAtDiscovery != "" {
		status = dt.scriptEvaluator.DetectionStatus(dt.context, recipe, recipeNames)
		durationMs = time.Since(start).Milliseconds()
		log.Debugf("ScriptEvaluation for recipe:%s completed in %dms with status:%s", recipe.Name, durationMs, status)
	}
	if recipe.Name == "logs-integration" {
		log.Debug("+++++++++++")
		log.Debugf(string(status))
		log.Debug("+++++++++++")
	}

	return &RecipeDetectionResult{
		recipe,
		status,
		durationMs,
	}
}
