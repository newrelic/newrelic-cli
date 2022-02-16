package install

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpproxy"

	"github.com/newrelic/newrelic-cli/internal/cli"
	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

type RecipeInstaller struct {
	types.InstallerContext
	discoverer                  Discoverer
	manifestValidator           *discovery.ManifestValidator
	recipeFetcher               recipes.RecipeFetcher
	recipeExecutor              execution.RecipeExecutor
	recipeValidator             RecipeValidator
	recipeFileFetcher           RecipeFileFetcher
	status                      *execution.InstallStatus
	prompter                    Prompter
	executionProgressIndicator  ux.ProgressIndicator
	validationProgressIndicator ux.ProgressIndicator
	licenseKeyFetcher           LicenseKeyFetcher
	configValidator             ConfigValidator
	recipeVarPreparer           RecipeVarPreparer
	recipeFilterer              RecipeFilterRunner
	agentValidator              *validation.AgentValidator
}

type RecipeInstallFunc func(ctx context.Context, i *RecipeInstaller, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) error

const (
	validationTimeout       = 5 * time.Minute
	validationInProgressMsg = "Checking for data in New Relic (this may take a few minutes)..."
)

func NewRecipeInstaller(ic types.InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstaller {
	checkNetwork(nrClient)

	var recipeFetcher recipes.RecipeFetcher

	if ic.LocalRecipes != "" {
		recipeFetcher = &recipes.LocalRecipeFetcher{
			Path: ic.LocalRecipes,
		}
	} else {
		recipeFetcher = recipes.NewEmbeddedRecipeFetcher()
	}

	mv := discovery.NewManifestValidator()
	ff := recipes.NewRecipeFileFetcher()
	ers := []execution.StatusSubscriber{
		execution.NewNerdStorageStatusReporter(&nrClient.NerdStorage),
		execution.NewTerminalStatusReporter(),
		execution.NewInstallEventsReporter(&nrClient.InstallEvents),
	}
	lkf := NewServiceLicenseKeyFetcher(&nrClient.NerdGraph)
	slg := execution.NewPlatformLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)

	d := discovery.NewPSUtilDiscoverer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	cv := diagnose.NewConfigValidator(nrClient)
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()
	sp := ux.NewSpinner()
	rvp := execution.NewRecipeVarProvider()
	rf := recipes.NewRecipeFilterRunner(ic, statusRollup)
	av := validation.NewAgentValidator()

	i := RecipeInstaller{
		discoverer:                  d,
		manifestValidator:           mv,
		recipeFetcher:               recipeFetcher,
		recipeExecutor:              re,
		recipeValidator:             v,
		recipeFileFetcher:           ff,
		status:                      statusRollup,
		prompter:                    p,
		executionProgressIndicator:  pi,
		validationProgressIndicator: sp,
		licenseKeyFetcher:           lkf,
		configValidator:             cv,
		recipeVarPreparer:           rvp,
		recipeFilterer:              rf,
		agentValidator:              av,
	}

	i.InstallerContext = ic

	return &i
}

func (i *RecipeInstaller) promptIfNotLatestCLIVersion(ctx context.Context) error {
	latestReleaseVersion, err := cli.GetLatestReleaseVersion(ctx)
	if err != nil {
		log.Debug(err)
		return nil
	}

	isLatestCLIVersion, err := cli.IsLatestVersion(ctx, latestReleaseVersion)
	if err != nil {
		log.Debug(err)
		return nil
	}

	if !isLatestCLIVersion {
		i.status.UpdateRequired = true

		cli.PrintUpdateCLIMessage(latestReleaseVersion)

		err := &types.UpdateRequiredError{
			Err:     fmt.Errorf(`%s`, cli.FormatUpdateVersionMessage(latestReleaseVersion)),
			Details: "UpdateRequiredError",
		}
		return err
	}

	return nil
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
	fmt.Println()

	log.Tracef("InstallerContext: %+v", i.InstallerContext)
	log.WithFields(log.Fields{
		"RecipesProvided":     i.RecipesProvided(),
		"RecipePathsProvided": i.RecipePathsProvided(),
		"RecipeNamesProvided": i.RecipeNamesProvided(),
	}).Debug("context summary")

	ctx, cancel := context.WithCancel(utils.SignalCtx)
	defer cancel()

	errChan := make(chan error)
	var err error

	// Test split service
	// treatment := split.Service.Get(split.VirtuosoCLITest)
	// log.Printf("Got treatment: %s for %s", treatment, split.VirtuosoCLITest)

	log.Printf("Validating connectivity to the New Relic platform...")
	if err = i.configValidator.Validate(utils.SignalCtx); err != nil {
		return err
	}

	if i.RecipesProvided() {
		i.status.SetTargetedInstall()
	}

	i.status.InstallStarted()

	// If not in a dev environemt, check to see if
	// the installed CLI is up to date.
	if !cli.IsDevEnvironment() {
		if err = i.promptIfNotLatestCLIVersion(ctx); err != nil {
			i.status.InstallComplete(err)
			return err
		}
	}

	go func(ctx context.Context) {
		errChan <- i.install(ctx)
	}(ctx)

	select {
	case <-ctx.Done():
		i.status.InstallCanceled()
		return nil
	case err = <-errChan:
		if errors.Is(err, types.ErrInterrupt) {
			i.status.InstallCanceled()
			return err
		}

		i.status.InstallComplete(err)

		return err
	}
}

