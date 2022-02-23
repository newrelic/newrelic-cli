package install

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

type StatusReporter interface {
	ReportStatus(status execution.RecipeStatusType, event execution.RecipeStatusEvent)
}

type BundleInstaller struct {
	installedRecipes map[string]bool
	ctx              context.Context
	manifest         *types.DiscoveryManifest
	statusReporter   StatusReporter
	recipeInstaller  RecipeInstallerInterface
	prompter         Prompter
}

//TODO should we revert this interface extraction? Was changed in order to mock via testify...
type RecipeInstallerInterface interface {
	promptIfNotLatestCLIVersion(ctx context.Context) error
	Install() error
	install(ctx context.Context) error
	assertDiscoveryValid(ctx context.Context, m *types.DiscoveryManifest) error
	discover(ctx context.Context) (*types.DiscoveryManifest, error)
	executeAndValidate(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, vars types.RecipeVars) (string, error)
	validateRecipeViaAllMethods(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest, vars types.RecipeVars) (string, error)
	executeAndValidateWithProgress(ctx context.Context, m *types.DiscoveryManifest, r *types.OpenInstallationRecipe, assumeYes bool) (string, error)
}

func NewBundleInstaller(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstallerInterface, statusReporter StatusReporter) *BundleInstaller {

	return &BundleInstaller{
		ctx:              ctx,
		manifest:         manifest,
		recipeInstaller:  recipeInstallerInterface,
		statusReporter:   statusReporter,
		installedRecipes: make(map[string]bool),
		prompter:         NewPrompter(),
	}
}

func NewPrompter() *ux.PromptUIPrompter {
	return ux.NewPromptUIPrompter()
}

func (bi *BundleInstaller) InstallStopOnError(bundle *recipes.Bundle, assumeYes bool) error {

	bi.reportBundleStatus(bundle)

	for _, br := range bundle.BundleRecipes {
		err := bi.InstallBundleRecipe(br, assumeYes)

		if err != nil {
			return err
		}
	}

	return nil
}

func (bi *BundleInstaller) InstallContinueOnError(bundle *recipes.Bundle, assumeYes bool) {
	bi.reportBundleStatus(bundle)

	installableBundleRecipes := bi.getInstallableBundleRecipes(bundle)
	if len(installableBundleRecipes) == 0 {
		return
	}

	if !assumeYes && bundle.IsAdditionalGuided() {
		//TODO: Should this be log instead of fmt?
		fmt.Println("\nWe've detected additional monitoring that can be configured by installing the following:")

		for _, bundleRecipe := range installableBundleRecipes {
			fmt.Println(bundleRecipe.Recipe.DisplayName)
		}

		prompter := ux.NewPromptUIPrompter()
		msg := "Continue installing? "
		isConfirmed, err := prompter.PromptYesNo(msg)

		if err != nil {
			log.Debug(err)
			isConfirmed = false
		}

		if !isConfirmed {
			for _, additionalRecipe := range installableBundleRecipes {
				skippedEvent := execution.NewRecipeStatusEvent(additionalRecipe.Recipe)
				bi.statusReporter.ReportStatus(execution.RecipeStatusTypes.SKIPPED, skippedEvent)
			}
			return
		}
	}

	for _, additionalRecipe := range installableBundleRecipes {
		err := bi.InstallBundleRecipe(additionalRecipe, assumeYes)
		if err != nil {
			log.Debugf("error installing recipe %v: %v", additionalRecipe.Recipe.Name, err)
		}
	}
}

func (bi *BundleInstaller) reportBundleStatus(bundle *recipes.Bundle) {
	for _, recipe := range bundle.BundleRecipes {
		for _, status := range recipe.DetectedStatuses {
			e := execution.RecipeStatusEvent{Recipe: *recipe.Recipe}
			bi.statusReporter.ReportStatus(status, e)
		}
	}
}

func (bi *BundleInstaller) InstalledRecipesCount() int {
	return len(bi.installedRecipes)
}

func (bi *BundleInstaller) InstallBundleRecipe(bundleRecipe *recipes.BundleRecipe, assumeYes bool) error {
	var err error

	for _, dr := range bundleRecipe.Dependencies {
		err = bi.InstallBundleRecipe(dr, assumeYes)
		if err != nil {
			return err
		}
	}

	recipeName := bundleRecipe.Recipe.Name
	if bi.installedRecipes[bundleRecipe.Recipe.Name] {
		return nil
	}

	bi.installedRecipes[recipeName] = true
	log.WithFields(log.Fields{
		"name": recipeName,
	}).Debug("installing recipe")

	_, err = bi.recipeInstaller.executeAndValidateWithProgress(bi.ctx, bi.manifest, bundleRecipe.Recipe, assumeYes)
	if err != nil {
		log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", recipeName, err)
		return err
	}
	log.Debugf("Done executing and validating with progress for recipe name %s.", recipeName)

	return nil
}

func (bi *BundleInstaller) getInstallableBundleRecipes(bundle *recipes.Bundle) []*recipes.BundleRecipe {
	var bundleRecipes []*recipes.BundleRecipe

	for _, bundleRecipe := range bundle.BundleRecipes {
		if !bundleRecipe.HasStatus(execution.RecipeStatusTypes.AVAILABLE) {
			//Skip if not available
			continue
		}
		if !bi.installedRecipes[bundleRecipe.Recipe.Name] {
			bundleRecipes = append(bundleRecipes, bundleRecipe)
		}
	}

	return bundleRecipes
}
