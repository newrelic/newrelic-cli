package recipes

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var coreRecipeMap = map[string]bool{
	types.InfraAgentRecipeName: true,
	types.LoggingRecipeName:    true,
}

type Bundler struct {
	AvailableRecipes    RecipeDetectionResults
	Context             context.Context
	cachedBundleRecipes map[string]*BundleRecipe
}

func NewBundler(context context.Context, availableRecipes RecipeDetectionResults) *Bundler {
	return &Bundler{
		Context:             context,
		AvailableRecipes:    availableRecipes,
		cachedBundleRecipes: make(map[string]*BundleRecipe),
	}
}

func (b *Bundler) CreateCoreBundle() *Bundle {
	return b.createBundle(b.getCoreRecipeNames(), BundleTypes.CORE)
}

func (b *Bundler) CreateAdditionalGuidedBundle() *Bundle {
	var recipes []string

	for _, d := range b.AvailableRecipes {
		if !coreRecipeMap[d.Recipe.Name] {
			recipes = append(recipes, d.Recipe.Name)
		}
	}

	return b.createBundle(recipes, BundleTypes.ADDITIONALGUIDED)
}

func (b *Bundler) CreateAdditionalTargetedBundle(recipes []string) *Bundle {
	return b.createBundle(recipes, BundleTypes.ADDITIONALTARGETED)
}

func (b *Bundler) getCoreRecipeNames() []string {
	coreRecipeNames := make([]string, 0, len(coreRecipeMap))
	for k := range coreRecipeMap {
		coreRecipeNames = append(coreRecipeNames, k)
	}
	return coreRecipeNames
}

func (b *Bundler) createBundle(recipes []string, bType BundleType) *Bundle {
	bundle := &Bundle{Type: bType}

	for _, r := range recipes {
		if d, ok := b.AvailableRecipes.GetRecipeDetection(r); ok {
			var bundleRecipe *BundleRecipe
			// OHIs with dual infra/super dependency: if super-agent is targeted, remove infrastructure-agent
			if utils.StringInSlice(types.SuperAgentRecipeName, recipes) {
				bundleRecipe = b.getBundleRecipeWithDependencies(d.Recipe, types.InfraAgentRecipeName)
			} else {
				// OHIs with dual infra/super dependency: if super-agent is not targeted, remove super-agent
				bundleRecipe = b.getBundleRecipeWithDependencies(d.Recipe, types.SuperAgentRecipeName)
			}
			if bundleRecipe != nil {
				log.Debugf("Adding bundle recipe:%s status:%+v dependencies:%+v", bundleRecipe.Recipe.Name, bundleRecipe.DetectedStatuses, bundleRecipe.Recipe.Dependencies)
				bundle.AddRecipe(bundleRecipe)
			}
		}
	}

	return bundle
}

func (b *Bundler) getBundleRecipeWithDependencies(recipe *types.OpenInstallationRecipe, depException string) *BundleRecipe {
	if br, ok := b.cachedBundleRecipes[recipe.Name]; ok {
		return br
	}

	bundleRecipe := &BundleRecipe{
		Recipe: recipe,
	}

	// For OHIs with dual infrastructure-agent/super-agent dependencies, remove the one not needed:
	//	if super-agent is targeted, infrastructure-agent is removed
	//	if super-agent is not targeted, super-agent is removed
	if recipe.IsOhi() && utils.StringInSlice(types.InfraAgentRecipeName, recipe.Dependencies) && utils.StringInSlice(types.SuperAgentRecipeName, recipe.Dependencies) && depException != "" {
		recipe.Dependencies = utils.RemoveFromSlice(recipe.Dependencies, depException)
	}

	for _, d := range recipe.Dependencies {
		if dt, ok := b.AvailableRecipes.GetRecipeDetection(d); ok {
			var dr *BundleRecipe
			if depException == "" {
				dr = b.getBundleRecipeWithDependencies(dt.Recipe, "")
			} else {
				dr = b.getBundleRecipeWithDependencies(dt.Recipe, depException)
			}
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

	if bundleRecipe.AreAllDependenciesAvailable() {
		if dt, ok := b.AvailableRecipes.GetRecipeDetection(recipe.Name); ok {
			bundleRecipe.AddDetectionStatus(dt.Status, dt.DurationMs)
			b.cachedBundleRecipes[recipe.Name] = bundleRecipe
			return bundleRecipe
		}
	}

	b.cachedBundleRecipes[recipe.Name] = nil
	return nil
}
