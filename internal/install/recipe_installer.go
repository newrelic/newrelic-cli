package install

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
)

const (
	validationTimeout = 5 * time.Minute
)

type RecipeInstall struct {
	types.InstallerContext
	recipeFetcher          recipes.RecipeFetcher
	recipeExecutor         execution.RecipeExecutor
	recipeValidator        RecipeValidator
	recipeFileFetcher      RecipeFileFetcher
	recipeLogForwarder     execution.LogForwarder
	status                 *execution.InstallStatus
	prompter               Prompter
	licenseKeyFetcher      LicenseKeyFetcher
	configValidator        ConfigValidator
	recipeVarPreparer      RecipeVarPreparer
	agentValidator         AgentValidator
	shouldInstallCore      func() bool
	bundlerFactory         func(ctx context.Context, availableRecipes recipes.RecipeDetectionResults) RecipeBundler
	bundleInstallerFactory func(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) RecipeBundleInstaller
	progressIndicator      ux.ProgressIndicator
	recipeDetectorFactory  func(ctx context.Context, repo *recipes.RecipeRepository) RecipeStatusDetector
	processEvaluator       recipes.ProcessEvaluatorInterface
}

type RecipeInstallFunc func(ctx context.Context, i *RecipeInstall, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) error

func NewRecipeInstaller(ic types.InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstall {
	checkNetwork(nrClient)

	var recipeFetcher recipes.RecipeFetcher

	if ic.LocalRecipes != "" {
		recipeFetcher = &recipes.LocalRecipeFetcher{
			Path: ic.LocalRecipes,
		}
	} else if len(ic.RecipePaths) > 0 {
		recipeFetcher = recipes.NewRecipeFileFetcher(ic.RecipePaths)
	} else {
		recipeFetcher = recipes.NewEmbeddedRecipeFetcher()
	}

	//mv := discovery.NewManifestValidator()
	ff := recipes.NewRecipeFileFetcher([]string{})
	lf := execution.NewRecipeLogForwarder()
	ers := []execution.StatusSubscriber{
		execution.NewNerdStorageStatusReporter(&nrClient.NerdStorage),
		execution.NewTerminalStatusReporter(),
		execution.NewInstallEventsReporter(&nrClient.InstallEvents),
	}
	lkf := NewServiceLicenseKeyFetcher(&nrClient.NerdGraph)
	slg := execution.NewPlatformLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)

	//d := discovery.NewPSUtilDiscoverer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	cv := diagnose.NewConfigValidator(nrClient)
	p := ux.NewPromptUIPrompter()
	rvp := execution.NewRecipeVarProvider()
	av := validation.NewAgentValidator()

	i := RecipeInstall{
		//discoverer:         d,
		//manifestValidator:  mv,
		recipeFetcher:      recipeFetcher,
		recipeExecutor:     re,
		recipeValidator:    v,
		recipeFileFetcher:  ff,
		recipeLogForwarder: lf,
		status:             statusRollup,
		prompter:           p,
		licenseKeyFetcher:  lkf,
		configValidator:    cv,
		recipeVarPreparer:  rvp,
		agentValidator:     av,
		progressIndicator:  ux.NewSpinnerProgressIndicator(),
		processEvaluator:   recipes.NewProcessEvaluator(),
	}

	i.InstallerContext = ic

	i.shouldInstallCore = func() bool {
		return os.Getenv("NEW_RELIC_CLI_SKIP_CORE") != "1"
	}

	i.bundlerFactory = func(ctx context.Context, availableRecipes recipes.RecipeDetectionResults) RecipeBundler {
		return recipes.NewBundler(ctx, availableRecipes)
	}

	i.bundleInstallerFactory = func(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) RecipeBundleInstaller {
		return NewBundleInstaller(ctx, manifest, recipeInstallerInterface, statusReporter)
	}
	i.recipeDetectorFactory = func(ctx context.Context, repo *recipes.RecipeRepository) RecipeStatusDetector {
		return recipes.NewRecipeDetector(ctx, repo, i.processEvaluator)
	}
	return &i
}

var discoverManifest = func(ctx context.Context, discoverer Discoverer) (*types.DiscoveryManifest, error) {
	log.Debug("discovering system information")
	//discoverer := discovery.NewPSUtilDiscoverer()
	m, err := discoverer.Discover(ctx)
	if err != nil {
		return nil, fmt.Errorf("there was an error discovering system info: %s", err)
	}

	return m, nil
}

