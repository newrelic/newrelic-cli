package install

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/http/httpproxy"

	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/packs"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

type RecipeInstaller struct {
	types.InstallerContext
	discoverer        Discoverer
	fileFilterer      FileFilterer
	manifestValidator *discovery.ManifestValidator
	recipeFetcher     recipes.RecipeFetcher
	recipeExecutor    execution.RecipeExecutor
	recipeValidator   RecipeValidator
	recipeFileFetcher RecipeFileFetcher
	status            *execution.InstallStatus
	prompter          Prompter
	progressIndicator ux.ProgressIndicator
	licenseKeyFetcher LicenseKeyFetcher
	configValidator   ConfigValidator
	recipeVarPreparer RecipeVarPreparer
	recipeFilterer    RecipeFilterRunner
	packsFetcher      PacksFetcher
	packsInstaller    PacksInstaller
}

type RecipeInstallFunc func(ctx context.Context, i *RecipeInstaller, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) error

var (
	recipeInstallFuncs map[string]RecipeInstallFunc = map[string]RecipeInstallFunc{
		"logs-integration": installLogging,
	}
)

func NewRecipeInstaller(ic types.InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstaller {
	checkNetwork(nrClient)

	var recipeFetcher recipes.RecipeFetcher

	if ic.LocalRecipes != "" {
		recipeFetcher = &recipes.LocalRecipeFetcher{
			Path: ic.LocalRecipes,
		}
	} else {
		recipeFetcher = recipes.NewServiceRecipeFetcher(&nrClient.NerdGraph)
	}

	mv := discovery.NewManifestValidator()
	ff := recipes.NewRecipeFileFetcher()
	ers := []execution.StatusSubscriber{
		execution.NewInstallEventsReporter(&nrClient.InstallEvents),
		execution.NewNerdStorageStatusReporter(&nrClient.NerdStorage),
		execution.NewTerminalStatusReporter(),
	}
	lkf := NewServiceLicenseKeyFetcher(&nrClient.NerdGraph)
	slg := execution.NewPlatformLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)

	d := discovery.NewPSUtilDiscoverer()
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	cv := diagnose.NewConfigValidator(nrClient)
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()
	rvp := execution.NewRecipeVarProvider()
	rf := recipes.NewRecipeFilterRunner(ic, statusRollup)
	spf := packs.NewServicePacksFetcher(&nrClient.NerdGraph, statusRollup)
	cpi := packs.NewServicePacksInstaller(nrClient, statusRollup)

	i := RecipeInstaller{
		discoverer:        d,
		fileFilterer:      gff,
		manifestValidator: mv,
		recipeFetcher:     recipeFetcher,
		recipeExecutor:    re,
		recipeValidator:   v,
		recipeFileFetcher: ff,
		status:            statusRollup,
		prompter:          p,
		progressIndicator: pi,
		licenseKeyFetcher: lkf,
		configValidator:   cv,
		recipeVarPreparer: rvp,
		recipeFilterer:    rf,
		packsFetcher:      spf,
		packsInstaller:    cpi,
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
	fmt.Println()

	log.Tracef("InstallerContext: %+v", i.InstallerContext)
	log.WithFields(log.Fields{
		"ShouldInstallInfraAgent":   i.ShouldInstallInfraAgent(),
		"ShouldInstallLogging":      i.ShouldInstallLogging(),
		"ShouldInstallIntegrations": i.ShouldInstallIntegrations(),
		"RecipesProvided":           i.RecipesProvided(),
		"RecipePathsProvided":       i.RecipePathsProvided(),
		"RecipeNamesProvided":       i.RecipeNamesProvided(),
	}).Debug("context summary")

	ctx, cancel := context.WithCancel(utils.SignalCtx)
	defer cancel()

	errChan := make(chan error)
	var err error

	// log.Printf("Validating connectivity to the New Relic platform...")
	// if err = i.configValidator.Validate(utils.SignalCtx); err != nil {
	// 	return err
	// }

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
	// Execute the discovery process, exiting on failure.
	m, err := i.discover(ctx)
	if err != nil {
		return err
	}

	if err = i.assertDiscoveryValid(ctx, m); err != nil {
		i.status.DiscoveryComplete(*m)
		return err
	}

	i.status.DiscoveryComplete(*m)

	repo := recipes.NewRecipeRepository(func() ([]types.OpenInstallationRecipe, error) {
		recipes, err2 := i.recipeFetcher.FetchRecipes(ctx)
		return recipes, err2
	})

	recipesForPlatform, err := repo.FindAll(*m)
	if err != nil {
		log.Debugf("should throw here %s", err)
		return err
	}
	log.Tracef("recipes found for platform: %v\n", recipesForPlatform)

	var targetedRecipes = recipesForPlatform
	if i.RecipesProvided() {
		targetedRecipes, err = i.fetchProvidedRecipe(m, recipesForPlatform)
		if err != nil {
			return err
		}
		log.Tracef("recipes supplied by user: %v\n", targetedRecipes)
	}

	filteredRecipes := i.recipeFilterer.RunFilterMultiple(ctx, targetedRecipes, m)
	log.Tracef("recipes after filtering: %v\n", filteredRecipes)

	if !i.RecipesProvided() {
		var selected, unselected []types.OpenInstallationRecipe
		selected, unselected, err = i.promptUserSelect(filteredRecipes)
		if err != nil {
			return err
		}
		log.Tracef("recipes selected by user: %v\n", selected)

		for _, r := range unselected {
			i.status.RecipeSkipped(execution.RecipeStatusEvent{Recipe: r})
		}

		filteredRecipes = selected
	}

	i.status.RecipesSelected(filteredRecipes)

	dependencies := resolveDependencies(filteredRecipes, recipesForPlatform)
	recipesToInstall := addIfMissing(filteredRecipes, dependencies)

	if err = i.installRecipes(ctx, m, recipesToInstall); err != nil {
		return err
	}

	if err = i.fetchAndInstallPacks(ctx, recipesToInstall); err != nil {
		return err
	}

	log.Debugf("Done installing.")
	return nil
}

func (i *RecipeInstaller) fetchAndInstallPacks(ctx context.Context, recipesToInstall []types.OpenInstallationRecipe) error {
	packs, err := i.packsFetcher.FetchPacks(ctx, recipesToInstall)
	if err != nil {
		// nolint: golint
		return fmt.Errorf("Failed to fetch observability packs: %s", err)
	}
	log.Debugf("Fetched Packs: %d", len(packs))

	if len(packs) > 0 {
		if err := i.packsInstaller.Install(ctx, packs); err != nil {
			// nolint: golint
			return fmt.Errorf("Failed to install observability pack: %s", err)
		}
	}

	return nil
}

func (i *RecipeInstaller) assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error {
	err := i.manifestValidator.Validate(m)
	if err != nil {
		return err
	}
	log.Debugf("Done asserting valid operating system for OS:%s and PlatformVersion:%s", m.OS, m.PlatformVersion)
	return nil
}

func (i *RecipeInstaller) installRecipes(ctx context.Context, m *types.DiscoveryManifest, recipes []types.OpenInstallationRecipe) error {
	log.WithFields(log.Fields{
		"recipe_count": len(recipes),
	}).Debug("installing recipes")
	var lastError error

	for _, r := range recipes {
		var err error

		log.WithFields(log.Fields{
			"name": r.Name,
		}).Debug("installing recipe")

		if f, ok := recipeInstallFuncs[r.Name]; ok {
			err = f(ctx, i, m, &r, recipes)
		} else {
			_, err = i.executeAndValidateWithProgress(ctx, m, &r)
		}

		if err != nil {
			if err == types.ErrInterrupt {
				return err
			}

			if r.Name == types.InfraAgentRecipeName || r.Name == types.LoggingRecipeName || i.RecipesProvided() {
				return err
			}

			lastError = err

			log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", r.Name, err)
			log.Warn(err)
			log.Warn(i.failMessage(r.DisplayName))
		}
		log.Debugf("Done executing and validating with progress for recipe name %s.", r.Name)
	}

	if lastError != nil {
		// Return last recipe error that was caught if any
		return lastError
	}

	if !i.status.WasSuccessful() {
		return fmt.Errorf("no recipes were installed")
	}

	return nil
}

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

	var entityGUID string
	var err error
	var validationDurationMilliseconds int64
	start := time.Now()
	if r.ValidationNRQL != "" {
		entityGUID, err = i.recipeValidator.ValidateRecipe(ctx, *m, *r)
		if err != nil {
			validationDurationMilliseconds = time.Since(start).Milliseconds()
			msg := fmt.Sprintf("encountered an error while validating receipt of data for %s: %s", r.Name, err)
			i.status.RecipeFailed(execution.RecipeStatusEvent{
				Recipe:                         *r,
				Msg:                            msg,
				ValidationDurationMilliseconds: validationDurationMilliseconds,
			})
			return "", errors.New(msg)
		}
	} else {
		log.Debugf("skipping validation due to missing validation query")
	}

	validationDurationMilliseconds = time.Since(start).Milliseconds()
	i.status.RecipeInstalled(execution.RecipeStatusEvent{
		Recipe:                         *r,
		EntityGUID:                     entityGUID,
		ValidationDurationMilliseconds: validationDurationMilliseconds,
	})

	return entityGUID, nil
}

