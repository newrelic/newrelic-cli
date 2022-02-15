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
type Bundle struct {
	BundleRecipes []*BundleRecipe
}

func (b *Bundle) AddRecipe(bundleRecipe *BundleRecipe) {

	if index := b.IndexOf(bundleRecipe.Recipe.Name); index != -1 {
		b.BundleRecipes[index] = bundleRecipe
	} else {
		b.BundleRecipes = append(b.BundleRecipes, bundleRecipe)
	}
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

//TODO: do we need an IndexOf?
func (b *Bundle) IndexOf(name string) int {
	for i := range b.BundleRecipes {
		if b.BundleRecipes[i].Recipe.Name == name {
			return i
		}
	}

	return -1
}

//TODO: Not sure if ths is needed
func (b *Bundle) Contains(recipe *types.OpenInstallationRecipe) bool {

	for i := range b.BundleRecipes {
		if b.BundleRecipes[i].Recipe == recipe {
			return true
		}
	}

	return false
}

func (b *Bundle) ContainsName(name string) bool {

	for i, _ := range b.BundleRecipes {
		if b.BundleRecipes[i].Recipe.Name == name {
			return true
		}
	}

	return false
}
