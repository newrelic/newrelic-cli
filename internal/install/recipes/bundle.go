package recipes

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
)

type Bundle struct {
	BundleRecipes []*BundleRecipe
	Type          BundleType
}

// OpenInstallationOperatingSystem - Operating System of target environment
type BundleType string

var BundleTypes = struct {
	// MacOS operating system
	CORE BundleType
	// Linux-based operating system
	ADDITIONAL_GUIDED BundleType
	// Windows operating system
	ADDITIONAL_TARGETED BundleType
}{
	// MacOS operating system
	CORE: "CORE",
	// Linux-based operating system
	ADDITIONAL_GUIDED: "ADDITIONAL_GUIDED",
	// Windows operating system
	ADDITIONAL_TARGETED: "ADDITIONAL_TARGETED",
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

func (b *Bundle) AvailableRecipeCount() int {

	count := 0
	for i := 0; i < len(b.BundleRecipes); i++ {
		if b.BundleRecipes[i].HasStatus(execution.RecipeStatusTypes.AVAILABLE) {
			count++
		}
	}

	return count
}

func (b *Bundle) PrintRecipes() {

	for i := 0; i < len(b.BundleRecipes); i++ {
		fmt.Printf("\n%v: %v\n", i, b.BundleRecipes[i].Recipe.Name)
	}
}
