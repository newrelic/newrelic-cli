package install

import (
	recipes "github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var coreBundleRecipeNames = []string{
	types.InfraAgentRecipeName,
	types.LoggingRecipeName,
	types.GoldenRecipeName,
}

type Bundler struct {
	RecipeRepository *recipes.RecipeRepository
}

func NewBundler(rr *recipes.RecipeRepository) *Bundler {
	return &Bundler{
		RecipeRepository: rr,
	}
}

func (b *Bundler) createCoreBundle() *Bundle {
	var coreRecipes []*types.OpenInstallationRecipe
	for _, recipeName := range coreBundleRecipeNames {
		if r := b.RecipeRepository.FindRecipeByName(recipeName); r != nil {
			coreRecipes = append(coreRecipes, r)
		}
	}

	return b.createBundle(coreRecipes)
}

func (b *Bundler) createBundle(recipes []*types.OpenInstallationRecipe) *Bundle {

	coreBundle := NewBundle(recipes)
	coreBundle = b.addBundleDependencies(coreBundle)

	// TODO: do detection here, and there

	return coreBundle
}

func (b *Bundler) addBundleDependencies(bundle *Bundle) *Bundle {

	dependencies := b.getBundleDependencies(bundle)
	bundle.AddRecipes(dependencies)

	return bundle
}

// This is a naive implementation that only resolves dependencies one level deep.
func (b *Bundler) getBundleDependencies(bundle *Bundle) []*types.OpenInstallationRecipe {
	var results []*types.OpenInstallationRecipe
	found := make(map[string]bool, 0)

	for _, br := range bundle.BundleRecipes {
		if len(br.recipe.Dependencies) > 0 {
			for _, d := range br.recipe.Dependencies {
				if r := b.RecipeRepository.FindRecipeByName(d); r != nil {
					if r != nil && !found[r.Name] {
						results = append(results, r)
						found[r.Name] = true
					}
				}
			}
		}
	}

	return results
}

// Control status
type BundleInstaller struct {
}

//Recipe Candidate, recipe + collection of status
//Recipe context, capturing recipe intall info, timing, status..etc.

// func (b *Bundler) createAdditionalBundle(recipes []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
// 	_, a := createBundles(coreBundleRecipeNames, recipes)
// 	return a
// }
