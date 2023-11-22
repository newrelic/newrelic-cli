package recipes

import (
	"context"
	"regexp"
	"strings"

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
			if dualDep, ok := getDualDependency(d.Recipe.Dependencies); ok {
				dep := updateDependency(dualDep, recipes)
				if dep != nil {
					d.Recipe.Dependencies = dep
				} else {
					log.Debugf("could not process update for dual dependency: %s", dualDep)
				}
			}
			bundleRecipe = b.getBundleRecipeWithDependencies(d.Recipe)
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
		if dt, ok := b.AvailableRecipes.GetRecipeDetection(d); ok {
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
		if dt, ok := b.AvailableRecipes.GetRecipeDetection(recipe.Name); ok {
			bundleRecipe.AddDetectionStatus(dt.Status, dt.DurationMs)
			b.cachedBundleRecipes[recipe.Name] = bundleRecipe
			return bundleRecipe
		}
	}

	b.cachedBundleRecipes[recipe.Name] = nil
	return nil
}

func getDualDependency(deps []string) (string, bool) {
	if len(deps) == 0 {
		return "", false
	}

	const dualRecipeDependencyRegex = `^.+\|\|.+$` // e.g.: infrastructure-agent-installer || super-agent
	r, _ := regexp.Compile(dualRecipeDependencyRegex)

	// Not yet considering the unlikely case of dealing with more than one recipe dependency line coming in the 'a || b' form
	for _, dep := range deps {
		if r.MatchString(dep) {
			return dep, true
		}
	}

	return "", false
}

func updateDependency(dualDep string, recipes []string) []string {
	var deps []string

	for _, dep := range strings.Split(dualDep, `||`) {
		dep = strings.TrimSpace(dep)
		if utils.StringInSlice(dep, recipes) {
			deps = []string{dep}
			break
		} else {
			deps = []string{types.InfraAgentRecipeName}
		}
	}

	return deps
}
