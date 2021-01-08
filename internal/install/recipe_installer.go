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
	infraAgentRecipeName = "infrastructure-agent-installer"
	loggingRecipeName    = "logs-integration"
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

// nolint:gocyclo
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
	fmt.Println()

	// Execute the discovery process if necessary, exiting on failure.
	m, err := i.discoverWithProgress()
	if err != nil {
		return i.fail(err)
	}

	var recipes []types.Recipe
	if i.RecipePathsProvided() {
		// Load the recipes from the provided file names.
		for _, n := range i.RecipePaths {
			log.Debugln(fmt.Sprintf("Attempting to match recipePath %s.", n))
			var recipe *types.Recipe
			recipe, err = i.recipeFromPath(n)
			if err != nil {
				log.Debugln(fmt.Sprintf("Error while building recipe from path, detail:%s.", err))
				return i.fail(err)
			}
			log.Debugln(fmt.Sprintf("Found recipe from path %s.", n))
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
	} else if i.ShouldRunDiscovery() {
		log.Debugln("Ask the recipe service for recommendations.")
		recipes, err = i.fetchRecommendationsWithStatus(m)
		if err != nil {
			log.Debugln(fmt.Sprintf("Error while finding recommendations, detail:%s.", err))
			return i.fail(err)
		}

		if len(recipes) == 0 {
			log.Debugln("No available integrations found.")
		}

		for _, r := range recipes {
			log.Debugf(fmt.Sprintf("Found available integration %s.", r.Name))
		}

		i.reportRecipesAvailable(recipes)
	}

	log.Debugf(fmt.Sprintf("AssumeYes:%t.", i.AssumeYes))
	log.Debugf(fmt.Sprintf("RecipeNames:%s.", i.RecipeNames))
	log.Debugf(fmt.Sprintf("RecipePaths:%s.", i.RecipePaths))
	log.Debugf(fmt.Sprintf("SkipDiscovery:%t.", i.SkipDiscovery))
	log.Debugf(fmt.Sprintf("SkipInfraInstall:%t.", i.SkipInfraInstall))
	log.Debugf(fmt.Sprintf("SkipLoggingInstall:%t.", i.SkipLoggingInstall))
	log.Debugf(fmt.Sprintf("SkipIntegrations:%t.", i.SkipIntegrations))
	log.Debugf(fmt.Sprintf("RecipesProvided:%t.", i.RecipesProvided()))

	// Install the Infrastructure Agent if requested, exiting on failure.
	var entityGUID string
	if i.ShouldInstallInfraAgent() {
		log.Debugf("Installing infrastructure agent...")
		entityGUID, err = i.installInfraAgent(m)
		if err != nil {
			log.Error(i.failMessage(infraAgentRecipeName))
			return i.fail(err)
		}
		log.Debugf("Done installing infrastructure agent.")
	} else {
		log.Debugf("Skipping installing infrastructure agent")
	}

	// Run the logging recipe if requested, exiting on failure.
	if i.ShouldInstallLogging() {
		log.Debugf("Installing logging...")
		if err = i.installLogging(m, recipes); err != nil {
			log.Error(i.failMessage(loggingRecipeName))
			return i.fail(err)
		}
		log.Debugf("Done installing logging.")
	} else {
		log.Debugf("Skipping installing logging")
	}

	recipeErrorsEncountered := false
	// Install integrations if necessary, continuing on failure with warnings.
	if i.ShouldInstallIntegrations() {
		log.Debugf("Installing integrations...")
		if recipeErrorsEncountered, err = i.installRecipesWithPrompts(m, recipes, entityGUID); err != nil {
			return err
		}
		log.Debugf("Done installing integrations.")
	} else {
		log.Debugf("Skipping installing integrations")
	}

	i.reportComplete()

	if recipeErrorsEncountered {
		return errors.New("one or more integrations failed to install, check the install log for more details")
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
	return nil
}

func (i *RecipeInstaller) installRecipesWithPrompts(m *types.DiscoveryManifest, recipes []types.Recipe, entityGUID string) (bool, error) {
	errorsEncountered := false
	log.Debugf("Installing recipes with prompts...")

	for _, r := range recipes {
		log.Debugf(fmt.Sprintf("Installing recipe %s with prompts...", r.Name))
		// The infra and logging install have their own install methods.  In the
		// case where the recommendations come back with either of these recipes,
		// we skip here to avoid duplicate installation.
		if !i.RecipesProvided() {
			if r.Name == infraAgentRecipeName || r.Name == loggingRecipeName {
				log.Debugf(fmt.Sprintf("Skipping recipe %s with prompts, matching either infra agent name %s or logging recipe name %s.", r.Name, infraAgentRecipeName, loggingRecipeName))
				continue
			}
		}

		var ok bool
		var err error

		// Skip prompting the user if the recipe has been asked for directly.
		if i.RecipesProvided() || i.AssumeYes {
			ok = true
		} else {
			log.Debugf("Checking user accepts install...")
			ok, err = i.userAcceptsInstall(r)
			if err != nil {
				log.Debugf(fmt.Sprintf("Done installing recipes with prompts, exception:%s", err))
				return true, err
			}
			log.Debugf(fmt.Sprintf("Done checking user accepts install ok:%t", ok))
		}

		if !ok {
			log.Debugf(fmt.Sprintf("skipping not ok recipe %s.", r.Name))
			i.reportRecipeSkipped(execution.RecipeStatusEvent{
				Recipe:     r,
				EntityGUID: entityGUID,
			})
			continue
		}

		log.Debugf(fmt.Sprintf("Executing and validating with progress for recipe name %s...", r.Name))
		_, err = i.executeAndValidateWithProgress(m, &r)
		if err != nil {
			errorsEncountered = true
			log.Debugf(fmt.Sprintf("Failed while executing and validating with progress for recipe name %s, detail:%s", r.Name, err))
			log.Warn(err)
			log.Warn(i.failMessage(r.Name))
		}
		log.Debugf(fmt.Sprintf("Done executing and validating with progress for recipe name %s.", r.Name))
	}

	log.Debugf(fmt.Sprintf("Done installing recipes with prompts, errorsEncountered:%t", errorsEncountered))
	return errorsEncountered, nil
}

func (i *RecipeInstaller) discoverWithProgress() (*types.DiscoveryManifest, error) {
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

func (i *RecipeInstaller) installLogging(m *types.DiscoveryManifest, recipes []types.Recipe) error {
	log.WithFields(log.Fields{
		"name": loggingRecipeName,
	}).Debug("fetching log recipe for install")

	r, err := i.fetch(m, loggingRecipeName)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"recipe_count": len(recipes),
	}).Debug("filtering log matches")
	logMatches, err := i.fileFilterer.Filter(utils.SignalCtx, recipes)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"possible_matches": len(logMatches),
	}).Debug("filtered log matches")

	var acceptedLogMatches []types.LogMatch
	var ok bool
	for _, match := range logMatches {
		ok, err = i.userAcceptsLogFile(match)
		if err != nil {
			return err
		}

		if ok {
			acceptedLogMatches = append(acceptedLogMatches, match)
		}
	}

	log.WithFields(log.Fields{
		"matches": acceptedLogMatches,
	}).Debug("matches accepted")

	// The struct to approximate the logging configuration file of the Infra Agent.
	type loggingConfig struct {
		Logs []types.LogMatch `yaml:"logs"`
	}

	r.AddVar("DISCOVERED_LOG_FILES", loggingConfig{Logs: acceptedLogMatches})

	_, err = i.executeAndValidateWithProgress(m, r)
	return err
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

	log.WithFields(log.Fields{
		"recipe_count": len(recipes),
	}).Debug("recipes received")

	return recipes, nil
}

