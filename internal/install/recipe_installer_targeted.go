package install

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func (i *RecipeInstaller) targetedInstall(m *types.DiscoveryManifest) error {
	var err error
	var recipes []types.Recipe

	if i.RecipePathsProvided() {
		// Load the recipes from the provided file names.
		for _, n := range i.RecipePaths {
			log.Debugln(fmt.Sprintf("Attempting to match recipePath %s.", n))
			var recipe *types.Recipe
			recipe, err = i.recipeFromPath(n)
			if err != nil {
				log.Debugln(fmt.Sprintf("Error while building recipe from path, detail:%s.", err))
				return err
			}

			log.WithFields(log.Fields{
				"name":         recipe.Name,
				"display_name": recipe.DisplayName,
				"path":         n,
			}).Debug("found recipe at path")

			recipes = append(recipes, *recipe)
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
					recipes = append(recipes, *r)
				} else {
					log.Debugln(fmt.Sprintf("Skipping recipe, name %s does not match.", r.Name))
				}
			}
		}
	}

	fmt.Printf("The installation will begin by installing the latest version of the New Relic Infrastructure agent, which is required for additional instrumentation.\n\n")

	// Fetch the infra agent recipe and mark it as available.
	infraAgentRecipe, err := i.fetchRecipeAndReportAvailable(m, infraAgentRecipeName)
	if err != nil {
		return err
	}

	// Show the user what will be installed.
	i.status.RecipesAvailable(recipes)
	i.status.RecipesSelected(append([]types.Recipe{*infraAgentRecipe}, recipes...))

	// Install the infra agent.
	log.Debugf("Installing infrastructure agent")
	_, err = i.executeAndValidateWithProgress(m, infraAgentRecipe)
	if err != nil {
		log.Error(i.failMessage(infraAgentRecipeName))
		return err
	}
	log.Debugf("Done installing infrastructure agent.")

	// Install the requested integrations.
	log.Debugf("Installing integrations")
	if err = i.installRecipes(m, recipes); err != nil {
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