func (i *RecipeInstaller) executeAndValidateWithProgress(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe) (string, error) {
	msg := fmt.Sprintf("Installing %s", r.Name)
	i.progressIndicator.Start(msg)
	defer func() { i.progressIndicator.Stop() }()

	if r.PreInstallMessage() != "" {
		fmt.Println(r.PreInstallMessage())
	}

	licenseKey, err := i.licenseKeyFetcher.FetchLicenseKey(ctx)
	if err != nil {
		return "", err
	}

	vars, err := i.recipeVarPreparer.Prepare(*m, *r, i.AssumeYes, licenseKey)
	if err != nil {
		return "", err
	}

	entityGUID, err := i.executeAndValidate(ctx, m, r, vars)
	if err != nil {
		i.progressIndicator.Fail(msg)
		return "", err
	}

	if r.PostInstallMessage() != "" {
		fmt.Println(r.PostInstallMessage())
	}

	i.progressIndicator.Success(msg)
	return entityGUID, nil
}

func (i *RecipeInstaller) failMessage(componentName string) error {
	searchURL := "https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/troubleshooting/not-seeing-data/"

	return fmt.Errorf("execution of %s failed, please see the following link for clues on how to resolve the issue: %s", componentName, searchURL)
}

func (i *RecipeInstaller) fetchProvidedRecipe(m *types.DiscoveryManifest, recipesForPlatform []types.OpenInstallationRecipe) ([]types.OpenInstallationRecipe, error) {
	var recipes []types.OpenInstallationRecipe

	// Load the recipes from the provided file names.
	for _, n := range i.RecipePaths {
		log.Debugln(fmt.Sprintf("Attempting to match recipePath %s.", n))
		recipe, err := i.recipeFromPath(n)
		if err != nil {
			log.Debugln(fmt.Sprintf("Error while building recipe from path, detail:%s.", err))
			return nil, err
		}

		log.WithFields(log.Fields{
			"name":         recipe.Name,
			"display_name": recipe.DisplayName,
			"path":         n,
		}).Debug("found recipe at path")

		recipes = append(recipes, *recipe)
	}

	// Load the recipes from the provided file names.
	for _, n := range i.RecipeNames {
		found := false
		log.Debugln(fmt.Sprintf("Attempting to match recipe name %s.", n))
		for _, r := range recipesForPlatform {
			if strings.EqualFold(r.Name, n) {
				log.WithFields(log.Fields{
					"name":         r.Name,
					"display_name": r.DisplayName,
				}).Debug("found recipe with name")
				recipes = append(recipes, r)
				found = true
				break
			}
		}
		if !found {
			log.Errorf("Could not find recipe with name %s.", n)
		}
	}

	return recipes, nil
}