var validateManifest = func(m *types.DiscoveryManifest) error {
	manifestValidator := discovery.NewManifestValidator()
	err := manifestValidator.Validate(m)
	if err != nil {
		return err
	}
	log.Debugf("Done asserting valid operating system for OS:%s and PlatformVersion:%s", m.OS, m.PlatformVersion)
	return nil
}

var getLatestCliVersionReleased = func(ctx context.Context) (string, error) {
	return cli.GetLatestReleaseVersion(ctx)
}

var isLatestCliVersionInstalled = func(ctx context.Context, version string) (bool, error) {
	return cli.IsLatestVersion(ctx, version)
}

func (i *RecipeInstall) promptIfNotLatestCLIVersion(ctx context.Context) error {
	latestCliVersion, err := getLatestCliVersionReleased(ctx)
	if err != nil {
		log.Debug(err)
		return nil
	}

	isMostRecentCliVersion, err := isLatestCliVersionInstalled(ctx, latestCliVersion)
	if err != nil {
		log.Debug(err)
		return nil
	}

	if !isMostRecentCliVersion {
		i.status.UpdateRequired = true

		cli.PrintUpdateCLIMessage(latestCliVersion)

		err := &types.UpdateRequiredError{
			Err:     fmt.Errorf(`%s`, cli.FormatUpdateVersionMessage(latestCliVersion)),
			Details: "UpdateRequiredError",
		}
		return err
	}

	return nil
}

func (i *RecipeInstall) ensureSingleConcurrentInstall(ctx context.Context) error {
	processes := i.processEvaluator.GetOrLoadProcesses(ctx)
	count := 0
	nameRegex := regexp.MustCompile(`(?i)newrelic(\.exe)?`)
	for _, p := range processes {
		name, err := p.Name()
		if err == nil && nameRegex.MatchString(name) {
			cmd, err := p.Cmd()
			cmdRegex := regexp.MustCompile(`(?i)newrelic(\.exe)? install`)
			if err == nil && cmdRegex.MatchString(cmd) {
				log.Debugf(fmt.Sprintf("EnsureSingleConcurrentInstall Matched:%s pid:%d", name, p.PID()))
				count++
			}
		}
	}
	if count > 1 {
		message := fmt.Sprintf("only 1 newrelic install command can run at one time, found %d currently executing. Please retry later, or terminate the other newrelic installations", count)
		return errors.New(message)
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
		"RecipePathsProvided": i.RecipePathsProvided(),
		"RecipeNamesProvided": i.RecipeNamesProvided(),
	}).Debug("context summary")

	if i.RecipeNamesProvided() {
		i.status.SetTargetedInstall(i.RecipeNames)
	}

	if err := i.ensureSingleConcurrentInstall(utils.SignalCtx); err != nil {
		log.Debug(err)
		return err
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		message := "This system is not supported for automatic installation, no host info. Please see our documentation for requirements."
		return errors.New(message)
	}

	i.status.InstallStarted()

	ctx, cancel := context.WithCancel(utils.SignalCtx)
	defer cancel()

	errChan := make(chan error)
	var err error

	err = i.connectToPlatform()
	if err != nil {
		i.status.InstallComplete(diagnose.ConnectionError{Err: err})
		return err
	}

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
			return
		}
		loaderChan <- nil
	}()

	i.progressIndicator.Start("Connecting to New Relic Platform")

	loaded := <-loaderChan

	if loaded == nil {
		i.progressIndicator.Success("Connecting to New Relic Platform")
	} else {
		i.progressIndicator.Fail("Connecting to New Relic Platform")
	}
	return loaded
}

