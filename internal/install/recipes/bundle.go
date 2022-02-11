package recipes

import (
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleRecipe struct {
	recipe   *types.OpenInstallationRecipe
	statuses []execution.RecipeStatusType
}

func NewBundleRecipe(recipe *types.OpenInstallationRecipe) *BundleRecipe {
	return &BundleRecipe{
		recipe: recipe,
	}
}

type Bundle struct {
	BundleRecipes []*BundleRecipe
}

func NewBundle(recipes []*types.OpenInstallationRecipe) *Bundle {

	bundle := &Bundle{}
	for i := 0; i < len(recipes); i++ {
		bundle.BundleRecipes = append(bundle.BundleRecipes,
			NewBundleRecipe(recipes[i]))
	}

	return bundle
}

func (b *Bundle) AddRecipes(recipes []*types.OpenInstallationRecipe) {
	for i := 0; i < len(recipes); i++ {
		b.AddRecipe(recipes[i])
	}
}

func (b *Bundle) AddRecipe(recipe *types.OpenInstallationRecipe) {
	if index, exists := b.ContainsName(recipe.Name); exists {
		b.BundleRecipes[index] = &BundleRecipe{
			recipe: recipe,
		}
	} else {
		newItem := []*BundleRecipe{
			{
				recipe: recipe,
			}}

		b.BundleRecipes = append(newItem, b.BundleRecipes...)
	}
}

func (b *Bundle) Contains(recipe *types.OpenInstallationRecipe) (int, bool) {

	for i, _ := range b.BundleRecipes {
		if b.BundleRecipes[i].recipe == recipe {
			return i, true
		}
	}

	return -1, false
}

func (b *Bundle) ContainsName(name string) (int, bool) {

	for i, _ := range b.BundleRecipes {
		if b.BundleRecipes[i].recipe.Name == name {
			return i, true
		}
	}

	return -1, false
}
