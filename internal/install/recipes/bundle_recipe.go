package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleRecipe struct {
	Recipe           *types.OpenInstallationRecipe
	Dependencies     []*BundleRecipe
	DetectedStatuses []execution.RecipeStatusType
}

func (br *BundleRecipe) AddStatus(newStatus execution.RecipeStatusType) {
	if br.HasStatus(newStatus) {
		return
	}
	if newStatus == execution.RecipeStatusTypes.AVAILABLE {
		br.DetectedStatuses = append(br.DetectedStatuses, execution.RecipeStatusTypes.DETECTED)
	}
	br.DetectedStatuses = append(br.DetectedStatuses, newStatus)
}

func (br *BundleRecipe) HasStatus(status execution.RecipeStatusType) bool {
	for _, detectedStatus := range br.DetectedStatuses {
		if detectedStatus == status {
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
