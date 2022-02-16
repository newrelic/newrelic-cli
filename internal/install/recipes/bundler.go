package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var coreBundleRecipeNames = []string{
	types.InfraAgentRecipeName,
	types.LoggingRecipeName,
	types.GoldenRecipeName,
}

type Bundler struct {
	RecipeRepository *RecipeRepository
	RecipeDetector   *RecipeDetector
	Context          context.Context
}

func NewBundler(context context.Context, rr *RecipeRepository) *Bundler {
	return newBundler(context, rr, NewRecipeDetector())
}

func newBundler(context context.Context, rr *RecipeRepository, rd *RecipeDetector) *Bundler {
	return &Bundler{
		Context:          context,
		RecipeRepository: rr,
		RecipeDetector:   rd,
	}
}

func (b *Bundler) CreateCoreBundle() *Bundle {
	var coreRecipes []*types.OpenInstallationRecipe
	for _, recipeName := range coreBundleRecipeNames {
		if r := b.RecipeRepository.FindRecipeByName(recipeName); r != nil {
			coreRecipes = append(coreRecipes, r)
		}
	}

	return b.CreateBundle(coreRecipes)
}

func (b *Bundler) CreateBundle(recipes []*types.OpenInstallationRecipe) *Bundle {

	bundle := &Bundle{}

	for _, r := range recipes {
		// recipe shouldn't have itself as dependency
		visited := map[string]bool{r.Name: true}
		bundleRecipe := b.getBundleRecipeWithDependencies(r, visited)
		bundle.AddRecipe(bundleRecipe)
	}

	b.RecipeDetector.DetectBundle(b.Context, bundle)

	return bundle
}

func (b *Bundler) CreateBundleRecipe(recipe *types.OpenInstallationRecipe) *BundleRecipe {

	visited := map[string]bool{recipe.Name: true}
	return b.getBundleRecipeWithDependencies(recipe, visited)
}

func (b *Bundler) getBundleRecipeWithDependencies(recipe *types.OpenInstallationRecipe, visited map[string]bool) *BundleRecipe {

	bundleRecipe := &BundleRecipe{
		Recipe: recipe,
	}

	for _, d := range recipe.Dependencies {
		if !visited[d] {
			visited[d] = true
			if r := b.RecipeRepository.FindRecipeByName(d); r != nil {
				dr := b.getBundleRecipeWithDependencies(r, visited)
				bundleRecipe.Dependencies = append(bundleRecipe.Dependencies, dr)
			}
		}
	}

	return bundleRecipe
}

// Control Status

//Recipe Candidate, recipe + collection of Status
//Recipe context, capturing recipe intall info, timing, Status..etc.

// func (b *Bundler) createAdditionalBundle(recipes []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
// 	_, a := createBundles(coreBundleRecipeNames, recipes)
// 	return a
// }