func (i *RecipeInstaller) install(ctx context.Context) error {
	installLibraryVersion := i.recipeFetcher.FetchLibraryVersion(ctx)
	log.Debugf("Using open-install-library version %s", installLibraryVersion)
	i.status.SetVersions(installLibraryVersion)

	// Execute the discovery process, exiting on failure.
	m, err := i.discover(ctx)
	if err != nil {
		return err
	}

	err = i.assertDiscoveryValid(ctx, m)
	i.status.DiscoveryComplete(*m)

	if err != nil {
		return err
	}

	repo := recipes.NewRecipeRepository(func() ([]types.OpenInstallationRecipe, error) {
		recipes, err2 := i.recipeFetcher.FetchRecipes(ctx)
		return recipes, err2
	}, m)

	//FIXME: need to fix

	bundler := recipes.NewBundler(ctx, repo)
	coreBundle := bundler.CreateCoreBundle()
	bundlerInstaller := NewBundleInstaller(ctx, m, i)
	err = bundlerInstaller.InstallStopOnError(coreBundle, true)
	if err != nil {
		log.Debugf("error installing core bundle: %s", err)
		return err
	}

	return nil

	// err = i.intallBundle(ctx, m, coreBundle)
	// if err != nil {
	// 	log.Debugf("Unable to load install core recipes, detail: %s", err)
	// 	return err
	// }

	// additionalBundle := bundler.createAdditionalBundle(coreBundle)

	// err = i.intallBundle(ctx, m, additionalBundle)
	// if err != nil {
	// 	log.Debugf("Unable to load install core recipes, detail: %s", err)
	// 	return err
	// }

	// return nil

	// var recipesForInstall []types.OpenInstallationRecipe
	// if i.RecipesProvided() {
	// 	recipesForInstall, err = i.fetchProvidedRecipe(m, recipesForPlatform)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Debugf("recipes provided:\n")
	// 	logRecipes(recipesForInstall)

	// 	if err = i.recipeFilterer.EnsureDoesNotFilter(ctx, recipesForInstall, m); err != nil {
	// 		return err
	// 	}

	// } else {
	// 	var selected, unselected []types.OpenInstallationRecipe

	// 	recipesForInstall = i.recipeFilterer.RunFilterAll(ctx, recipesForPlatform, m)
	// 	log.Debugf("recipes after filtering:")
	// 	logRecipes(recipesForInstall)

	// 	selected, unselected, err = i.promptUserSelect(recipesForInstall)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Tracef("recipes selected by user: %v\n", selected)

	// 	for _, r := range unselected {
	// 		i.status.RecipeSkipped(execution.RecipeStatusEvent{Recipe: r})
	// 	}

	// 	recipesForInstall = selected
	// }

	// i.status.RecipesSelected(recipesForInstall)

	// dependencies := resolveDependencies(recipesForInstall, recipesForPlatform)
	// recipesForInstall = addIfMissing(recipesForInstall, dependencies)

	// if err = i.installRecipes(ctx, m, recipesForInstall); err != nil {
	// 	return err
	// }

	// log.Debugf("Done installing.")
	// return nil
}

