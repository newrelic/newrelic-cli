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
	HasSuperInstalled   bool
}

func NewBundler(context context.Context, availableRecipes RecipeDetectionResults) *Bundler {
	return &Bundler{
		Context:             context,
		AvailableRecipes:    availableRecipes,
		cachedBundleRecipes: make(map[string]*BundleRecipe),
		HasSuperInstalled:   false,
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

// TODO: Add doc for the following function
// dependency match - agent running
func (b *Bundler) createBundle(recipes []string, bType BundleType) *Bundle {
	bundle := &Bundle{Type: bType}

	for _, r := range recipes {
		if d, ok := b.AvailableRecipes.GetRecipeDetection(r); ok {
			var bundleRecipe *BundleRecipe
			if dualDep, ok := detectDependencies(d.Recipe.Dependencies); ok {
				dep := b.updateDependency(dualDep, recipes)
				if !utils.StringInSlice(dep[0], recipes) && b.HasSuperInstalled {
					d.Recipe.Dependencies = nil
				}
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

// getBundleRecipeWithDependencies retrieves a bundle recipe with its resolved dependencies.
// Parameters: recipe (types.OpenInstallationRecipe): The OpenInstallationRecipe object representing the bundle recipe.
// Returns: *BundleRecipe: A pointer to a BundleRecipe object containing the recipe and its resolved dependencies,
// or nil if there are missing dependencies or errors.
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

// detectDependencies evaluates if a recipe's dependency comes in the form 'recipe-a || recipe-b' and
// if detected it returns that dependency line content along with a 'true' found value.
func detectDependencies(deps []string) (string, bool) {
	if len(deps) == 0 {
		return "", false
	}

	const dualRecipeDependencyRegex = `^.+\|\|.+$` // e.g.: infrastructure-agent-installer || super-agent
	r, _ := regexp.Compile(dualRecipeDependencyRegex)

	// Not yet considering the unlikely case of dealing with more than one recipe dependency line coming in the 'recipe-a || recipe-b' form
	for _, dep := range deps {
		if r.MatchString(dep) {
			return dep, true
		}
	}

	return "", false
}

// updateDependency updates a recipe's dependency with the first one of the form 'recipe-a || recipe-b' that is found in the targeted
// recipes (e.g.: newrelic install -n recipe-a,recipe-c). If none of the recipe's dependency in the form 'recipe-a || recipe-b' are found
// in the targeted recipes, the first one of those in that same form 'recipe-a || recipe-b' is used. The final result is that recipe's
// dependency will change from the form 'recipe-a || recipe-b' to, for example, 'recipe-a' only.
func (b *Bundler) updateDependency(dualDep string, recipes []string) []string {
	var (
		splitDeps   = strings.Split(dualDep, `||`)
		hasSuperDep bool
	)

	if len(splitDeps) <= 1 {
		return nil
	}
	// TODO: Update the doc
	for _, dep := range splitDeps {
		dep = strings.TrimSpace(dep)
		if dep == types.SuperAgentRecipeName {
			hasSuperDep = true
			break
		}
	}
	if hasSuperDep && b.HasSuperInstalled {
		return []string{types.SuperAgentRecipeName}
	}

	for _, dep := range splitDeps {
		dep = strings.TrimSpace(dep)
		if utils.StringInSlice(dep, recipes) {
			return []string{dep}
		}
	}

	return []string{strings.TrimSpace(splitDeps[0])}
}
