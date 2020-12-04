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

const (
	infraAgentRecipeName = "Infrastructure Agent Installer"
	loggingRecipeName    = "Logs integration"
)

func (i *recipeInstaller) install() {

	log.Infoln("Welcome to New Relic. Let's install some instrumentation.")
	log.Infoln("Questions? Read more about our installation process at https://docs.newrelic.com/install-newrelic.")

	// Execute the discovery process, exiting on failure.
	m := i.discoverFatal()

	// Run the infra agent recipe, exiting on failure.
	i.installInfraAgentFatal(m)

	// Run the logging recipe, exiting on failure.
	i.installLoggingFatal(m)

	// Retrieve a list of recipes to execute.
	var recipes []recipe
	if i.autoDiscoveryMode {
		// Ask the recipe service for recommendations.
		recipes = i.fetchRecommendationsFatal(m)
	} else {
		// Execute the requested recipes.
		for _, n := range i.recipeFriendlyNames {
			r := i.fetchWarn(m, n)
			recipes = append(recipes, *r)
		}
	}

	// Execute and validate each of the recipes in the collection.
	for _, r := range recipes {
		i.executeAndValidateWarn(m, &r)
	}

	log.Infoln("Success! Your data is available in New Relic.")
	log.Infoln("Go to New Relic to confirm and start exploring your data.")
}

func (i *recipeInstaller) discoverFatal() *discoveryManifest {
	m, err := i.discoverer.discover(utils.SignalCtx)
	if err != nil {
		log.Fatalf("Could not install New Relic.  There was an error discovering system info: %s", err)
	}

	return m
}

func (i *recipeInstaller) installInfraAgentFatal(m *discoveryManifest) {
	i.fetchExecuteAndValidateFatal(m, infraAgentRecipeName)
}

func (i *recipeInstaller) installLoggingFatal(m *discoveryManifest) {
	i.fetchExecuteAndValidateFatal(m, loggingRecipeName)
}

func (i *recipeInstaller) fetchRecommendationsFatal(m *discoveryManifest) []recipe {
	recipes, err := i.recipeFetcher.fetchRecommendations(utils.SignalCtx, m)
	if err != nil {
		log.Fatalf("Could not install New Relic. Error retrieving recipe recommendations: %s", err)
	}

	return recipes
}

func (i *recipeInstaller) fetchExecuteAndValidateFatal(m *discoveryManifest, recipeName string) {
	r := i.fetchFatal(m, recipeName)
	i.executeAndValidateFatal(m, r)
}

func (i *recipeInstaller) fetchWarn(m *discoveryManifest, recipeName string) *recipe {
	r, err := i.recipeFetcher.fetchRecipe(utils.SignalCtx, m, infraAgentRecipeName)
	if err != nil {
		log.Warnf("Could not install %s. Error retrieving recipe: %s", recipeName, err)
		return nil
	}

	if r == nil {
		log.Warnf("Recipe %s not found. Skipping installation.", recipeName)
	}

	return r
}

func (i *recipeInstaller) fetchFatal(m *discoveryManifest, recipeName string) *recipe {
	r, err := i.recipeFetcher.fetchRecipe(utils.SignalCtx, m, infraAgentRecipeName)
	if err != nil {
		log.Fatalf("Could not install %s. Error retrieving recipe: %s", recipeName, err)
	}

	if r == nil {
		log.Fatalf("Recipe %s not found.", recipeName)
	}

	return r
}

func (i *recipeInstaller) executeAndValidate(m *discoveryManifest, r *recipe) (bool, error) {
	// Execute the recipe steps.
	log.Infof("Installing %s...\n", r.Metadata.Name)
	if err := i.recipeExecutor.execute(utils.SignalCtx, *m, *r); err != nil {
		return false, fmt.Errorf("encountered an error while executing %s: %s", r.Metadata.Name, err)
	}
	log.Infof("Installing %s...success\n", r.Metadata.Name)

	log.Info("Listening for data...")
	ok, err := i.recipeValidator.validate(utils.SignalCtx, *r)
	if err != nil {
		return false, fmt.Errorf("encountered an error while validating receipt of data for %s: %s", r.Metadata.Name, err)
	}

	if !ok {
		log.Infoln("failed.")
		return false, nil
	}

	log.Infoln("success.")
	return true, nil
}

func (i *recipeInstaller) executeAndValidateFatal(m *discoveryManifest, r *recipe) {
	ok, err := i.executeAndValidate(m, r)
	if err != nil {
		log.Fatalf("Could not install %s: %s", r.Metadata.Name, err)
	}

	if !ok {
		log.Fatalf("Could not detect data from %s.", r.Metadata.Name)
	}
}

func (i *recipeInstaller) executeAndValidateWarn(m *discoveryManifest, r *recipe) {
	ok, err := i.executeAndValidate(m, r)
	if err != nil {
		log.Warnf("Could not install %s: %s", r.Metadata.Name, err)
	}

	if !ok {
		log.Warnf("Could not detect data from %s.", r.Metadata.Name)
	}
}