func (i *RecipeInstall) install(ctx context.Context) error {

	installLibraryVersion := i.recipeFetcher.FetchLibraryVersion(ctx)
	log.Debugf("Using open-install-library version %s", installLibraryVersion)
	i.status.SetVersions(installLibraryVersion)

	// Execute the discovery process and validate, exiting on failure.
	discoveryManifest, discoveryError := discoverManifest(ctx, discovery.NewPSUtilDiscoverer())
	if discoveryError != nil {
		return discoveryError
	}
	i.status.DiscoveryComplete(*discoveryManifest)
	invalidManifestError := validateManifest(discoveryManifest)
	if invalidManifestError != nil {
		return invalidManifestError
	}

	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		recipes, err2 := i.recipeFetcher.FetchRecipes(ctx)
		return recipes, err2
	}, discoveryManifest)

	i.printStartInstallingMessage(repo)

	recipeDetector := i.recipeDetectorFactory(ctx, repo)
	availableRecipes, unavailableRecipes, err := recipeDetector.GetDetectedRecipes()
	if err != nil {
		return err
	}

	i.reportRecipeStatuses(availableRecipes, unavailableRecipes)

	if len(availableRecipes) == 0 && !i.RecipeNamesProvided() {
		fmt.Println("This system is not supported by any available recipes for automatic installation. Please see our documentation for requirements.")
		return &types.UncaughtError{
			Err: fmt.Errorf("no recipes found supporting this system"),
		}
	}

	bundler := i.bundlerFactory(ctx, availableRecipes)
	bundleInstaller := i.bundleInstallerFactory(ctx, discoveryManifest, i, i.status)

	cbErr := i.installCoreBundle(bundler, bundleInstaller)
	if cbErr != nil {
		return cbErr
	}

	abErr := i.installAdditionalBundle(bundler, bundleInstaller, repo)
	if abErr != nil {
		return abErr
	}

	i.reportRecipeRecommendations(availableRecipes)

	log.Debugf("Done installing.")
	return nil
}

func (i *RecipeInstall) printStartInstallingMessage(repo *recipes.RecipeRepository) {
	message := "\n\nInstalling New Relic"
	if i.RecipeNamesProvided() && len(i.RecipeNames) > 0 {
		r := repo.FindRecipeByName(i.RecipeNames[0])
		if r != nil {
			message = fmt.Sprintf("%s %s", message, r.DisplayName)
		}
	}
	fmt.Println(message)
}

func (i *RecipeInstall) reportRecipeStatuses(availableRecipes recipes.RecipeDetectionResults,
	unavailableRecipes recipes.RecipeDetectionResults) {

	for _, d := range unavailableRecipes {
		e := execution.RecipeStatusEvent{Recipe: *d.Recipe, ValidationDurationMs: d.DurationMs}
		i.status.ReportStatus(d.Status, e)
	}
	for _, d := range availableRecipes {
		e := execution.RecipeStatusEvent{Recipe: *d.Recipe, ValidationDurationMs: d.DurationMs}
		i.status.ReportStatus(execution.RecipeStatusTypes.DETECTED, e)
	}
}

func (i *RecipeInstall) reportRecipeRecommendations(availableRecipes recipes.RecipeDetectionResults) {
	recipeRecommendations := i.getRecipeRecommendations(availableRecipes)
	for _, e := range recipeRecommendations {
		i.status.ReportStatus(execution.RecipeStatusTypes.RECOMMENDED, e)
	}
}

// Report recommendations for recipes that are available but have not been attempted to be installed
func (i *RecipeInstall) getRecipeRecommendations(availableRecipes recipes.RecipeDetectionResults) []execution.RecipeStatusEvent {
	isTargetedInstall := i.RecipeNamesProvided() && len(i.RecipeNames) > 0
	recommendations := []execution.RecipeStatusEvent{}

	if isTargetedInstall {
		for _, d := range availableRecipes {
			recipeName := d.Recipe.Name
			installed := i.status.RecipeHasStatus(recipeName, execution.RecipeStatusTypes.INSTALLED)
			failed := i.status.RecipeHasStatus(recipeName, execution.RecipeStatusTypes.FAILED)
			unsupported := i.status.RecipeHasStatus(recipeName, execution.RecipeStatusTypes.UNSUPPORTED)

			if installed || failed || unsupported {
				continue
			}
			e := execution.RecipeStatusEvent{Recipe: *d.Recipe, ValidationDurationMs: d.DurationMs}
			recommendations = append(recommendations, e)
		}
	}
	return recommendations
}

func (i *RecipeInstall) isTargetInstallRecipe(recipeName string) bool {
	for _, r := range i.RecipeNames {
		if r == recipeName {
			return true
		}
	}
	return false
}

