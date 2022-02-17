package install

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type StatusReporter interface {
	ReportStatus(status execution.RecipeStatusType, recipe types.OpenInstallationRecipe)
}

type BundleInstaller struct {
	installedRecipes map[string]bool
	ctx              context.Context
	manifest         *types.DiscoveryManifest
	statusReporter   StatusReporter
	recipeInstaller  *RecipeInstaller
}

func NewBundleInstaller(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstaller *RecipeInstaller, statusReporter StatusReporter) *BundleInstaller {

	return &BundleInstaller{
		ctx:              ctx,
		manifest:         manifest,
		recipeInstaller:  recipeInstaller,
		statusReporter:   statusReporter,
		installedRecipes: make(map[string]bool),
	}
}

func (bi *BundleInstaller) InstallStopOnError(bundle *recipes.Bundle, assumeYes bool) error {

	bi.reportStatus(bundle)

	for _, br := range bundle.BundleRecipes {
		err := bi.installBundleRecipe(br, assumeYes)

		if err != nil {
			return err
		}
	}

	return nil
}

func (bi *BundleInstaller) InstallContinueOnError(bundle *recipes.Bundle, assumeYes bool) {

	bi.reportStatus(bundle)
	//TODO does this need to `reportStatus` like InstallStopOnError?
	// bi.reportStatus(bundle)

	for _, br := range bundle.BundleRecipes {
		err := bi.installBundleRecipe(br, assumeYes)
		log.Debugf("error installing recipe %v: %v", br.Recipe.Name, err)
	}
}

func (bi *BundleInstaller) reportStatus(bundle *recipes.Bundle) {

	for _, recipe := range bundle.BundleRecipes {
		for _, status := range recipe.RecipeStatuses {
			//TODO: should we keep track if status is already reported/saved
			bi.statusReporter.ReportStatus(status.Status, *recipe.Recipe)
		}
	}
}

func (bi *BundleInstaller) installBundleRecipe(bundleRecipe *recipes.BundleRecipe, assumeYes bool) error {

	// no dependencies
	var err error

	for _, dr := range bundleRecipe.Dependencies {
		err = bi.installBundleRecipe(dr, assumeYes)
		if err != nil {
			return err
		}
	}

	var withAvailableToInstallStatus = bundleRecipe.HasStatus(execution.RecipeStatusTypes.AVAILABLE)

	if _, found := bi.installedRecipes[bundleRecipe.Recipe.Name]; !found && withAvailableToInstallStatus {
		recipeName := bundleRecipe.Recipe.Name
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
	}

	//TODO: actual install here
	return nil
}

// Installer bundle no prompting
// Error handling with core bundle, additional
// TODO: Log Match needs to be reviewed, needs to log match process on the host
// TODO: maybe log match dont need detection, just look for all logs
