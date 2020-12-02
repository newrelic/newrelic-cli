package install

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

type recipeInstaller struct {
	installContext
	discoverer      discoverer
	recipeFetcher   recipeFetcher
	recipeExecutor  recipeExecutor
	recipeValidator recipeValidator
}

func newRecipeInstaller(
	ic installContext,
	d discoverer,
	f recipeFetcher,
	e recipeExecutor,
	v recipeValidator,
) *recipeInstaller {
	i := recipeInstaller{
		discoverer:      d,
		recipeFetcher:   f,
		recipeExecutor:  e,
		recipeValidator: v,
	}

	i.autoDiscoveryMode = ic.autoDiscoveryMode
	i.interactiveMode = ic.interactiveMode
	i.recipeFriendlyNames = ic.recipeFriendlyNames

	return &i
}

type installContext struct {
	interactiveMode     bool
	autoDiscoveryMode   bool
	recipeFriendlyNames []string
}

func (i *recipeInstaller) install() error {
	// Execute the discovery process.
	log.Info("Running discovery...")
	m, err := i.discoverer.discover(utils.SignalCtx)
	if err != nil {
		return err
	}

	var recipes []recipe
	if i.autoDiscoveryMode {
		// Retrieve the relevant recipes.
		log.Info("Retrieving recipes...")
		recipes, err = i.recipeFetcher.fetchRecommendations(utils.SignalCtx, m)
		if err != nil {
			return err
		}
	} else {
		// Search for the relevant recipes.
		for _, n := range i.recipeFriendlyNames {
			log.Infof("Retrieving recipe %s...", n)

			r, err := i.recipeFetcher.fetchRecipe(utils.SignalCtx, m, n)
			if err != nil {
				return fmt.Errorf("error retrieving recipe %s: %s", n, err)
			}

			if r == nil {
				return fmt.Errorf("recipe %s not found", n)
			}
		}
	}

	// Iterate through the recipe collection, executing recipe steps and
	// validating execution for each recipe.
	for _, r := range recipes {
		log.Infof("Executing %s...", r.Metadata.Name)

		// Execute the recipe steps.
		if err := i.recipeExecutor.execute(utils.SignalCtx, *m, r); err != nil {
			return fmt.Errorf("encountered an error while executing %s: %s", r.Metadata.Name, err)
		}

		log.Infof("Validating %s...", r.Metadata.Name)

		// Execute the recipe steps.
		ok, err := i.recipeValidator.validate(utils.SignalCtx, r)
		if err != nil {
			return fmt.Errorf("encountered an error while validating receipt of data for %s: %s", r.Metadata.Name, err)
		}

		if !ok {
			log.Warnf("Data could not be found for %s.", r.Metadata.Name)
			continue
		}

		log.Infof("Data is being sent to New Relic for recipe %s.", r.Metadata.Name)
	}

	return nil
}
