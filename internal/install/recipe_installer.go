package install

import (
	"errors"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

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
	status            *execution.StatusRollup
	prompter          ux.Prompter
	progressIndicator ux.ProgressIndicator
}

func NewRecipeInstaller(ic InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstaller {
	rf := recipes.NewServiceRecipeFetcher(&nrClient.NerdGraph)
	pf := discovery.NewRegexProcessFilterer(rf)
	ff := recipes.NewRecipeFileFetcher()
	ers := []execution.StatusReporter{
		execution.NewNerdStorageStatusReporter(&nrClient.NerdStorage),
		execution.NewTerminalStatusReporter(),
	}
	statusRollup := execution.NewStatusRollup(ers)

	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	p := ux.NewPromptUIPrompter()
	// s := ux.NewSpinner()
	s := ux.NewPlainProgress()

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		recipeFetcher:     rf,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
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

	log.Tracef("InstallerContext: %+v", i.InstallerContext)
	log.WithFields(log.Fields{
		"ShouldRunDiscovery":        i.ShouldRunDiscovery(),
		"ShouldInstallInfraAgent":   i.ShouldInstallInfraAgent(),
		"ShouldInstallLogging":      i.ShouldInstallLogging(),
		"ShouldInstallIntegrations": i.ShouldInstallIntegrations(),
		"RecipesProvided":           i.RecipesProvided(),
		"RecipePathsProvided":       i.RecipePathsProvided(),
		"RecipeNamesProvided":       i.RecipeNamesProvided(),
	}).Debug("context summary")

	// Execute the discovery process, exiting on failure.
	m, err := i.discover()
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
	} else if i.ShouldRunDiscovery() {
		recipes, err = i.fetchRecommendations(m)
		if err != nil {
			log.Debugf("error fetching recommendations: %s", err)
			return i.fail(err)
		}

		if len(recipes) == 0 {
			log.Debug("no available integrations found")
		}

		if i.SkipIntegrations {
			for _, r := range recipes {
				i.status.ReportRecipeSkipped(execution.RecipeStatusEvent{Recipe: r})
			}
		}
	}

	infraAgentRecipe, err := i.fetchRecipeAndReport(m, infraAgentRecipeName)
	if err != nil {
		return err
	}

	loggingRecipe, err := i.fetchRecipeAndReport(m, loggingRecipeName)
	if err != nil {
		return err
	}

	// Include Logging and Infra recipes in the report
	recipesForReport := []types.Recipe{}

	if i.ShouldInstallInfraAgent() {
		if infraAgentRecipe != nil {
			recipesForReport = append(recipesForReport, *infraAgentRecipe)
		}
	}

	if i.ShouldInstallLogging() {
		if loggingRecipe != nil {
			recipesForReport = append(recipesForReport, *loggingRecipe)
		}
	}

	if !i.SkipIntegrations {
		recipesForReport = append(recipesForReport, recipes...)
	}

	// Report discovered recipes as available
	i.status.ReportRecipesAvailable(recipesForReport)

	// If guided installation, prompt user to determine what to skip and what to include
	if !i.RecipesProvided() {
		var filteredRecipes []types.Recipe
		filteredRecipes, err = i.filterSkippedRecipes(recipesForReport)
		if err != nil {
			return i.fail(err)
		}

		recipes = filteredRecipes
	}

	var entityGUID string

	if i.ShouldInstallInfraAgent() {
		log.Debugf("Installing infrastructure agent")
		entityGUID, err = i.executeAndValidateWithProgress(m, infraAgentRecipe)
		if err != nil {
			log.Error(i.failMessage(infraAgentRecipeName))
			return i.fail(err)
		}
		log.Debugf("Done installing infrastructure agent.")
	}

	if i.ShouldInstallLogging() {
		log.Debugf("Installing logging")
		if err = i.installLogging(m, loggingRecipe, recipes); err != nil {
			log.Error(i.failMessage(loggingRecipeName))
			return i.fail(err)
		}
		log.Debugf("Done installing logging.")
	}

	// Install integrations if necessary, continuing on failure with warnings.
	if i.ShouldInstallIntegrations() {
		log.Debugf("Installing integrations")
		if err = i.installRecipes(m, recipes, entityGUID); err != nil {
			return err
		}
		log.Debugf("Done installing integrations.")
	}

	i.status.ReportComplete()

	return nil
}

