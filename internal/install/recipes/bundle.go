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
	ADDITIONALGUIDED BundleType
	// Windows operating system
	ADDITIONALTARGETED BundleType
}{
	// MacOS operating system
	CORE: "CORE",
	// Linux-based operating system
	ADDITIONALGUIDED: "ADDITIONALGUIDED",
	// Windows operating system
	ADDITIONALTARGETED: "ADDITIONALTARGETED",
}

func (b *Bundle) AddRecipe(bundleRecipe *BundleRecipe) {
	if b.ContainsName(bundleRecipe.Recipe.Name) {
		return
	}
	b.BundleRecipes = append(b.BundleRecipes, bundleRecipe)
}

func (b *Bundle) IsAdditionalGuided() bool {
	return b.Type == BundleTypes.ADDITIONALGUIDED
}

func (b *Bundle) IsAdditionalTargeted() bool {
	return b.Type == BundleTypes.ADDITIONALTARGETED
}

func (b *Bundle) ContainsName(name string) bool {

	for i := range b.BundleRecipes {
		if b.BundleRecipes[i].Recipe.Name == name {
			return true
		}
	}

	return false
}

func (b *Bundle) GetBundleRecipe(name string) *BundleRecipe {

	for _, r := range b.BundleRecipes {
		if r.Recipe.Name == name {
			return r
		}
	}

	return nil
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

func (b *Bundle) RemoveBundleRecipe(name string) {

	for i, r := range b.BundleRecipes {
		if r.Recipe.Name == name {
			b.BundleRecipes = append(b.BundleRecipes[0:i], b.BundleRecipes[i+1:]...)
			return
		}
	}
}

func (b *Bundle) String() string {
	result := fmt.Sprintf("%s %s", b.Type, b.BundleRecipes)
	return fmt.Sprintf("{%s}", result)
}
