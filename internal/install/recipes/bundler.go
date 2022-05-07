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
	Detections          map[string]*RecipeDetection
	Context             context.Context
	cachedBundleRecipes map[string]*BundleRecipe
}

func NewBundler(context context.Context, detections map[string]*RecipeDetection) *Bundler {
	return &Bundler{
		Context:             context,
		Detections:          detections,
		cachedBundleRecipes: make(map[string]*BundleRecipe),
	}
}

func (b *Bundler) CreateCoreBundle() *Bundle {
	return b.createBundle(b.getCoreRecipeNames(), BundleTypes.CORE)
}

func (b *Bundler) CreateAdditionalGuidedBundle() *Bundle {
	var recipes []string

	for _, d := range b.Detections {
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
		if d, ok := b.Detections[r]; ok {
			bundleRecipe := b.getBundleRecipeWithDependencies(d.Recipe)
			if bundleRecipe != nil {
				log.Debugf("Adding bundle recipe:%s status:%+v dependencies:%+v", bundleRecipe.Recipe.Name, bundleRecipe.DetectedStatuses, bundleRecipe.Recipe.Dependencies)
				bundle.AddRecipe(bundleRecipe)
			}
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

	for _, d := range recipe.Dependencies {
		if dt, ok := b.Detections[d]; ok {
			dr := b.getBundleRecipeWithDependencies(dt.Recipe)
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
		if dt, ok := b.Detections[recipe.Name]; ok {
			if dt.Status == execution.RecipeStatusTypes.AVAILABLE {
				bundleRecipe.AddDetectionStatus(dt.Status, dt.DurationMs)
				b.cachedBundleRecipes[recipe.Name] = bundleRecipe
				return bundleRecipe
			}
		}
	}

	b.cachedBundleRecipes[recipe.Name] = nil
	return nil
}
