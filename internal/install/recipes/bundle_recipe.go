package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleRecipe struct {
	Recipe       *types.OpenInstallationRecipe
	Dependencies []*BundleRecipe
	//keep track of reported status
	//optional: datetime instead of if it's saved
	//can have method all the status is reported/saved
	Statuses []execution.RecipeStatusType
}

func (br *BundleRecipe) AddStatus(status execution.RecipeStatusType) {
	if br.HasStatus(status) {
		return
	}
	if status == execution.RecipeStatusTypes.AVAILABLE {
		br.Statuses = append(br.Statuses, execution.RecipeStatusTypes.DETECTED)
	}
	br.Statuses = append(br.Statuses, status)
}

func (br *BundleRecipe) HasStatus(status execution.RecipeStatusType) bool {
	for _, value := range br.Statuses {
		if value == status {
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
