package recipes

import (
	"context"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
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
	return b.createBundle(b.GetCoreRecipeNames(), BundleTypes.CORE)
}

func (b *Bundler) CreateAdditionalGuidedBundle() *Bundle {
	var recipes []string

	for _, d := range b.AvailableRecipes {
		if !coreRecipeMap[d.Recipe.Name] {
			if !strings.EqualFold(d.Recipe.Name, types.AgentControlRecipeName) {
				recipes = append(recipes, d.Recipe.Name)
			}
		}
	}

	return b.createBundle(recipes, BundleTypes.ADDITIONALGUIDED)
}

func (b *Bundler) CreateAdditionalTargetedBundle(recipes []string) *Bundle {
	return b.createBundle(recipes, BundleTypes.ADDITIONALTARGETED)
}

func (b *Bundler) GetCoreRecipeNames() []string {
	coreRecipeNames := make([]string, 0, len(coreRecipeMap))
	for k := range coreRecipeMap {
		coreRecipeNames = append(coreRecipeNames, k)
	}
	return coreRecipeNames
}

// createBundle creates a new bundle based on the given recipes and bundle type
// It iterates over the recipes, detects dependencies, and adds recipes to the bundle
func (b *Bundler) createBundle(recipes []string, bType BundleType) *Bundle {
	bundle := &Bundle{Type: bType}

	for _, r := range recipes {
		if d, ok := b.AvailableRecipes.GetRecipeDetection(r); ok {
			var bundleRecipe *BundleRecipe

			if dualDep, ok := detectDependencies(d.Recipe.Dependencies); ok {
				dep := b.updateDependency(dualDep, recipes)
				if dep != nil {
					log.Debugf("Found dual dependency and selected : %s", dep)
					d.Recipe.Dependencies = dep
				} else {
					log.Debugf("could not process update for dual dependency: %s", dualDep)
				}
			}

			bundleRecipe = b.getBundleRecipeWithDependencies(d.Recipe)

			if bundleRecipe != nil {
				isAgentControlTargetedInput := utils.StringInSlice(types.AgentControlRecipeName, recipes)
				if dep, found := findRecipeDependency(bundleRecipe, types.AgentControlRecipeName); !isAgentControlTargetedInput && found {
					log.Debugf("updating the dependency status for %s", dep)
					dep.AddDetectionStatus(execution.RecipeStatusTypes.INSTALLED, 0)
				}
				if !isAgentControlTargetedInput && b.HasSuperInstalled {
					log.Debugf("Agent Control found, skipping")
					superRecipe := &BundleRecipe{
						Recipe: &types.OpenInstallationRecipe{
							Name: types.AgentControlRecipeName,
						},
					}
					superRecipe.AddDetectionStatus(execution.RecipeStatusTypes.INSTALLED, 0)
					bundle.AddRecipe(superRecipe)
				}
				log.Debugf("Adding bundle recipe:%s status:%+v dependencies:%+v", bundleRecipe.Recipe.Name, bundleRecipe.DetectedStatuses, bundleRecipe.Recipe.Dependencies)
				bundle.AddRecipe(bundleRecipe)
			}
		}
	}

	return bundle
}

// findRecipeDependency recursively searches for a recipe dependency
func findRecipeDependency(recipe *BundleRecipe, name string) (*BundleRecipe, bool) {
	if strings.EqualFold(recipe.Recipe.Name, name) {
		return recipe, true
	}
	for _, dep := range recipe.Dependencies {
		if found, ok := findRecipeDependency(dep, name); ok {
			return found, true
		}
	}
	return nil, false
}

// getBundleRecipeWithDependencies returns a BundleRecipe for the given recipe, including its dependencies.
// If the recipe is already cached, it is returned immediately. Otherwise, the function recursively
// fetches the dependencies and builds the BundleRecipe. If any dependencies are missing, the
// bundle recipe is invalidated and nil is returned.
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
		log.Debugf("A dependency is missing, invalidating the bundle recipe %s", recipe.Name)
		b.cachedBundleRecipes[recipe.Name] = nil
		return nil
	}

	if bundleRecipe.AreAllDependenciesAvailable() {
		if dt, ok := b.AvailableRecipes.GetRecipeDetection(recipe.Name); ok {
			bundleRecipe.AddDetectionStatus(dt.Status, dt.DurationMs)
			b.cachedBundleRecipes[recipe.Name] = bundleRecipe
			log.Debugf("bundle recipe with dependency %v", bundleRecipe)
			return bundleRecipe
		}
	}

	b.cachedBundleRecipes[recipe.Name] = nil
	log.Debugf("Returning nil for the recipe: %s", recipe.Name)
	return nil
}

// detectDependencies evaluates if a recipe's dependency comes in the form 'recipe-a || recipe-b' and
// if detected it returns that dependency line content along with a 'true' found value.
func detectDependencies(deps []string) (string, bool) {
	if len(deps) == 0 {
		return "", false
	}

	const dualRecipeDependencyRegex = `^.+\|\|.+$` // e.g.: infrastructure-agent-installer || agent-control
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
// dependency will change from the form 'recipe-a || recipe-b' to the first recipe in that list, for example, 'recipe-a' only.
// If the dependency is a super dependency (i.e., '||' is present), and the super dependency is installed, the function will return the super dependency.
// If the dependency is not a super dependency and is found in the targeted recipes, the function will return that dependency.
// If the dependency is not a super dependency and is not found in the targeted recipes, the function will return the first part of the dependency.
func (b *Bundler) updateDependency(dualDep string, recipes []string) []string {
	var (
		splitDeps   = strings.Split(dualDep, `||`)
		hasSuperDep bool
	)

	if len(splitDeps) <= 1 {
		return nil
	}
	for _, dep := range splitDeps {
		dep = strings.TrimSpace(dep)
		if dep == types.AgentControlRecipeName {
			hasSuperDep = true
			break
		}
	}
	if hasSuperDep && b.HasSuperInstalled {
		return []string{types.AgentControlRecipeName}
	}

	for _, dep := range splitDeps {
		dep = strings.TrimSpace(dep)
		if utils.StringInSlice(dep, recipes) {
			return []string{dep}
		}
	}

	return []string{strings.TrimSpace(splitDeps[0])}
}

// IsCore checks if a recipe is a core recipe
func (b *Bundler) IsCore(recipeName string) bool {
	return coreRecipeMap[recipeName]
}
