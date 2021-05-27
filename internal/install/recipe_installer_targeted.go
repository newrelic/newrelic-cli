package install

import (
	"context"
	"fmt"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

func (i *RecipeInstaller) resolveRecipeDependencies(ctx context.Context, recipe types.OpenInstallationRecipe, manifest *types.DiscoveryManifest) ([]*types.OpenInstallationRecipe, error) {
	dependencies := []*types.OpenInstallationRecipe{}

	if len(recipe.Dependencies) == 0 {
		return dependencies, nil
	}

	for _, recipeName := range recipe.Dependencies {
		recipe, err := i.fetchRecipeAndReportAvailable(ctx, manifest, recipeName)
		if err != nil {
			return dependencies, err
		}

		if recipe != nil {
			dependencies = append(dependencies, recipe)
		}
	}

	return dependencies, nil
}

func (i *RecipeInstaller) collectRecipes(m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	var recipes []types.OpenInstallationRecipe

	if i.RecipePathsProvided() {
		// Load the recipes from the provided file names.
		for _, n := range i.RecipePaths {
			// Early continue when skipInfra is set
			if i.SkipInfra && n == types.InfraAgentRecipeName {
				continue
			}

			log.Debugln(fmt.Sprintf("Attempting to match recipePath %s.", n))
			recipe, err := i.recipeFromPath(n)
			if err != nil {
				log.Debugln(fmt.Sprintf("Error while building recipe from path, detail:%s.", err))
				return nil, err
			}

			log.WithFields(log.Fields{
				"name":         recipe.Name,
				"display_name": recipe.DisplayName,
				"path":         n,
			}).Debug("found recipe at path")

			recipes = append(recipes, *recipe)
		}
	} else if i.RecipeNamesProvided() {
		// matchedRecipes, err := i.recipeFetcher.FetchRecommendations(utils.SignalCtx, m)
		// if err != nil {
		// 	log.Debugf("error retrieving recipe recommendations: %s", err)
		// 	return recipes, err
		// }

		// Fetch the provided recipes from the recipe service.
		for _, n := range i.RecipeNames {
			// Early continue when skipInfra is set
			if i.SkipInfra && n == types.InfraAgentRecipeName {
				continue
			}

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

			recipes = append(recipes, *r)
		}
	}

	return recipes, nil
}

func (i *RecipeInstaller) targetedInstall(ctx context.Context, m *types.DiscoveryManifest) error {
	var recipes []types.OpenInstallationRecipe

	i.status.SetTargetedInstall()

	providedRecipes, err := i.collectRecipes(m)
	log.Print("\n\n **************************** \n")
	log.Printf("\n providedRecipes:  %+v \n", len(providedRecipes))
	log.Print("\n **************************** \n\n")
	if err != nil {
		return err
	}

	if len(providedRecipes) == 0 {
		fmt.Println("Nothing to install.")
		return nil
	}

	for _, r := range providedRecipes {
		dependencies, err := i.resolveRecipeDependencies(ctx, r, m)
		if err != nil {
			return err
		}

		for _, d := range dependencies {
			if i.SkipInfra && types.InfraAgentRecipeName == d.Name {
				continue
			} else {
				recipes = append(recipes, *d)
			}
		}
		recipes = append(recipes, r)
	}

	// Show the user what will be installed.
	i.status.RecipesAvailable(recipes)
	i.status.RecipesSelected(recipes)

	log.Print("\n\n **************************** \n")
	log.Printf("\n recipes:  %+v \n", len(recipes))
	log.Print("\n **************************** \n\n")

	for _, r := range recipes {
		if r.Name == types.InfraAgentRecipeName {
			continue
		}

		processesToMatch := r.ProcessMatch

		processesDiscovered := m.Processes

		log.Print("\n\n **************************** \n")
		log.Printf("\n processesToMatch:  %+v \n", processesToMatch)
		log.Printf("\n processesDiscovered:  %+v \n", processesDiscovered)
		log.Print("\n **************************** \n\n")
		time.Sleep(3 * time.Second)

		// if r == nil {
		// 	msg := fmt.Sprintf("please rerun `newrelic install` without `-n %s`", r)
		// 	err := fmt.Errorf("could not install because a dependent process was not detected, %s", msg)
		// 	return err
		// }
	}

	// Install the requested integrations.
	log.Debugf("Installing integrations")
	if err := i.installRecipes(ctx, m, recipes); err != nil {
		return err
	}

	log.Debugf("Done installing integrations.")

	return nil
}

func (i *RecipeInstaller) recipeFromPath(recipePath string) (*types.OpenInstallationRecipe, error) {
	recipeURL, parseErr := url.Parse(recipePath)
	if parseErr == nil && recipeURL.Scheme != "" {
		f, err := i.recipeFileFetcher.FetchRecipeFile(recipeURL)
		if err != nil {
			return nil, fmt.Errorf("could not fetch file %s: %s", recipePath, err)
		}
		return f, nil
	}

	f, err := i.recipeFileFetcher.LoadRecipeFile(recipePath)
	if err != nil {
		return nil, fmt.Errorf("could not load file %s: %s", recipePath, err)
	}

	return f, nil
}

func (i *RecipeInstaller) fetchWarn(m *types.DiscoveryManifest, recipeName string) *types.OpenInstallationRecipe {
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

func getRecipe(recipeName string, recipes []types.OpenInstallationRecipe) *types.OpenInstallationRecipe {
	for _, r := range recipes {
		if recipeName == r.Name {
			return &r
		}
	}
	return nil
}
