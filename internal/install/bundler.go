package install

import (
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var coreBundleRecipeNames = []string{
	types.InfraAgentRecipeName,
	types.LoggingRecipeName,
	types.GoldenRecipeName,
}

type Bundler struct {
	RecipeRepository *RecipeRepository
}

func NewBundler(rr *RecipeRepository) *Bundler {
	return &Bundler{
		RecipeRepository: rr,
	}
}

func (b *Bundler) createCoreBundle() []types.OpenInstallationRecipe {
	var core []types.OpenInstallationRecipe
	for _, recipeName := range coreBundleRecipeNames {
		// TODO: implement FindRecipeByName
		if r := b.RecipeRepository.FindRecipeByName(recipeName); r != nil {
			core = append(core, *r)
		}
	}

	// TODO: continue here
	return core

}

// func (b *Bundler) createAdditionalBundle(recipes []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
// 	_, a := createBundles(coreBundleRecipeNames, recipes)
// 	return a
// }

	return nil
}
