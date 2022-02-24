package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var coreRecipeMap = map[string]bool{
	types.InfraAgentRecipeName: true,
	types.LoggingRecipeName:    true,
	types.GoldenRecipeName:     true,
}

type Bundler struct {
	RecipeRepository *RecipeRepository
	RecipeDetector   *RecipeDetector
	Context          context.Context
}

func NewBundler(context context.Context, rr *RecipeRepository) *Bundler {
	return &Bundler{
		Context:          context,
		RecipeRepository: rr,
		RecipeDetector:   NewRecipeDetector(),
	}
}

func (b *Bundler) CreateCoreBundle() *Bundle {
	var recipes []*types.OpenInstallationRecipe

	for _, recipeName := range b.getCoreRecipeNames() {
		if r := b.RecipeRepository.FindRecipeByName(recipeName); r != nil {
			recipes = append(recipes, r)
		}
	}

	return b.createBundle(recipes, BundleTypes.CORE)
}

func (b *Bundler) CreateAdditionalGuidedBundle() *Bundle {
	var recipes []*types.OpenInstallationRecipe

	allRecipes, _ := b.RecipeRepository.FindAll()
	for _, recipe := range allRecipes {
		if !coreRecipeMap[recipe.Name] {
			recipes = append(recipes, recipe)
		}
	}

	return b.createBundle(recipes, BundleTypes.ADDITIONALGUIDED)
}

func (b *Bundler) CreateAdditionalTargetedBundle(recipeNames []string) *Bundle {
	var recipes []*types.OpenInstallationRecipe

	for _, recipeName := range recipeNames {
		if coreRecipeMap[recipeName] {
			continue
		}
		if r := b.RecipeRepository.FindRecipeByName(recipeName); r != nil {
			recipes = append(recipes, r)
		}
	}

	return b.createBundle(recipes, BundleTypes.ADDITIONALTARGETED)
}

func (b *Bundler) CreateAdditionalTargetedPathBundle(recipes []*types.OpenInstallationRecipe) *Bundle {
	return b.createBundle(recipes, BundleTypes.ADDITIONALTARGETED)
}

func (b *Bundler) getCoreRecipeNames() []string {
	coreRecipeNames := make([]string, 0, len(coreRecipeMap))
	for k := range coreRecipeMap {
		coreRecipeNames = append(coreRecipeNames, k)
	}
	return coreRecipeNames
}

func (b *Bundler) createBundle(recipes []*types.OpenInstallationRecipe, bType BundleType) *Bundle {
	bundle := &Bundle{Type: bType}

	for _, r := range recipes {
		// recipe shouldn't have itself as dependency
		visited := map[string]bool{r.Name: true}
		bundleRecipe := b.getBundleRecipeWithDependencies(r, visited)

		if bundleRecipe != nil {
			log.Debugf("Adding bundle recipe:%s status:%+v dependencies:%+v", bundleRecipe.Recipe.Name, bundleRecipe.DetectedStatuses, bundleRecipe.Recipe.Dependencies)
			bundle.AddRecipe(bundleRecipe)
		}
	}

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

	//this is the parent
	//FIXME: don't like returning nil
	b.RecipeDetector.detectBundleRecipe(b.Context, bundleRecipe)
	if bundleRecipe.HasStatus(execution.RecipeStatusTypes.NULL) {
		return nil
	}

	for _, d := range recipe.Dependencies {
		if !visited[d] {
			visited[d] = true
			if r := b.RecipeRepository.FindRecipeByName(d); r != nil {
				dr := b.getBundleRecipeWithDependencies(r, visited)
				if dr != nil {
					bundleRecipe.Dependencies = append(bundleRecipe.Dependencies, dr)
				}
			}
		}
	}

	return bundleRecipe
}