func (i *RecipeInstall) installAdditionalBundle(bundler RecipeBundler, bundleInstaller RecipeBundleInstaller, repo *recipes.RecipeRepository) error {

	var additionalBundle *recipes.Bundle
	if i.RecipeNamesProvided() {
		additionalBundle = bundler.CreateAdditionalTargetedBundle(i.RecipeNames)
		i.reportUnsupportedTargetedRecipes(additionalBundle, repo)
		log.Debugf("Additional Targeted bundle recipes:%s", additionalBundle)
	} else {
		additionalBundle = bundler.CreateAdditionalGuidedBundle()
		log.Debugf("Additional Guided bundle recipes:%s", additionalBundle)
	}

	bundleInstaller.InstallContinueOnError(additionalBundle, i.AssumeYes)

	if bundleInstaller.InstalledRecipesCount() == 0 {
		return &types.UncaughtError{
			Err: fmt.Errorf("no recipes were installed"),
		}
	} else if len(i.RecipeNames) > len(additionalBundle.BundleRecipes) {
		return &types.UncaughtError{
			Err: fmt.Errorf("one or more selected recipes could not be installed"),
		}
	}

	return nil
}

func (i *RecipeInstall) installCoreBundle(bundler RecipeBundler, bundleInstaller RecipeBundleInstaller) error {

	if i.shouldInstallCore() {
		coreBundle := bundler.CreateCoreBundle()
		log.Debugf("Core bundle recipes:%s", coreBundle)
		err := bundleInstaller.InstallStopOnError(coreBundle, i.AssumeYes)
		if err != nil {
			log.Debugf("error installing core bundle:%s", err)
			return err
		}
	} else {
		log.Debugf("Skipping core bundle")
	}

	return nil
}

func (i *RecipeInstall) reportUnsupportedTargetedRecipes(bundle *recipes.Bundle, repo *recipes.RecipeRepository) {
	for _, recipeName := range i.RecipeNames {
		br := bundle.GetBundleRecipe(recipeName)
		if br == nil {
			recipe := repo.FindRecipeByName(recipeName)
			if recipe == nil {
				recipe = &types.OpenInstallationRecipe{Name: recipeName, DisplayName: recipeName}
			}
			unsupportedEvent := execution.NewRecipeStatusEvent(recipe)
			i.status.RecipeUnsupported(unsupportedEvent)
		} else {
			if !br.HasStatus(execution.RecipeStatusTypes.AVAILABLE) {
				ds := &recipes.DetectedStatusType{
					Status:     execution.RecipeStatusTypes.UNSUPPORTED,
					DurationMs: 0,
				}
				br.DetectedStatuses = append(br.DetectedStatuses, ds)
			}
		}
	}
}

