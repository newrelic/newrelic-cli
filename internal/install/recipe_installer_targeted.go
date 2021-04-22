package install

import (
	"context"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func (i *RecipeInstaller) resolveRecipeDependencies(recipe types.Recipe, manifest *types.DiscoveryManifest) []*types.Recipe {
	log.Print("\n\n **************************** \n")
	log.Printf("\n resolveRecipeDependencies:  %+v \n", recipe.Dependencies)
	log.Print("\n **************************** \n\n")

	if len(recipe.Dependencies) == 0 {
		return []*types.Recipe{}
	}

	dependencies := []*types.Recipe{}
	for _, d := range recipe.Dependencies {
		manifest.OS = "linux"
		recipe := i.fetchWarn(manifest, d)

		log.Print("\n\n **************************** \n")
		log.Printf("\n fetched recipe:  %+v \n", recipe)
		log.Print("\n **************************** \n\n")

		if recipe != nil {
			dependencies = append(dependencies, recipe)
		}
	}

	return dependencies
}

func (i *RecipeInstaller) targetedInstall(ctx context.Context, m *types.DiscoveryManifest) error {
	var recipes []types.Recipe
	var infraAgentRecipe *types.Recipe

	i.status.SetTargetedInstall()

	if i.RecipePathsProvided() {
		// Load the recipes from the provided file names.
		for _, n := range i.RecipePaths {
			log.Debugln(fmt.Sprintf("Attempting to match recipePath %s.", n))
			recipe, err := i.recipeFromPath(n)
			if err != nil {
				log.Debugln(fmt.Sprintf("Error while building recipe from path, detail:%s.", err))
				return err
			}

			log.WithFields(log.Fields{
				"name":         recipe.Name,
				"display_name": recipe.DisplayName,
				"path":         n,
			}).Debug("found recipe at path")

			log.Print("\n\n **************************** \n")
			log.Printf("\n manifest:            %+v \n", *m)
			log.Printf("\n recipeDependencies:  %+v \n", recipe)

			// recipeDependencies := i.resolveRecipeDependencies(*recipe, m)

			log.Print("\n **************************** \n\n")
			time.Sleep(3 * time.Second)

			if !i.SkipInfra && recipe.Name == types.InfraAgentRecipeName {
				infraAgentRecipe = recipe
			} else {
				recipes = append(recipes, *recipe)
			}
		}
	} else if i.RecipeNamesProvided() {
		// Fetch the provided recipes from the recipe service.
		for _, n := range i.RecipeNames {
			log.Debugln(fmt.Sprintf("Attempting to match recipeName %s.", n))
			r := i.fetchWarn(m, n)
			if r != nil {
				// Skip anything that was returned by the service if it does not match the requested name.
				if r.Name == n {
					log.Debugln(fmt.Sprintf("Found recipe from name %s.", n))
					if r.Name == types.InfraAgentRecipeName {
						infraAgentRecipe = r
					} else {
						recipes = append(recipes, *r)
					}
				} else {
					log.Debugln(fmt.Sprintf("Skipping recipe, name %s does not match.", r.Name))
				}
			}
		}
	}

	if !i.SkipInfra && infraAgentRecipe == nil {
		fmt.Printf("The installation will begin by installing the latest version of the New Relic Infrastructure agent, which is required for additional instrumentation.\n\n")
		// Fetch the infra agent recipe and mark it as available.
		recipe, err := i.fetchRecipeAndReportAvailable(ctx, m, types.InfraAgentRecipeName)
		if err != nil {
			return err
		}
		infraAgentRecipe = recipe
	}

	// Show the user what will be installed.
	i.status.RecipesAvailable(recipes)

	if !i.SkipInfra {
		i.status.RecipesSelected(append([]types.Recipe{*infraAgentRecipe}, recipes...))

		// Install the infra agent.
		log.Debugf("Installing infrastructure agent")
		_, err := i.executeAndValidateWithProgress(ctx, m, infraAgentRecipe)
		if err != nil {
			log.Error(i.failMessage(types.InfraAgentRecipeName))
			return err
		}
		log.Debugf("Done installing infrastructure agent.")
	}

	// log.Print("\n\n **************************** \n")
	// log.Printf("\n RECIPES:  %+v \n", recipes)
	// log.Print("\n **************************** \n\n")
	// time.Sleep(2 * time.Second)

	// Install the requested integrations.
	log.Debugf("Installing integrations")
	if err := i.installRecipes(ctx, m, recipes); err != nil {
		return err
	}

	log.Debugf("Done installing integrations.")

	return nil
}

func (i *RecipeInstaller) recipeFromPath(recipePath string) (*types.Recipe, error) {
	recipeURL, parseErr := url.Parse(recipePath)
	if parseErr == nil && recipeURL.Scheme != "" {
		f, err := i.recipeFileFetcher.FetchRecipeFile(recipeURL)
		if err != nil {
			return nil, fmt.Errorf("could not fetch file %s: %s", recipePath, err)
		}
		return finalizeRecipe(f)
	}

	f, err := i.recipeFileFetcher.LoadRecipeFile(recipePath)
	if err != nil {
		return nil, fmt.Errorf("could not load file %s: %s", recipePath, err)
	}

	return finalizeRecipe(f)
}

func finalizeRecipe(f *recipes.RecipeFile) (*types.Recipe, error) {
	r, err := f.ToRecipe()
	if err != nil {
		return nil, fmt.Errorf("could not finalize recipe %s: %s", f.Name, err)
	}
	return r, nil
}

func (i *RecipeInstaller) fetchWarn(m *types.DiscoveryManifest, recipeName string) *types.Recipe {
	r, err := i.recipeFetcher.FetchRecipe(utils.SignalCtx, m, recipeName)
	if err != nil {
		log.Warnf("Could not install %s. Error retrieving recipe: %s", recipeName, err)
		return nil
	}

	if r == nil {
		log.Warnf("Recipe %s not found. Skipping installation.", recipeName)
	}

	return r
}