func (i *RecipeInstaller) installRecipes(m *types.DiscoveryManifest, recipes []types.Recipe, entityGUID string) error {
	log.WithFields(log.Fields{
		"recipe_count": len(recipes),
	}).Debug("installing recipes")

	isHostTarget := func(recipe types.Recipe) bool {
		for _, target := range recipe.InstallTargets {
			if target.Type != types.OpenInstallationTargetTypeTypes.HOST {
				return false
			}
		}

		return true
	}

	for _, r := range recipes {
		// The infra and logging install have their own install methods.  In the
		// case where the recommendations come back with either of these recipes,
		// we skip here to avoid duplicate installation.
		if !i.RecipesProvided() {
			if r.Name == infraAgentRecipeName || r.Name == loggingRecipeName {
				log.WithFields(log.Fields{
					"name": r.Name,
				}).Debug("skipping special recipe")

				continue
			}
		}

		// Skip non-HOST recipes, but report status of them.
		if !isHostTarget(r) {
			log.WithFields(log.Fields{
				"name": r.Name,
			}).Debug("skipping non-HOST recipe")

			i.status.ReportRecipeRecommended(execution.RecipeStatusEvent{
				Recipe:     r,
				EntityGUID: entityGUID,
			})
			continue
		}

		var err error

		log.WithFields(log.Fields{
			"name": r.Name,
		}).Debug("installing recipe")

		_, err = i.executeAndValidateWithProgress(m, &r)
		if err != nil {
			if serr, ok := err.(*types.ErrInterrupt); ok {
				return serr
			}

			log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", r.Name, err)
			log.Warn(err)
			log.Warn(i.failMessage(r.Name))

			if len(recipes) == 1 {
				return err
			}
		}
		log.Debugf("Done executing and validating with progress for recipe name %s.", r.Name)
	}

	return nil
}

