package recipes

import "fmt"

type Bundle struct {
	BundleRecipes []*BundleRecipe
}

func (b *Bundle) AddRecipe(bundleRecipe *BundleRecipe) {
	if b.ContainsName(bundleRecipe.Recipe.Name) {
		return
	}
	b.BundleRecipes = append(b.BundleRecipes, bundleRecipe)
}

func (b *Bundle) ContainsName(name string) bool {

	for i := range b.BundleRecipes {
		if b.BundleRecipes[i].Recipe.Name == name {
			return true
		}
	}

	return false
}

// Returns all recipes flatten with dependencies
func (b *Bundle) Flatten() map[string]bool {

	results := make(map[string]bool)
	for i := 0; i < len(b.BundleRecipes); i++ {
		recipeMap := b.BundleRecipes[i].Flatten()
		for key := range recipeMap {
			results[key] = true
		}
	}

	return results
}

func (b *Bundle) PrintRecipes() {

	for i := 0; i < len(b.BundleRecipes); i++ {
		fmt.Printf("\n%v: %v\n", i, b.BundleRecipes[i].Recipe.Name)
	}
}