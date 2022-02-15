package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleRecipe struct {
	Recipe       *types.OpenInstallationRecipe
	Dependencies []*BundleRecipe
	//maybe timestamp
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
