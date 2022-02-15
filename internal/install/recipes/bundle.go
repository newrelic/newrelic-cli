package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

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
