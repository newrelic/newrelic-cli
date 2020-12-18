package install

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type recipeInstaller struct {
	installContext
	discoverer        discoverer
	fileFilterer      fileFilterer
	recipeFetcher     recipeFetcher
	recipeExecutor    recipeExecutor
	recipeValidator   recipeValidator
	recipeFileFetcher recipeFileFetcher
	statusReporter    executionStatusReporter
}

func newRecipeInstaller(
	ic installContext,
	d discoverer,
	l fileFilterer,
	f recipeFetcher,
	e recipeExecutor,
	v recipeValidator,
	ff recipeFileFetcher,
	er executionStatusReporter,
) *recipeInstaller {
	i := recipeInstaller{
		discoverer:        d,
		fileFilterer:      l,
		recipeFetcher:     f,
		recipeExecutor:    e,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		statusReporter:    er,
	}

	i.recipePaths = ic.recipePaths
	i.recipeNames = ic.recipeNames
	i.skipDiscovery = ic.skipDiscovery
	i.skipInfraInstall = ic.skipInfraInstall
	i.skipIntegrations = ic.skipIntegrations
	i.skipLoggingInstall = ic.skipLoggingInstall

	return &i
}

const (
	infraAgentRecipeName = "Infrastructure Agent Installer"
	loggingRecipeName    = "Logs integration"
	checkMark            = "\u2705"
	boom                 = "\u1F4A5"
)

func (i *recipeInstaller) install() error {
	fmt.Printf(`
	Welcome to New Relic. Let's install some instrumentation.

	Questions? Read more about our installation process at
	https://docs.newrelic.com/

	`)

	// Execute the discovery process if necessary, exiting on failure.
	var m *discoveryManifest
	var err error
	if i.ShouldRunDiscovery() {
		m, err = i.discover()
		if err != nil {
			return i.fail(err)
		}
	}

	var recipes []recipe
	if i.RecipePathsProvided() {
		// Load the recipes from the provided file names.
		for _, n := range i.recipePaths {
			var recipe *recipe
			recipe, err = i.recipeFromPath(n)
			if err != nil {
				return i.fail(err)
			}

			recipes = append(recipes, *recipe)
		}
	} else if i.RecipeNamesProvided() {
		// Fetch the provided recipes from the recipe service.
		for _, n := range i.recipeNames {
			r := i.fetchWarn(m, n)
			recipes = append(recipes, *r)
		}
	} else {
		// Ask the recipe service for recommendations.
		log.Debugln("Fetching recommendations.")
		recipes, err = i.fetchRecommendations(m)
		if err != nil {
			return i.fail(err)
		}

		if len(recipes) == 0 {
			log.Debugln("No available integrations found.")
		}

		for _, r := range recipes {
			log.Debugf("Found available integration %s.", r.Name)
		}

		i.reportRecipesAvailable(recipes)
	}

	// Install the Infrastructure Agent if requested, exiting on failure.
	if i.ShouldInstallInfraAgent() {
		if err := i.installInfraAgent(m); err != nil {
			return i.fail(err)
		}
	}

	// Run the logging recipe if requested, exiting on failure.
	if i.ShouldInstallLogging() {
		if err := i.installLogging(m, recipes); err != nil {
			return i.fail(err)
		}
	}

	// Install integrations if necessary, continuing on failure with warnings.
	if i.ShouldInstallIntegrations() {
		for _, r := range recipes {
			err := i.executeAndValidateWithStatus(m, &r)
			if err != nil {
				log.Warn(err)
			}
		}
	}

	profile := credentials.DefaultProfile()
	fmt.Printf(`
	Success! Your data is available in New Relic.

	Go to New Relic to confirm and start exploring your data.
	https://one.newrelic.com/launcher/nrai.launcher?platform[accountId]=%d
	`, profile.AccountID)

	fmt.Println()

	i.reportComplete()
	return nil
}