func (i *RecipeInstaller) fetchExecuteAndValidate(m *types.DiscoveryManifest, recipeName string) (string, error) {
	log.WithFields(log.Fields{
		"name": recipeName,
	}).Trace("fetching recipe by name")

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

func (i *RecipeInstaller) executeAndValidate(m *types.DiscoveryManifest, r *types.Recipe, vars types.RecipeVars) (string, error) {
	// Execute the recipe steps.
	if err := i.recipeExecutor.Execute(utils.SignalCtx, *m, *r, vars); err != nil {
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
	vars, err := i.recipeExecutor.Prepare(utils.SignalCtx, *m, *r, i.AssumeYes)
	if err != nil {
		return "", fmt.Errorf("could not prepare recipe %s", err)
	}

	i.progressIndicator.Start(fmt.Sprintf("Installing %s...", r.Name))
	defer func() { i.progressIndicator.Stop() }()

	entityGUID, err := i.executeAndValidate(m, r, vars)
	if err != nil {
		i.progressIndicator.Fail()
		return "", fmt.Errorf("could not install recipe %s: %s", r.Name, err)
	}

	i.progressIndicator.Success()
	return entityGUID, nil
}

func (i *RecipeInstaller) userAccepts(msg string) (bool, error) {
	if i.AssumeYes {
		return true, nil
	}

	val, err := i.prompter.PromptYesNo(msg)
	if err != nil {
		return false, err
	}

	return val, nil
}

func (i *RecipeInstaller) userAcceptsLogFile(match types.LogMatch) (bool, error) {
	if i.AssumeYes {
		return true, nil
	}

	msg := fmt.Sprintf("Files have been found at the following pattern: %s\nDo you want to watch them? [Yes/No]", match.File)
	return i.userAccepts(msg)
}

func (i *RecipeInstaller) userAcceptsInstall(r types.Recipe) (bool, error) {
	if i.AssumeYes {
		return true, nil
	}

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
