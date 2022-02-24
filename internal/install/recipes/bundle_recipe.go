package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleRecipe struct {
	Recipe           *types.OpenInstallationRecipe
	Dependencies     []*BundleRecipe
	DetectedStatuses []*DetectedStatusType
}

type DetectedStatusType struct {
	Status     execution.RecipeStatusType
	DurationMs int64
}

func (br *BundleRecipe) AddDetectionStatus(newStatus execution.RecipeStatusType, durationMs int64) {
	if br.HasStatus(newStatus) {
		return
	}
	if newStatus == execution.RecipeStatusTypes.AVAILABLE {
		ds := &DetectedStatusType{
			Status:     execution.RecipeStatusTypes.DETECTED,
			DurationMs: durationMs,
		}
		br.DetectedStatuses = append(br.DetectedStatuses, ds)
	}
	br.DetectedStatuses = append(br.DetectedStatuses, &DetectedStatusType{Status: newStatus, DurationMs: durationMs})
}

func (br *BundleRecipe) HasStatus(status execution.RecipeStatusType) bool {
	for _, detectedStatus := range br.DetectedStatuses {
		if detectedStatus.Status == status {
			return true
		}
	}
	return false
}

//TODO: might not need!
func (br *BundleRecipe) Flatten() map[string]bool {

	results := make(map[string]bool)
	br.flatten(results)

	return results
}

func (br *BundleRecipe) flatten(recipeMap map[string]bool) {

	if _, ok := recipeMap[br.Recipe.Name]; !ok {
		recipeMap[br.Recipe.Name] = true
	}

	for _, d := range br.Dependencies {
		d.flatten(recipeMap)
	}
}