func (i *RecipeInstaller) assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error {
	err := i.manifestValidator.Validate(m)
	if err != nil {
		return err
	}
	log.Debugf("Done asserting valid operating system for OS:%s and PlatformVersion:%s", m.OS, m.PlatformVersion)
	return nil
}

// func (i *RecipeInstaller) installRecipes(ctx context.Context, m *types.DiscoveryManifest, recipes []types.OpenInstallationRecipe) error {
// 	log.WithFields(log.Fields{
// 		"recipe_count": len(recipes),
// 	}).Debug("installing recipes")
// 	var lastError error

// 	for _, r := range recipes {
// 		var err error

// 		log.WithFields(log.Fields{
// 			"name": r.Name,
// 		}).Debug("installing recipe")

// 		_, err = i.executeAndValidateWithProgress(ctx, m, &r)

// 		if err != nil {
// 			if err == types.ErrInterrupt {
// 				return err
// 			}

// 			if r.Name == types.InfraAgentRecipeName || r.Name == types.LoggingRecipeName || i.RecipesProvided() {
// 				return err
// 			}

// 			lastError = err

// 			log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", r.Name, err)
// 			log.Debug(err)
// 		}
// 		log.Debugf("Done executing and validating with progress for recipe name %s.", r.Name)
// 	}

// 	if lastError != nil {
// 		// Return last recipe error that was caught if any
// 		return lastError
// 	}

// 	if !i.status.WasSuccessful() {
// 		return &types.UncaughtError{
// 			Err: fmt.Errorf("no recipes were installed"),
// 		}
// 	}

// 	return nil
// }

func (i *RecipeInstaller) discover(ctx context.Context) (*types.DiscoveryManifest, error) {
	log.Debug("discovering system information")

	m, err := i.discoverer.Discover(ctx)
	if err != nil {
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

	return m, nil
}

func (i *RecipeInstaller) executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars) (string, error) {
	i.status.RecipeInstalling(execution.RecipeStatusEvent{Recipe: *r})

	// Execute the recipe steps.
	if err := i.recipeExecutor.Execute(ctx, *r, vars); err != nil {
		if err == types.ErrInterrupt {
			return "", err
		}

		if e, ok := err.(*types.UnsupportedOperatingSytemError); ok {
			i.status.RecipeUnsupported(execution.RecipeStatusEvent{
				Recipe: *r,
				Msg:    e.Error(),
			})

			return "", err
		}

		msg := fmt.Sprintf("execution failed for %s: %s", r.Name, err)
		se := execution.RecipeStatusEvent{
			Recipe: *r,
			Msg:    msg,
		}

		if e, ok := err.(types.GoTaskError); ok {
			e.SetError(msg)
			se.TaskPath = e.TaskPath()
		} else {
			err = errors.New(msg)
		}

		i.status.RecipeFailed(se)
		return "", err
	}

	validationStart := time.Now()
	entityGUID, err := i.validateRecipeViaAllMethods(ctx, r, m, vars)
	validationDurationMs := time.Since(validationStart).Milliseconds()
	if err != nil {
		validationErr := fmt.Errorf("encountered an error while validating receipt of data for %s: %w", r.Name, err)
		i.status.RecipeFailed(execution.RecipeStatusEvent{
			Recipe:               *r,
			Msg:                  validationErr.Error(),
			ValidationDurationMs: validationDurationMs,
		})

		return "", validationErr
	}

	i.status.RecipeInstalled(execution.RecipeStatusEvent{
		Recipe:               *r,
		EntityGUID:           entityGUID,
		ValidationDurationMs: validationDurationMs,
	})

	return entityGUID, nil
}

type validationFunc func() (string, error)

