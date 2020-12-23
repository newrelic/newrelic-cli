package install

import (
	"errors"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

const (
	infraAgentRecipeName = "Infrastructure Agent Installer"
	loggingRecipeName    = "Logs integration"
)

type RecipeInstaller struct {
	InstallerContext
	discoverer        discovery.Discoverer
	fileFilterer      discovery.FileFilterer
	recipeFetcher     recipes.RecipeFetcher
	recipeExecutor    execution.RecipeExecutor
	recipeValidator   validation.RecipeValidator
	recipeFileFetcher recipes.RecipeFileFetcher
	statusReporter    execution.StatusReporter
	prompter          ux.Prompter
	progressIndicator ux.ProgressIndicator
}

func NewRecipeInstaller(ic InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstaller {
	rf := recipes.NewServiceRecipeFetcher(&nrClient.NerdGraph)
	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	er := execution.NewNerdStorageStatusReporter(&nrClient.NerdStorage)
	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	p := ux.NewPromptUIPrompter()
	s := ux.NewSpinner()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		statusReporter:    er,
		prompter:          p,
		progressIndicator: s,
	}

	i.InstallerContext = ic

	return &i
}

func (i *RecipeInstaller) Install() error {
	fmt.Printf(`
	 _   _                 ____      _ _
	| \ | | _____      __ |  _ \ ___| (_) ___
	|  \| |/ _ \ \ /\ / / | |_) / _ | | |/ __|
	| |\  |  __/\ V  V /  |  _ |  __| | | (__
	|_| \_|\___| \_/\_/   |_| \_\___|_|_|\___|

	Welcome to New Relic. Let's install some instrumentation.

	Questions? Read more about our installation process at
	https://docs.newrelic.com/

	`)

	// Execute the discovery process if necessary, exiting on failure.
	var m *types.DiscoveryManifest
	var err error
	if i.ShouldRunDiscovery() {
		m, err = i.discoverWithStatus()
		if err != nil {
			return i.fail(err)
		}
	}

	var recipes []types.Recipe
	if i.RecipePathsProvided() {
		// Load the recipes from the provided file names.
		for _, n := range i.RecipePaths {
			var recipe *types.Recipe
			recipe, err = i.recipeFromPath(n)
			if err != nil {
				return i.fail(err)
			}

			recipes = append(recipes, *recipe)
		}
	} else if i.RecipeNamesProvided() {
		// Fetch the provided recipes from the recipe service.
		for _, n := range i.RecipeNames {
			r := i.fetchWarn(m, n)
			recipes = append(recipes, *r)
		}
	} else {
		// Ask the recipe service for recommendations.
		recipes, err = i.fetchRecommendationsWithStatus(m)
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
	var entityGUID string
	if i.ShouldInstallInfraAgent() {
		entityGUID, err = i.installInfraAgent(m)
		if err != nil {
			log.Error(i.failMessage(infraAgentRecipeName))
			return i.fail(err)
		}

	}

	// Run the logging recipe if requested, exiting on failure.
	if i.ShouldInstallLogging() {
		_, err := i.installLogging(m, recipes)
		if err != nil {
			log.Error(i.failMessage(loggingRecipeName))
			return i.fail(err)
		}
	}

	// Install integrations if necessary, continuing on failure with warnings.
	if i.ShouldInstallIntegrations() {
		for _, r := range recipes {
			if !i.userAcceptsInstall(r) {
				log.Infof("Skipping %s.", r.Name)
				i.reportRecipeSkipped(execution.RecipeStatusEvent{
					Recipe:     r,
					EntityGUID: entityGUID,
				})
				continue
			}

			_, err := i.executeAndValidateWithProgress(m, &r)
			if err != nil {
				log.Warn(err)
				log.Warn(i.failMessage(r.Name))
			}
		}
	}

	msg := `
	Success! Your data is available in New Relic.

	Go to New Relic to confirm and start exploring your data.`

	profile := credentials.DefaultProfile()
	if profile != nil {
		msg += fmt.Sprintf(`
		https://one.newrelic.com/launcher/nrai.launcher?platform[accountId]=%d`, profile.AccountID)
	}

	fmt.Println(msg)

	i.reportComplete()
	return nil
}

func (i *RecipeInstaller) discoverWithStatus() (*types.DiscoveryManifest, error) {
	i.progressIndicator.Start("Discovering system information...")
	defer func() {
		i.progressIndicator.Stop()
	}()

	m, err := i.discoverer.Discover(utils.SignalCtx)
	if err != nil {
		i.progressIndicator.Fail()
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

	i.progressIndicator.Success()

	return m, nil
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

func (i *RecipeInstaller) installInfraAgent(m *types.DiscoveryManifest) (string, error) {
	return i.fetchExecuteAndValidate(m, infraAgentRecipeName)
}

func (i *RecipeInstaller) installLogging(m *types.DiscoveryManifest, recipes []types.Recipe) (string, error) {
	r, err := i.fetch(m, loggingRecipeName)
	if err != nil {
		return "", err
	}

	logMatches, err := i.fileFilterer.Filter(utils.SignalCtx, recipes)
	if err != nil {
		return "", err
	}

	var acceptedLogMatches []types.LogMatch
	for _, match := range logMatches {
		if i.userAcceptsLogFile(match) {
			acceptedLogMatches = append(acceptedLogMatches, match)
		}
	}

	// The struct to approximate the logging configuration file of the Infra Agent.
	type loggingConfig struct {
		Logs []types.LogMatch `yaml:"logs"`
	}

	r.AddVar("DISCOVERED_LOG_FILES", loggingConfig{Logs: acceptedLogMatches})

	return i.executeAndValidateWithProgress(m, r)
}

func (i *RecipeInstaller) fetchRecommendationsWithStatus(m *types.DiscoveryManifest) ([]types.Recipe, error) {
	i.progressIndicator.Start("Fetching recommended recipes...")
	defer func() {
		i.progressIndicator.Stop()
	}()

	recipes, err := i.recipeFetcher.FetchRecommendations(utils.SignalCtx, m)
	if err != nil {
		i.progressIndicator.Fail()
		return nil, fmt.Errorf("error retrieving recipe recommendations: %s", err)
	}

	i.progressIndicator.Success()

	return recipes, nil
}

func (i *RecipeInstaller) fetchExecuteAndValidate(m *types.DiscoveryManifest, recipeName string) (string, error) {
	r, err := i.fetch(m, recipeName)
	if err != nil {
		return "", err
	}

	entityGUID, err := i.executeAndValidateWithProgress(m, r)
	if err != nil {
		return "", err
	}

	return entityGUID, nil
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

func (i *RecipeInstaller) fetch(m *types.DiscoveryManifest, recipeName string) (*types.Recipe, error) {
	r, err := i.recipeFetcher.FetchRecipe(utils.SignalCtx, m, recipeName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recipe %s: %s", recipeName, err)
	}

	if r == nil {
		return nil, fmt.Errorf("recipe %s not found", recipeName)
	}

	return r, nil
}

func (i *RecipeInstaller) executeAndValidate(m *types.DiscoveryManifest, r *types.Recipe) (string, error) {
	// Execute the recipe steps.
	if err := i.recipeExecutor.Execute(utils.SignalCtx, *m, *r); err != nil {
		msg := fmt.Sprintf("encountered an error while executing %s: %s", r.Name, err)
		i.reportRecipeFailed(execution.RecipeStatusEvent{
			Recipe: *r,
			Msg:    msg,
		})
		return "", errors.New(msg)
	}

	var entityGUID string
	var err error
	if r.ValidationNRQL != "" {
		entityGUID, err = i.recipeValidator.Validate(utils.SignalCtx, *m, *r)
		if err != nil {
			msg := fmt.Sprintf("encountered an error while validating receipt of data for %s: %s", r.Name, err)
			i.reportRecipeFailed(execution.RecipeStatusEvent{
				Recipe: *r,
				Msg:    msg,
			})
			return "", errors.New(msg)
		}

		i.reportRecipeInstalled(execution.RecipeStatusEvent{
			Recipe:     *r,
			EntityGUID: entityGUID,
		})
	} else {
		log.Debugf("Skipping validation due to missing validation query.")
	}

	return entityGUID, nil
}

func (i *RecipeInstaller) reportRecipesAvailable(recipes []types.Recipe) {
	if err := i.statusReporter.ReportRecipesAvailable(recipes); err != nil {
		log.Errorf("Could not report recipe execution status: %s", err)
	}
}

func (i *RecipeInstaller) reportRecipeInstalled(e execution.RecipeStatusEvent) {
	if err := i.statusReporter.ReportRecipeInstalled(e); err != nil {
		log.Errorf("Error writing recipe status for recipe %s: %s", e.Recipe.Name, err)
	}
}

func (i *RecipeInstaller) reportRecipeFailed(e execution.RecipeStatusEvent) {
	if err := i.statusReporter.ReportRecipeFailed(e); err != nil {
		log.Errorf("Error writing recipe status for recipe %s: %s", e.Recipe.Name, err)
	}
}

func (i *RecipeInstaller) reportRecipeSkipped(e execution.RecipeStatusEvent) {
	if err := i.statusReporter.ReportRecipeSkipped(e); err != nil {
		log.Errorf("Error writing recipe status for recipe %s: %s", e.Recipe.Name, err)
	}
}

func (i *RecipeInstaller) reportComplete() {
	if err := i.statusReporter.ReportComplete(); err != nil {
		log.Errorf("Error writing execution status: %s", err)
	}
}

func (i *RecipeInstaller) executeAndValidateWithProgress(m *types.DiscoveryManifest, r *types.Recipe) (string, error) {
	i.progressIndicator.Start(fmt.Sprintf("Installing %s...", r.Name))
	defer func() { i.progressIndicator.Stop() }()

	entityGUID, err := i.executeAndValidate(m, r)
	if err != nil {
		i.progressIndicator.Fail()
		return "", fmt.Errorf("could not install %s: %s", r.Name, err)
	}

	i.progressIndicator.Success()
	return entityGUID, nil
}

func (i *RecipeInstaller) userAccepts(msg string) bool {
	val, err := i.prompter.PromptYesNo(msg)
	if err != nil {
		log.Error(err)
	}

	return val
}

func (i *RecipeInstaller) userAcceptsLogFile(match types.LogMatch) bool {
	msg := fmt.Sprintf("Files have been found at the following pattern: %s\nDo you want to watch them? [Yes/No]", match.File)
	return i.userAccepts(msg)
}

func (i *RecipeInstaller) userAcceptsInstall(r types.Recipe) bool {
	msg := fmt.Sprintf("Would you like to enable %s?", r.Name)
	return i.userAccepts(msg)
}

func (i *RecipeInstaller) fail(err error) error {
	i.reportComplete()
	return err
}

func (i *RecipeInstaller) failMessage(componentName string) error {

	u, _ := url.Parse("https://docs.newrelic.com/search#")
	q := u.Query()
	q.Set("query", componentName)
	u.RawQuery = q.Encode()

	searchURL := u.String()

	return fmt.Errorf("execution of %s failed, please see the following link for clues on how to resolve the issue: %s", componentName, searchURL)
}