// installing recipe
func (i *RecipeInstall) executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars, assumeYes bool) (string, error) {
	i.status.RecipeInstalling(execution.RecipeStatusEvent{Recipe: *r})

	// Execute the recipe steps.
	if err := i.recipeExecutor.Execute(ctx, *r, vars); err != nil {
		if err == types.ErrInterrupt {
			return "", err
		}

		if e, ok := err.(*types.UnsupportedOperatingSystemError); ok {
			i.status.RecipeUnsupported(execution.RecipeStatusEvent{
				Recipe:   *r,
				Msg:      e.Error(),
				Metadata: i.recipeExecutor.GetOutput().Metadata(),
			})

			return "", err
		}

		// Needed to move some of the failure CLI output/messaging here.  We want to capture metadata on when a user
		// opts in/out of sending cli logs, and this must occur before we build the RecipeStatusEvent to post to NR
		msg := fmt.Sprintf("execution failed for %s: %s", r.Name, err)
		i.optInToSendLogsAndUpdateRecipeMetadata(r.Name)
		se := execution.RecipeStatusEvent{
			Recipe:   *r,
			Msg:      msg,
			Metadata: i.recipeExecutor.GetOutput().Metadata(),
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

	entityGUID := i.recipeExecutor.GetOutput().EntityGUID()
	if entityGUID != "" {
		log.Debugf("Found entityGuid from recipe execution:%s", entityGUID)

		i.status.RecipeInstalled(execution.RecipeStatusEvent{
			Recipe:     *r,
			EntityGUID: entityGUID,
			Metadata:   i.recipeExecutor.GetOutput().Metadata(),
		})

		return entityGUID, nil
	}

	// show validation spinner if we need to validate and has no other spinner (Spinner is show when assume yes)
	if !assumeYes {
		msg := fmt.Sprintf("Validating %s", r.DisplayName)
		i.progressIndicator.ShowSpinner(!assumeYes)
		i.progressIndicator.Start(msg)
	}

	validationStart := time.Now()
	entityGUID, err := i.validateRecipeViaAllMethods(ctx, r, m, vars, assumeYes)
	validationDurationMs := time.Since(validationStart).Milliseconds()
	if err != nil {
		validationErr := fmt.Errorf("encountered an error while validating receipt of data for %s: %w", r.Name, err)
		i.status.RecipeFailed(execution.RecipeStatusEvent{
			Recipe:               *r,
			Msg:                  validationErr.Error(),
			ValidationDurationMs: validationDurationMs,
			Metadata:             i.recipeExecutor.GetOutput().Metadata(),
		})

		return "", validationErr
	}

	i.status.RecipeInstalled(execution.RecipeStatusEvent{
		Recipe:               *r,
		EntityGUID:           entityGUID,
		ValidationDurationMs: validationDurationMs,
		Metadata:             i.recipeExecutor.GetOutput().Metadata(),
	})

	return entityGUID, nil
}

func (i *RecipeInstall) optInToSendLogsAndUpdateRecipeMetadata(recipeName string) {
	i.progressIndicator.Fail("Installing " + recipeName)
	recipeOutput := i.recipeExecutor.GetRecipeOutput()
	logCaptureEnabledForRecipe := i.recipeExecutor.GetOutput().IsCapturedCliOutput()
	if len(recipeOutput) > 0 && logCaptureEnabledForRecipe {
		userOptIn := i.recipeLogForwarder.PromptUserToSendLogs(os.Stdin)
		i.recipeLogForwarder.SetUserOptedIn(userOptIn)
		i.recipeExecutor.GetOutput().AddMetadata("SendLogsOptIn", strconv.FormatBool(userOptIn))
	}
}

type validationFunc func() (string, error)

// Post install validation
func (i *RecipeInstall) validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars, assumeYes bool) (string, error) {
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
		log.Debugf("no validationUrl defined, skipping")
	}

	// Add NRQL validation if configured
	if r.ValidationNRQL != "" {
		validationFuncs = append(validationFuncs, func() (string, error) {
			return i.recipeValidator.ValidateRecipe(timeoutCtx, *m, *r, vars)
		})
	} else {
		log.Debugf("no validationNRQL defined, skipping")
	}

	if len(validationFuncs) == 0 {
		return "", nil
	}

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
			return entityGUID, nil
		case err := <-validationErrorChan:
			validationErrors = append(validationErrors, err)
			log.Debugf("validation error encountered: %s", err)

			if len(validationErrors) == len(validationFuncs) {
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
			return
		}

		vars, err := i.recipeVarPreparer.Prepare(*m, *r, assumeYes, licenseKey)
		if err != nil {
			errorChan <- err
			return
		}

		vars["assumeYes"] = fmt.Sprintf("%v", assumeYes)

		entityGUID, err := i.executeAndValidate(ctx, m, r, vars, assumeYes)

		if err != nil {
			errorChan <- err
			return
		}
		successChan <- entityGUID
	}()

	i.progressIndicator.ShowSpinner(assumeYes)
	i.progressIndicator.Start(msg)

	for {
		select {
		case entityGUID := <-successChan:
			i.progressIndicator.Success("Installing " + r.DisplayName)

			return entityGUID, nil
		case err := <-errorChan:
			if errors.Is(err, types.ErrInterrupt) {
				i.progressIndicator.Canceled("Installing " + r.DisplayName)
			} else {
				// progressIndicator has already been called; we need to finish i.e. message about logs being sent
				// and actually post logs to NR if the user has opted-in
				i.finishHandlingFailure(r.DisplayName)
			}
			log.Debugf("install error encountered: %s", err)
			return "", err
		}
	}
}

func (i *RecipeInstall) finishHandlingFailure(recipeName string) {
	if i.recipeLogForwarder.HasUserOptedIn() {
		i.progressIndicator.Start("Sending logs to New Relic")
		i.recipeLogForwarder.SendLogsToNewRelic(recipeName, i.recipeExecutor.GetRecipeOutput())
		i.progressIndicator.Success("Complete!")
	}
}

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
