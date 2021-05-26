package install

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

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
	InstallerContext
	discoverer        discovery.Discoverer
	fileFilterer      discovery.FileFilterer
	manifestValidator *discovery.ManifestValidator
	recipeFetcher     recipes.RecipeFetcher
	recipeExecutor    execution.RecipeExecutor
	recipeValidator   validation.RecipeValidator
	recipeFileFetcher recipes.RecipeFileFetcher
	status            *execution.InstallStatus
	prompter          ux.Prompter
	progressIndicator ux.ProgressIndicator
	licenseKeyFetcher LicenseKeyFetcher
	configValidator   diagnose.ConfigValidator
	recipeVarProvider execution.RecipeVarProvider
}

func NewRecipeInstaller(ic InstallerContext, nrClient *newrelic.NewRelic) *RecipeInstaller {

	var recipeFetcher recipes.RecipeFetcher

	if ic.LocalRecipes != "" {
		recipeFetcher = &recipes.LocalRecipeFetcher{
			Path: ic.LocalRecipes,
		}

	} else {
		recipeFetcher = recipes.NewServiceRecipeFetcher(&nrClient.NerdGraph)
	}

	pf := discovery.NewRegexProcessFilterer(recipeFetcher)
	mv := discovery.NewManifestValidator()
	ff := recipes.NewRecipeFileFetcher()
	ers := []execution.StatusSubscriber{
		execution.NewNerdStorageStatusReporter(&nrClient.NerdStorage),
		execution.NewTerminalStatusReporter(),
	}
	lkf := NewServiceLicenseKeyFetcher(&nrClient.NerdGraph)
	slg := execution.NewConcreteSuccessLinkGenerator()
	statusRollup := execution.NewInstallStatus(ers, slg)

	d := discovery.NewPSUtilDiscoverer(pf)
	gff := discovery.NewGlobFileFilterer()
	re := execution.NewGoTaskRecipeExecutor()
	v := validation.NewPollingRecipeValidator(&nrClient.Nrdb)
	cv := diagnose.NewConcreteConfigValidator(nrClient)
	p := ux.NewPromptUIPrompter()
	pi := ux.NewPlainProgress()
	rvp := execution.NewConcreteRecipeVarProvider()

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
		recipeVarProvider: rvp,
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
		"ShouldRunDiscovery":        i.ShouldRunDiscovery(),
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

	log.Printf("Validating connectivity to the New Relic platform...")
	if err = i.configValidator.ValidateConfig(ctx); err != nil {
		return err
	}

	go func(ctx context.Context) {
		errChan <- i.discoverAndRun(ctx)
	}(ctx)

	select {
	case <-ctx.Done():
		i.status.InstallCanceled()
		return nil
	case err = <-errChan:
		if err == types.ErrInterrupt {
			i.status.InstallCanceled()
			return err
		}

		i.status.InstallComplete(err)

		return err
	}
}

func (i *RecipeInstaller) discoverAndRun(ctx context.Context) error {
	// Execute the discovery process, exiting on failure.
	m, err := i.discover(ctx)
	if err != nil {
		return err
	}

	i.status.DiscoveryComplete(*m)

	err = i.assertDiscoveryValid(ctx, m)
	if err != nil {
		return err
	}

	if i.RecipesProvided() {
		// Run the targeted (AKA stitched path) installer.
		return i.targetedInstall(ctx, m)
	}

	// Run the guided installer.
	return i.guidedInstall(ctx, m)
}

func (i *RecipeInstaller) assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error {
	err := i.manifestValidator.Execute(m)
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

	for _, r := range recipes {
		var err error

		log.WithFields(log.Fields{
			"name": r.Name,
		}).Debug("installing recipe")

		_, err = i.executeAndValidateWithProgress(ctx, m, &r)
		if err != nil {
			if err == types.ErrInterrupt {
				return err
			}

			log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", r.Name, err)
			log.Warn(err)
			log.Warn(i.failMessage(r.DisplayName))

			if len(recipes) == 1 {
				return err
			}
		}
		log.Debugf("Done executing and validating with progress for recipe name %s.", r.Name)
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

	vars, err := i.recipeVarProvider.Prepare(*m, *r, i.AssumeYes, licenseKey)
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

func (i *RecipeInstaller) fetchRecipeAndReportAvailable(ctx context.Context, m *types.DiscoveryManifest, recipeName string) (*types.OpenInstallationRecipe, error) {
	log.WithFields(log.Fields{
		"name": recipeName,
	}).Debug("fetching recipe for install")

	r, err := i.fetch(ctx, m, recipeName)
	if err != nil {
		return nil, err
	}

	i.status.RecipeAvailable(*r)

	return r, nil
}

func (i *RecipeInstaller) fetch(ctx context.Context, m *types.DiscoveryManifest, recipeName string) (*types.OpenInstallationRecipe, error) {
	r, err := i.recipeFetcher.FetchRecipe(ctx, m, recipeName)
	if err != nil {
		log.Errorf("error retrieving recipe %s: %s", recipeName, err)
		return nil, err
	}

	if r == nil {
		return nil, fmt.Errorf("recipe %s not found", recipeName)
	}

	return r, nil
}
