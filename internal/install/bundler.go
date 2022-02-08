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

func (b *Bundler) createCoreBundle() []types.OpenInstallationRecipe {
	var core []types.OpenInstallationRecipe
	for _, recipeName := range coreBundleRecipeNames {
		if r := b.RecipeRepository.FindRecipeByName(recipeName); r != nil {
			core = append(core, *r)
		}
	}

	// TODO: continue here
	return core
}

//Recipe Candidate, recipe + collection of status
//Recipe context, capturing recipe intall info, timing, status..etc.

// func (b *Bundler) createAdditionalBundle(recipes []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
// 	_, a := createBundles(coreBundleRecipeNames, recipes)
// 	return a
// }
