package recipes

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

var coreRecipeMap = map[string]bool{
	types.InfraAgentRecipeName: true,
	types.LoggingRecipeName:    true,
}

type Bundler struct {
	RecipeRepository    *RecipeRepository
	RecipeDetector      *RecipeDetector
	Context             context.Context
	cachedBundleRecipes map[string]*BundleRecipe
}

func NewBundler(context context.Context, rr *RecipeRepository) *Bundler {
	return &Bundler{
		Context:             context,
		RecipeRepository:    rr,
		RecipeDetector:      NewRecipeDetector(),
		cachedBundleRecipes: make(map[string]*BundleRecipe),
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

	allRecipes, _ := b.RecipeRepository.FindAll()
	for _, recipe := range allRecipes {
		recipes = append(recipes, recipe)
	}

	bundle := b.createBundle(recipes, BundleTypes.ADDITIONALTARGETED)

	for _, br := range bundle.BundleRecipes {
		if !Find(recipeNames, br.Recipe.Name) {
			br.DetectedStatuses = br.RemoveStatus(execution.RecipeStatusTypes.AVAILABLE)
		}
	}

	return bundle
}

func Find(recipeNames []string, recipeSearchingFor string) bool {
	for _, i := range recipeNames {
		if i == recipeSearchingFor {
			return true
		}
	}
	return false
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
		bundleRecipe := b.getBundleRecipeWithDependencies(r)
		if bundleRecipe != nil {
			b.RecipeDetector.detectBundleRecipe(b.Context, bundleRecipe)
			log.Debugf("Adding bundle recipe:%s status:%+v dependencies:%+v", bundleRecipe.Recipe.Name, bundleRecipe.DetectedStatuses, bundleRecipe.Recipe.Dependencies)
			bundle.AddRecipe(bundleRecipe)
		}
	}

	return bundle
}

func (b *Bundler) getBundleRecipeWithDependencies(recipe *types.OpenInstallationRecipe) *BundleRecipe {
	if br, ok := b.cachedBundleRecipes[recipe.Name]; ok {
		return br
	}

	bundleRecipe := &BundleRecipe{
		Recipe: recipe,
	}
	b.cachedBundleRecipes[recipe.Name] = bundleRecipe

	for _, d := range recipe.Dependencies {
		if r := b.RecipeRepository.FindRecipeByName(d); r != nil {
			dr := b.getBundleRecipeWithDependencies(r)
			if dr != nil {
				bundleRecipe.Dependencies = append(bundleRecipe.Dependencies, dr)
				continue
			} else {
				log.Debugf("dependent bundle recipe %s not found, skipping recipe %s", d, recipe.Name)
			}
		} else {
			log.Debugf("dependent recipe %s not found, skipping recipe %s", d, recipe.Name)
		}
		// A dependency is missing, invalidating the bundle recipe
		b.cachedBundleRecipes[recipe.Name] = nil
		return nil
	}

	return bundleRecipe
}
