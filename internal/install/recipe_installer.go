package install

import (
	"context"
	"errors"
	"fmt"
	"os"
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

type RecipeInstall struct {
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
	agentValidator              *validation.AgentValidator
}

type RecipeInstallFunc func(ctx context.Context, i *RecipeInstall, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) error

const (
	validationTimeout       = 5 * time.Minute
	validationInProgressMsg = "Checking for data in New Relic (this may take a few minutes)..."
)

var statusRollup *execution.InstallStatus

func NewRecipeInstaller(ic types.InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstall {
	checkNetwork(nrClient)

	var recipeFetcher recipes.RecipeFetcher

	if ic.LocalRecipes != "" {
		recipeFetcher = &recipes.LocalRecipeFetcher{
			Path: ic.LocalRecipes,
		}
	} else if len(ic.RecipePaths) > 0 {
		recipeFetcher = &recipes.RecipeFileFetcher{
			Paths: ic.RecipePaths,
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
	statusRollup = execution.NewInstallStatus(ers, slg)

	d := discovery.NewPSUtilDiscoverer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	cv := diagnose.NewConfigValidator(nrClient)
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()
	sp := ux.NewSpinner()
	rvp := execution.NewRecipeVarProvider()
	av := validation.NewAgentValidator()

	i := RecipeInstall{
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
		agentValidator:              av,
	}

	i.InstallerContext = ic

	return &i
}

func (i *RecipeInstall) promptIfNotLatestCLIVersion(ctx context.Context) error {
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

func (i *RecipeInstall) Install() error {
	fmt.Printf(`
_   _                 ____      _ _
| \ | | _____      __ |  _ \ ___| (_) ___
|  \| |/ _ \ \ /\ / / | |_) / _ | | |/ __|
| |\  |  __/\ V  V /  |  _ |  __| | | (__
|_| \_|\___| \_/\_/   |_| \_\___|_|_|\___|

Welcome to New Relic. Let's set up full stack observability for your environment. 
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

	err = i.connectToPlatform()

	if err != nil {
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

func (i *RecipeInstall) connectToPlatform() error {
	loaderChan := make(chan error)

	go func() {
		err := i.configValidator.Validate(utils.SignalCtx)
		if err != nil {
			loaderChan <- err
		}
		loaderChan <- nil
	}()

	welcomeScreenProgressBar := ux.NewSpinnerProgressIndicator()
	welcomeScreenProgressBar.Start("Connecting to New Relic Platform")

	loaded := <-loaderChan

	if loaded == nil {
		welcomeScreenProgressBar.Success("Connecting to New Relic Platform")
	} else {
		welcomeScreenProgressBar.Fail("Connecting to New Relic Platform")
	}
	return loaded
}

func OSEnvVariableGetter(name string) string {
	return os.Getenv(name)
}

//TODO: needs to skipcore, skipcore with assume yes, not skipping core
var EnvVariableGetter = OSEnvVariableGetter

func (i *RecipeInstall) install(ctx context.Context) error {
	installLibraryVersion := i.recipeFetcher.FetchLibraryVersion(ctx)
	log.Debugf("Using open-install-library version %s", installLibraryVersion)
	i.status.SetVersions(installLibraryVersion)

	fmt.Println("\n\nInstalling New Relic")
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

	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		recipes, err2 := i.recipeFetcher.FetchRecipes(ctx)
		return recipes, err2
	}, m)

	//FIXME: need to fix
	bundler := recipes.NewBundler(ctx, repo)
	bundleInstaller := NewBundleInstaller(ctx, m, i, statusRollup)

	installCoreBundle := EnvVariableGetter("NEW_RELIC_CLI_SKIP_CORE") != "1"

	if installCoreBundle {
		coreBundle := bundler.CreateCoreBundle()
		err = bundleInstaller.InstallStopOnError(coreBundle, true)
		if err != nil {
			log.Debugf("error installing core bundle: %s", err)
			return err
		}
	}

	//FIXME: additional install mock, just hack together code for install to check flow, needs to be refactor
	var additionalBundle *recipes.Bundle
	if i.RecipeNamesProvided() {
		additionalBundle = bundler.CreateAdditionalTargetedBundle(i.RecipeNames)
	} else {
		additionalBundle = bundler.CreateAdditionalGuidedBundle()
	}
	bundleInstaller.InstallContinueOnError(additionalBundle, i.AssumeYes)

	if bundleInstaller.InstalledRecipesCount() == 0 {
		return &types.UncaughtError{
			Err: fmt.Errorf("no recipes were installed"),
		}
	}

	log.Debugf("Done installing.")

	return nil

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

func (i *RecipeInstall) assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error {
	err := i.manifestValidator.Validate(m)
	if err != nil {
		return err
	}
	log.Debugf("Done asserting valid operating system for OS:%s and PlatformVersion:%s", m.OS, m.PlatformVersion)
	return nil
}

func (i *RecipeInstall) discover(ctx context.Context) (*types.DiscoveryManifest, error) {
	log.Debug("discovering system information")

	m, err := i.discoverer.Discover(ctx)
	if err != nil {
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

	return m, nil
}

// intalling recipe
func (i *RecipeInstall) executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars) (string, error) {
	i.status.RecipeInstalling(execution.RecipeStatusEvent{Recipe: *r})

	// Execute the recipe steps.
	if err := i.recipeExecutor.Execute(ctx, *r, vars); err != nil {
		if err == types.ErrInterrupt {
			return "", err
		}

		if e, ok := err.(*types.UnsupportedOperatingSystemError); ok {
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

// Post install validation
func (i *RecipeInstall) validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars) (string, error) {
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

	log.Debug(validationInProgressMsg)
	//i.validationProgressIndicator.Start(validationInProgressMsg)
	//defer i.validationProgressIndicator.Stop()

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

// Installing recipe
func (i *RecipeInstall) executeAndValidateWithProgress(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, assumeYes bool) (string, error) {

	fmt.Println()
	msg := fmt.Sprintf("Installing %s", r.DisplayName)

	errorChan := make(chan error)
	successChan := make(chan string)

	go func() {
		licenseKey, err := i.licenseKeyFetcher.FetchLicenseKey(ctx)
		if err != nil {
			errorChan <- err
		}

		vars, err := i.recipeVarPreparer.Prepare(*m, *r, assumeYes, licenseKey)
		if err != nil {
			errorChan <- err
		}

		vars["assumeYes"] = fmt.Sprintf("%v", assumeYes)

		entityGUID, err := i.executeAndValidate(ctx, m, r, vars)

		if err != nil {
			errorChan <- err
		}
		successChan <- entityGUID
	}()

	installProgressBar := ux.NewSpinnerProgressIndicator()
	installProgressBar.AssumeYes = assumeYes
	installProgressBar.Start(msg)

	for {
		select {
		case entityGUID := <-successChan:
			installProgressBar.Success("Installing " + r.DisplayName)

			return entityGUID, nil
		case err := <-errorChan:
			if errors.Is(err, types.ErrInterrupt) {
				installProgressBar.Canceled("Installing " + r.DisplayName)
			} else {
				installProgressBar.Fail("Installing " + r.DisplayName)
			}
			log.Debugf("install error encountered: %s", err)
			return "", err
		}
	}
}

// func (i *RecipeInstall) fetchProvidedRecipe(m *types.DiscoveryManifest, recipesForPlatform []types.OpenInstallationRecipe) ([]types.OpenInstallationRecipe, error) {
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