func (i *recipeInstaller) discover() (*discoveryManifest, error) {
	s := newSpinner()
	s.Suffix = " Discovering system information..."

	s.Start()
	defer func() {
		s.Stop()
		fmt.Println(s.Suffix)
	}()

	m, err := i.discoverer.discover(utils.SignalCtx)
	if err != nil {
		s.FinalMSG = boom
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

	s.FinalMSG = checkMark

	return m, nil
}

func (i *recipeInstaller) recipeFromPath(recipePath string) (*recipe, error) {
	recipeURL, parseErr := url.Parse(recipePath)
	if parseErr == nil && recipeURL.Scheme != "" {
		f, err := i.recipeFileFetcher.fetchRecipeFile(recipeURL)
		if err != nil {
			return nil, fmt.Errorf("could not fetch file %s: %s", recipePath, err)
		}
		return finalizeRecipe(f)
	}

	f, err := i.recipeFileFetcher.loadRecipeFile(recipePath)
	if err != nil {
		return nil, fmt.Errorf("could not load file %s: %s", recipePath, err)
	}
	return finalizeRecipe(f)
}

func finalizeRecipe(f *recipeFile) (*recipe, error) {
	r, err := f.ToRecipe()
	if err != nil {
		return nil, fmt.Errorf("could not finalize recipe %s: %s", f.Name, err)
	}
	return r, nil
}

func (i *recipeInstaller) installInfraAgent(m *discoveryManifest) error {
	return i.fetchExecuteAndValidate(m, infraAgentRecipeName)
}

func (i *recipeInstaller) installLogging(m *discoveryManifest, recipes []recipe) error {
	r, err := i.fetch(m, loggingRecipeName)
	if err != nil {
		return err
	}

	logMatches, err := i.fileFilterer.filter(utils.SignalCtx, recipes)
	if err != nil {
		return err
	}

	var acceptedLogMatches []logMatch
	for _, match := range logMatches {
		if userAcceptLogFile(match) {
			acceptedLogMatches = append(acceptedLogMatches, match)
		}
	}

	// The struct to approximate the logging configuration file of the Infra Agent.
	type loggingConfig struct {
		Logs []logMatch `yaml:"logs"`
	}

	r.AddVar("DISCOVERED_LOG_FILES", loggingConfig{Logs: acceptedLogMatches})

	return i.executeAndValidateWithStatus(m, r)
}

func (i *recipeInstaller) fetchRecommendations(m *discoveryManifest) ([]recipe, error) {
	s := newSpinner()
	s.Suffix = " Fetching recommended recipes..."

	s.Start()
	defer func() {
		s.Stop()
		fmt.Println(s.Suffix)
	}()

	recipes, err := i.recipeFetcher.fetchRecommendations(utils.SignalCtx, m)
	if err != nil {
		s.FinalMSG = boom
		return nil, fmt.Errorf("error retrieving recipe recommendations: %s", err)
	}

	s.FinalMSG = checkMark

	return recipes, nil
}

func (i *recipeInstaller) fetchExecuteAndValidate(m *discoveryManifest, recipeName string) error {
	r, err := i.fetch(m, recipeName)
	if err != nil {
		return err
	}

	if err := i.executeAndValidateWithStatus(m, r); err != nil {
		return err
	}

	return nil
}

func (i *recipeInstaller) fetchWarn(m *discoveryManifest, recipeName string) *recipe {
	r, err := i.recipeFetcher.fetchRecipe(utils.SignalCtx, m, recipeName)
	if err != nil {
		log.Warnf("Could not install %s. Error retrieving recipe: %s", recipeName, err)
		return nil
	}

	if r == nil {
		log.Warnf("Recipe %s not found. Skipping installation.", recipeName)
	}

	return r
}

func (i *recipeInstaller) fetch(m *discoveryManifest, recipeName string) (*recipe, error) {
	r, err := i.recipeFetcher.fetchRecipe(utils.SignalCtx, m, recipeName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recipe %s: %s", recipeName, err)
	}

	if r == nil {
		return nil, fmt.Errorf("recipe %s not found", recipeName)
	}

	return r, nil
}

func (i *recipeInstaller) executeAndValidate(m *discoveryManifest, r *recipe) error {
	// Execute the recipe steps.
	if err := i.recipeExecutor.execute(utils.SignalCtx, *m, *r); err != nil {
		msg := fmt.Sprintf("encountered an error while executing %s: %s", r.Name, err)
		i.reportRecipeFailed(recipeStatusEvent{*r, msg, ""})
		return errors.New(msg)
	}

	if r.ValidationNRQL != "" {
		entityGUID, err := i.recipeValidator.validate(utils.SignalCtx, *m, *r)
		if err != nil {
			msg := fmt.Sprintf("encountered an error while validating receipt of data for %s: %s", r.Name, err)
			i.reportRecipeFailed(recipeStatusEvent{*r, msg, ""})
			return errors.New(msg)
		}

		i.reportRecipeInstalled(recipeStatusEvent{*r, "", entityGUID})
	} else {
		log.Debugf("Skipping validation due to missing validation query.")
	}

	return nil
}

func (i *recipeInstaller) reportRecipesAvailable(recipes []recipe) {
	if err := i.statusReporter.reportRecipesAvailable(recipes); err != nil {
		log.Errorf("Could not report recipe execution status: %s", err)
	}
}

func (i *recipeInstaller) reportRecipeInstalled(e recipeStatusEvent) {
	if err := i.statusReporter.reportRecipeInstalled(e); err != nil {
		log.Errorf("Error writing recipe status for recipe %s: %s", e.recipe.Name, err)
	}
}

func (i *recipeInstaller) reportRecipeFailed(e recipeStatusEvent) {
	if err := i.statusReporter.reportRecipeFailed(e); err != nil {
		log.Errorf("Error writing recipe status for recipe %s: %s", e.recipe.Name, err)
	}
}

func (i *recipeInstaller) reportComplete() {
	if err := i.statusReporter.reportComplete(); err != nil {
		log.Errorf("Error writing execution status: %s", err)
	}
}

func (i *recipeInstaller) executeAndValidateWithStatus(m *discoveryManifest, r *recipe) error {
	s := newSpinner()
	s.Suffix = fmt.Sprintf(" Installing %s...", r.Name)

	s.Start()
	defer func() {
		s.Stop()
		fmt.Println(s.Suffix)
	}()

	err := i.executeAndValidate(m, r)
	if err != nil {
		s.FinalMSG = boom
		return fmt.Errorf("could not install %s: %s", r.Name, err)
	}

	s.FinalMSG = checkMark
	return nil
}

func userAcceptLogFile(match logMatch) bool {
	msg := fmt.Sprintf("Files have been found at the following pattern: %s\nDo you want to watch them? [Yes/No]", match.File)

	prompt := promptui.Select{
		Label: msg,
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Errorf("prompt failed: %s", err)
		return false
	}

	return result == "Yes"
}

func (i *recipeInstaller) fail(err error) error {
	log.Error(err)

	if err = i.statusReporter.reportComplete(); err != nil {
		log.Error(err)
	}

	return err
}

func newSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[14], 100*time.Millisecond)
}
