package install

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

type installContext struct {
	// nolint:structcheck,unused
	interactiveMode     bool
	autoDiscoveryMode   bool
	recipeFriendlyNames []string
}

func install(client *newrelic.NewRelic, ic installContext) error {
	rf := newServiceRecipeFetcher(&client.NerdGraph)
	pf := newRegexProcessFilterer(rf)
	var v recipeValidator = newPollingRecipeValidator(&client.Nrdb)
	var e recipeExecutor = newGoTaskRecipeExecutor()
	var d discoverer = newPSUtilDiscoverer(pf)

	// Execute the discovery process.
	log.Info("Running discovery...")
	m, err := d.discover(utils.SignalCtx)
	if err != nil {
		return err
	}

	var recipes []recipe
	if ic.autoDiscoveryMode {
		// Retrieve the relevant recipes.
		log.Info("Retrieving recipes...")
		recipes, err = rf.fetchRecommendations(utils.SignalCtx, m)
		if err != nil {
			return err
		}
	} else {
		// Search for the relevant recipes.
		for _, n := range ic.recipeFriendlyNames {
			log.Infof("Retrieving recipe %s...", n)

			r, err := rf.fetchRecipe(utils.SignalCtx, m, n)
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
		if err := e.execute(utils.SignalCtx, *m, r); err != nil {
			return fmt.Errorf("encountered an error while executing %s: %s", r.Metadata.Name, err)
		}

		log.Infof("Validating %s...", r.Metadata.Name)

		// Execute the recipe steps.
		ok, err := v.validate(utils.SignalCtx, r)
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