func (i *RecipeInstaller) validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars) (string, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, validationTimeout)
	defer cancel()

	entityGUIDChan := make(chan string)
	validationErrorChan := make(chan error)
	validationErrors := []error{}

	validationFuncs := []validationFunc{}

	// Add agent validation if configured
	hasValidationURL := r.ValidationURL != ""
	isAbsoluteURL := utils.IsAbsoluteURL(r.ValidationURL)
	if hasValidationURL && isAbsoluteURL {
		validationFuncs = append(validationFuncs, func() (string, error) {
			return i.agentValidator.Validate(timeoutCtx, r.ValidationURL)
		})
	} else {
		log.Debugf("skipping agent validation due to lack of validationUrl")
	}

	// Add NRQL validation if configured
	if r.ValidationNRQL != "" {
		validationFuncs = append(validationFuncs, func() (string, error) {
			return i.recipeValidator.ValidateRecipe(timeoutCtx, *m, *r, vars)
		})
	} else {
		log.Debugf("skipping NRQL validation due to lack of validationNRQL")
	}

	if len(validationFuncs) == 0 {
		log.Debugf("skipping recipe validation since no validation targets were configured")
		return "", nil
	}

	i.validationProgressIndicator.Start(validationInProgressMsg)
	defer i.validationProgressIndicator.Stop()

	for _, f := range validationFuncs {
		go func(fn validationFunc) {
			entityGUID, err := fn()
			if err != nil {
				validationErrorChan <- err
				return
			}

			entityGUIDChan <- entityGUID
		}(f)
	}

	for {
		select {
		case entityGUID := <-entityGUIDChan:
			i.validationProgressIndicator.Success("")
			return entityGUID, nil
		case err := <-validationErrorChan:
			validationErrors = append(validationErrors, err)
			log.Debugf("validation error encountered: %s", err)

			if len(validationErrors) == len(validationFuncs) {
				i.validationProgressIndicator.Fail("")
				return "", fmt.Errorf("no validation was successful.  most recent validation error: %w", err)
			}
		case <-timeoutCtx.Done():
			return "", fmt.Errorf("timed out waiting for validation to succeed")
		}
	}
}

func (i *RecipeInstaller) executeAndValidateWithProgress(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, assumeYes bool) (string, error) {
	msg := fmt.Sprintf("Installing %s", r.DisplayName)
	i.executionProgressIndicator.Start(msg)
	defer func() { i.executionProgressIndicator.Stop() }()

	if r.PreInstallMessage() != "" {
		fmt.Println(r.PreInstallMessage())
	}

	licenseKey, err := i.licenseKeyFetcher.FetchLicenseKey(ctx)
	if err != nil {
		return "", err
	}

	vars, err := i.recipeVarPreparer.Prepare(*m, *r, assumeYes, licenseKey)
	if err != nil {
		return "", err
	}

	entityGUID, err := i.executeAndValidate(ctx, m, r, vars)
	if err != nil {
		if errors.Is(err, types.ErrInterrupt) {
			i.executionProgressIndicator.Canceled(msg)
		} else {
			i.executionProgressIndicator.Fail(msg)
		}
		return "", err
	}

	if r.PostInstallMessage() != "" {
		fmt.Println(r.PostInstallMessage())
	}

	i.executionProgressIndicator.Success(msg)
	return entityGUID, nil
}

// func (i *RecipeInstaller) fetchProvidedRecipe(m *types.DiscoveryManifest, recipesForPlatform []types.OpenInstallationRecipe) ([]types.OpenInstallationRecipe, error) {
// 	var recipes []types.OpenInstallationRecipe

// 	// Load the recipes from the provided file names.
// 	for _, n := range i.RecipePaths {
// 		log.Debugln(fmt.Sprintf("Attempting to match recipePath %s.", n))
// 		recipe, err := i.recipeFromPath(n)
// 		if err != nil {
// 			log.Debugln(fmt.Sprintf("Error while building recipe from path, detail:%s.", err))
// 			return nil, err
// 		}

// 		log.WithFields(log.Fields{
// 			"name":         recipe.Name,
// 			"display_name": recipe.DisplayName,
// 			"path":         n,
// 		}).Debug("found recipe at path")

// 		recipes = append(recipes, *recipe)
// 	}

// 	// Load the recipes from the provided file names.
// 	for _, n := range i.RecipeNames {
// 		found := false
// 		log.Debugln(fmt.Sprintf("Attempting to match recipe name %s.", n))
// 		for _, r := range recipesForPlatform {
// 			if strings.EqualFold(r.Name, n) {
// 				log.WithFields(log.Fields{
// 					"name":         r.Name,
// 					"display_name": r.DisplayName,
// 				}).Debug("found recipe with name")
// 				recipes = append(recipes, r)
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			log.Errorf("Could not find recipe with name %s.", n)
// 		}
// 	}

// 	return recipes, nil
// }