func (i *RecipeInstaller) recipeFromPath(recipePath string) (*types.OpenInstallationRecipe, error) {
	recipeURL, parseErr := url.Parse(recipePath)
	if parseErr == nil && recipeURL.Scheme != "" {
		f, err := i.recipeFileFetcher.FetchRecipeFile(recipeURL)
		if err != nil {
			return nil, fmt.Errorf("could not fetch file %s: %s", recipePath, err)
		}
		return f, nil
	}

	f, err := i.recipeFileFetcher.LoadRecipeFile(recipePath)
	if err != nil {
		return nil, fmt.Errorf("could not load file %s: %s", recipePath, err)
	}

	return f, nil
}

func (i *RecipeInstaller) promptUserSelect(recipes []types.OpenInstallationRecipe) ([]types.OpenInstallationRecipe, []types.OpenInstallationRecipe, error) {
	if len(recipes) == 0 {
		return []types.OpenInstallationRecipe{}, []types.OpenInstallationRecipe{}, nil
	}

	if i.AssumeYes {
		return recipes, []types.OpenInstallationRecipe{}, nil
	}

	var selectedRecipes, unselectedRecipes []types.OpenInstallationRecipe

	names := []string{}
	selected := []string{}
	for _, r := range recipes {
		if r.Name != types.InfraAgentRecipeName {
			names = append(names, r.DisplayName)
		} else {
			fmt.Printf("The guided installation will begin by installing the latest version of the New Relic Infrastructure agent, which is required for additional instrumentation.\n\n")
		}
	}

	if len(names) > 0 {
		var promptErr error
		selected, promptErr = i.prompter.MultiSelect("Please choose from the following instrumentation to be installed:", names)
		if promptErr != nil {
			return nil, nil, promptErr
		}
		fmt.Println()
	}

	for _, r := range recipes {
		if utils.StringInSlice(r.DisplayName, selected) || r.Name == types.InfraAgentRecipeName {
			selectedRecipes = append(selectedRecipes, r)
		} else {
			unselectedRecipes = append(unselectedRecipes, r)
		}
	}

	return selectedRecipes, unselectedRecipes, nil
}

func findRecipeInRecipes(name string, recipes []types.OpenInstallationRecipe) *types.OpenInstallationRecipe {
	for _, r := range recipes {
		if r.Name == name {
			return &r
		}
	}

	return nil
}

// This is a naive implementation that only resolves dependencies one level deep.
func resolveDependencies(recipes []types.OpenInstallationRecipe, recipesForPlatform []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
	var results []types.OpenInstallationRecipe

	for _, r := range recipes {
		if len(r.Dependencies) > 0 {
			for _, n := range r.Dependencies {
				d := findRecipeInRecipes(n, recipesForPlatform)
				if d != nil {
					results = append(results, *d)
				}
			}
		}
	}

	return results
}

// This is a naive implementation that only resolves dependencies one level deep.
func addIfMissing(recipes []types.OpenInstallationRecipe, dependencies []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
	var results []types.OpenInstallationRecipe

	for _, d := range dependencies {
		if found := findRecipeInRecipes(d.Name, results); found == nil {
			results = append(results, d)
		}
	}

	for _, r := range recipes {
		if found := findRecipeInRecipes(r.Name, results); found == nil {
			results = append(results, r)
		}
	}

	return results
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
