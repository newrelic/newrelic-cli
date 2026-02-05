package recipes

import (
	"context"
	"os"
	"sort"
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

	// Skip autodiscovery if env var is set and recipes are targeted
	if os.Getenv("NEW_RELIC_SKIP_AUTODISCOVERY") == "1" &&
		(dt.installerContext.RecipeNamesProvided() || dt.installerContext.RecipePathsProvided()) {
		return isTargeted
	}

	// Check recipe's discoveryMode setting
	if len(recipe.PreInstall.DiscoveryMode) == 1 &&
		(recipe.PreInstall.DiscoveryMode[0] == types.OpenInstallationDiscoveryModeTypes.TARGETED) {
		return isTargeted
	}

	return true
}

func (dt *RecipeDetector) detectRecipe(recipe *types.OpenInstallationRecipe) *RecipeDetectionResult {
	start := time.Now()

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

	return &RecipeDetectionResult{
		recipe,
		status,
		durationMs,
	}
}
