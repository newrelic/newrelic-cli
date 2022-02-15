package recipes

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
