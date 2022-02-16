package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"time"
)

type BundleRecipe struct {
	Recipe         *types.OpenInstallationRecipe
	Dependencies   []*BundleRecipe
	RecipeStatuses []RecipeStatus
}

type RecipeStatus struct {
	Status     execution.RecipeStatusType
	StatusTime time.Time
}

func (br *BundleRecipe) AddStatus(newStatus execution.RecipeStatusType, statusTime time.Time) {
	if br.HasStatus(newStatus) {
		return
	}
	if newStatus == execution.RecipeStatusTypes.AVAILABLE {
		br.RecipeStatuses = append(br.RecipeStatuses, RecipeStatus{Status: execution.RecipeStatusTypes.DETECTED, StatusTime: statusTime})
	}
	br.RecipeStatuses = append(br.RecipeStatuses, RecipeStatus{Status: newStatus, StatusTime: statusTime})
}

func (br *BundleRecipe) HasStatus(status execution.RecipeStatusType) bool {
	for _, value := range br.RecipeStatuses {
		if value.Status == status {
			return true
		}
	}
	return false
}

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