func (i *RecipeInstaller) discover() (*types.DiscoveryManifest, error) {
	log.Debug("discovering system information")

	m, err := i.discoverer.Discover(utils.SignalCtx)
	if err != nil {
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

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

func (i *RecipeInstaller) fetchRecipeAndReport(m *types.DiscoveryManifest, recipeName string) (*types.Recipe, error) {
	log.WithFields(log.Fields{
		"name": recipeName,
	}).Debug("fetching recipe for install")

	r, err := i.fetch(m, recipeName)
	if err != nil {
		return nil, err
	}

	switch recipeName {
	case infraAgentRecipeName:
		if i.SkipInfraInstall {
			i.status.ReportRecipeSkipped(execution.RecipeStatusEvent{Recipe: *r})
		} else {
			i.status.ReportRecipeAvailable(*r)
		}
	case loggingRecipeName:
		if i.SkipLoggingInstall {
			i.status.ReportRecipeSkipped(execution.RecipeStatusEvent{Recipe: *r})
		} else {
			i.status.ReportRecipeAvailable(*r)
		}
	default:
		i.status.ReportRecipeAvailable(*r)
	}

	return r, nil
}

func (i *RecipeInstaller) installLogging(m *types.DiscoveryManifest, r *types.Recipe, recipes []types.Recipe) error {
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

func (i *RecipeInstaller) fetchRecommendations(m *types.DiscoveryManifest) ([]types.Recipe, error) {
	log.Debug("fetching recommended recipes")

	recipes, err := i.recipeFetcher.FetchRecommendations(utils.SignalCtx, m)
	if err != nil {
		return nil, fmt.Errorf("error retrieving recipe recommendations: %s", err)
	}

	names := []string{}
	if log.IsLevelEnabled(log.DebugLevel) {
		for _, r := range recipes {
			names = append(names, r.Name)
		}
	}

	log.WithFields(log.Fields{
		"names":        names,
		"recipe_count": len(recipes),
	}).Debug("recipes received")

	return recipes, nil
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
		if serr, ok := err.(*types.ErrInterrupt); ok {
			return "", serr
		}

		msg := fmt.Sprintf("encountered an error while executing %s: %s", r.Name, err)
		i.status.ReportRecipeFailed(execution.RecipeStatusEvent{
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
			i.status.ReportRecipeFailed(execution.RecipeStatusEvent{
				Recipe: *r,
				Msg:    msg,
			})
			return "", errors.New(msg)
		}

		i.status.ReportRecipeInstalled(execution.RecipeStatusEvent{
			Recipe:     *r,
			EntityGUID: entityGUID,
		})
	} else {
		log.Debugf("Skipping validation due to missing validation query.")
	}

	return entityGUID, nil
}

func (i *RecipeInstaller) executeAndValidateWithProgress(m *types.DiscoveryManifest, r *types.Recipe) (string, error) {
	if r.PreInstallMessage() != "" {
		fmt.Println(r.PreInstallMessage())
	}

	vars, err := i.recipeExecutor.Prepare(utils.SignalCtx, *m, *r, i.AssumeYes)
	if err != nil {
		return "", err
	}

	i.progressIndicator.Start(fmt.Sprintf("Installing %s", r.Name))
	defer func() { i.progressIndicator.Stop() }()
	i.status.ReportRecipeInstalling(execution.RecipeStatusEvent{Recipe: *r})

	entityGUID, err := i.executeAndValidate(m, r, vars)
	if err != nil {
		i.progressIndicator.Fail()
		return "", err
	}

	if r.PostInstallMessage() != "" {
		fmt.Println(r.PostInstallMessage())
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
	msg := fmt.Sprintf("Files have been found at the following pattern: %s Do you want to watch them?", match.File)
	return i.userAccepts(msg)
}

func (i *RecipeInstaller) fail(err error) error {
	i.status.ReportComplete()
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

func (i *RecipeInstaller) recipeInRecipes(recipe types.Recipe, recipes []types.Recipe) bool {
	for _, r := range recipes {
		if recipe.Name == r.Name {
			return true
		}
	}

	return false
}

func (i *RecipeInstaller) filterSkippedRecipes(recipes []types.Recipe) ([]types.Recipe, error) {
	installRecommendations, err := i.userAccepts("Would you like to install the recommended instrumentation?")
	if err != nil {
		return nil, err
	}

	var filteredRecipes []types.Recipe
	if !installRecommendations {
		// Start off skipping everything
		i.SkipInfraInstall = true
		i.SkipLoggingInstall = true
		i.SkipIntegrations = true

		reportedDisplayNames := []string{}
		for _, r := range recipes {
			reportedDisplayNames = append(reportedDisplayNames, r.DisplayName)
		}

		selectedRecipeNames, promptErr := i.prompter.MultiSelect("Please choose what instrumentation you would like to install:", reportedDisplayNames)
		if promptErr != nil {
			return nil, promptErr
		}

		// Find the intersection of the selected recipes by DisplayName and the
		// recipesForReport, which is the infra and logging recipe + the OHIs.
		for _, selectedRecipeName := range selectedRecipeNames {
			for _, r := range recipes {
				if r.DisplayName == selectedRecipeName {
					// Disable skip if logging was selected
					if r.Name == loggingRecipeName {
						i.SkipLoggingInstall = false
					}

					// Disable skip if infra was selected
					if r.Name == infraAgentRecipeName {
						i.SkipInfraInstall = false
					}

					// Disable skip if what we have is neither logging no infra
					if r.Name != infraAgentRecipeName && r.Name != loggingRecipeName {
						i.SkipIntegrations = false
					}

					filteredRecipes = append(filteredRecipes, r)
				}
			}
		}

		log.Debug("skipping recipes that were not selected")
		for _, r := range recipes {
			if !i.recipeInRecipes(r, filteredRecipes) {
				i.status.ReportRecipeSkipped(execution.RecipeStatusEvent{
					Recipe: r,
				})
			}
		}

		recipes = filteredRecipes
	}

	return recipes, nil
}