// func (i *RecipeInstaller) recipeFromPath(recipePath string) (*types.OpenInstallationRecipe, error) {
// 	recipeURL, parseErr := url.Parse(recipePath)
// 	if parseErr == nil && recipeURL.Scheme != "" && strings.HasPrefix(strings.ToLower(recipeURL.Scheme), "http") {
// 		f, err := i.recipeFileFetcher.FetchRecipeFile(recipeURL)
// 		if err != nil {
// 			return nil, fmt.Errorf("could not fetch file %s: %s", recipePath, err)
// 		}
// 		return f, nil
// 	}

// 	f, err := i.recipeFileFetcher.LoadRecipeFile(recipePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not load file %s: %s", recipePath, err)
// 	}

// 	return f, nil
// }

// func (i *RecipeInstaller) promptUserSelect(recipes []types.OpenInstallationRecipe) ([]types.OpenInstallationRecipe, []types.OpenInstallationRecipe, error) {
// 	if len(recipes) == 0 {
// 		return []types.OpenInstallationRecipe{}, []types.OpenInstallationRecipe{}, nil
// 	}

// 	if i.AssumeYes {
// 		return recipes, []types.OpenInstallationRecipe{}, nil
// 	}

// 	var selectedRecipes, unselectedRecipes []types.OpenInstallationRecipe

// 	names := []string{}
// 	selected := []string{}
// 	for _, r := range recipes {
// 		if r.Name != types.InfraAgentRecipeName {
// 			names = append(names, r.DisplayName)
// 		} else {
// 			fmt.Printf("The guided installation will begin by installing the latest version of the New Relic Infrastructure agent, which is required for additional instrumentation.\n\n")
// 		}
// 	}

// 	if len(names) > 0 {
// 		var promptErr error
// 		selected, promptErr = i.prompter.MultiSelect("Please choose from the following instrumentation to be installed:", names)
// 		if promptErr != nil {
// 			return nil, nil, promptErr
// 		}
// 		fmt.Println()
// 	}

// 	for _, r := range recipes {
// 		if utils.StringInSlice(r.DisplayName, selected) || r.Name == types.InfraAgentRecipeName {
// 			selectedRecipes = append(selectedRecipes, r)
// 		} else {
// 			unselectedRecipes = append(unselectedRecipes, r)
// 		}
// 	}

// 	return selectedRecipes, unselectedRecipes, nil
// }

// func logRecipes(recipes []types.OpenInstallationRecipe) {
// 	for _, r := range recipes {
// 		log.Debugf("%s", r.ToShortDisplayString())
// 	}
// }

// func findRecipeInRecipes(name string, recipes []types.OpenInstallationRecipe) *types.OpenInstallationRecipe {
// 	for _, r := range recipes {
// 		if r.Name == name {
// 			return &r
// 		}
// 	}

// 	return nil
// }

// This is a naive implementation that only resolves dependencies one level deep.
// func resolveDependencies(recipes []types.OpenInstallationRecipe, recipesForPlatform []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
// 	var results []types.OpenInstallationRecipe

// 	for _, r := range recipes {
// 		if len(r.Dependencies) > 0 {
// 			for _, n := range r.Dependencies {
// 				d := findRecipeInRecipes(n, recipesForPlatform)
// 				if d != nil {
// 					results = append(results, *d)
// 				}
// 			}
// 		}
// 	}

// 	return results
// }

// This is a naive implementation that only resolves dependencies one level deep.
// func addIfMissing(recipes []types.OpenInstallationRecipe, dependencies []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
// 	var results []types.OpenInstallationRecipe

// 	for _, d := range dependencies {
// 		if found := findRecipeInRecipes(d.Name, results); found == nil {
// 			results = append(results, d)
// 		}
// 	}

// 	for _, r := range recipes {
// 		if found := findRecipeInRecipes(r.Name, results); found == nil {
// 			results = append(results, r)
// 		}
// 	}

// 	return results
// }

func checkNetwork(nrClient *newrelic.NewRelic) {
	err := nrClient.TestEndpoints()
	if err != nil {

		proxyConfig := httpproxy.FromEnvironment()

		log.Debugf("proxyConfig: %+v", proxyConfig)
		if proxyConfig.HTTPProxy != "" || proxyConfig.HTTPSProxy != "" || proxyConfig.NoProxy != "" {
			log.Warn("Proxy settings have been configured but we are still unable to connect to the New Relic platform.  You may need to adjust your proxy environment variables.  https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md#using-an-http-proxy")
		}

		log.Error(err)
	}
}
