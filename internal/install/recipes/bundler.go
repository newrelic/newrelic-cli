package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"

	log "github.com/sirupsen/logrus"

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

func (b *Bundler) CreateAdditionalBundle() *Bundle {

	coreRecipeNameMap := make(map[string]bool)
	for _, recipeName := range coreBundleRecipeNames {
		coreRecipeNameMap[recipeName] = true
	}

	//TODO: check dependency graph at install? Not counting dependency here could cause it to prompt during additional install
	//TODO: Target-install logic here
	var additionalRecipes []*types.OpenInstallationRecipe
	recipes, _ := b.RecipeRepository.FindAll()

	for _, recipe := range recipes {
		if !coreRecipeNameMap[recipe.Name] {
			additionalRecipes = append(additionalRecipes, recipe)
		}
	}

	return b.CreateBundle(additionalRecipes)
}

func (b *Bundler) CreateBundle(recipes []*types.OpenInstallationRecipe) *Bundle {

	bundle := &Bundle{}

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
