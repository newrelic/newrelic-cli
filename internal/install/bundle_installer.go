package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	log "github.com/sirupsen/logrus"
)

type BundleInstaller struct {
	installedRecipes map[string]bool
	ctx              context.Context
	manifest         *types.DiscoveryManifest
	recipeInstaller  *RecipeInstaller
}

var (
	customRecipeInstallFuncs = map[string]RecipeInstallFunc{
		"logs-integration": installLogging,
	}
)

func NewBundleInstaller(recipeInstaller *RecipeInstaller) *BundleInstaller {

	return &BundleInstaller{
		recipeInstaller: recipeInstaller,
	}
}

func (bi *BundleInstaller) InstallStopOnError(bundle *recipes.Bundle) error {

	bi.ReportStatus(bundle)

	for _, br := range bundle.BundleRecipes {
		err := bi.InstallBundleRecipe(br)

		if err != nil {
			return err
		}
	}

	return nil
}

func (bi *BundleInstaller) ReportStatus(bundle *recipes.Bundle) {

	for _, recipe := range bundle.BundleRecipes {
		for _, status := range recipe.Statuses {
			bi.recipeInstaller.status.ReportStatus(status, *recipe.Recipe)
		}
	}
}

func (bi *BundleInstaller) InstallContinueOnError(bundle *recipes.Bundle) {

	for _, br := range bundle.BundleRecipes {
		bi.InstallBundleRecipe(br)
	}
}

func (bi *BundleInstaller) InstallBundleRecipe(bundleRecipe *recipes.BundleRecipe) error {

	// no dependencies
	var err error

	if len(bundleRecipe.Dependencies) == 0 {
		if _, ok := bi.installedRecipes[bundleRecipe.Recipe.Name]; !ok {
			recipeName := bundleRecipe.Recipe.Name
			bi.installedRecipes[recipeName] = true

			log.WithFields(log.Fields{
				"name": recipeName,
			}).Debug("installing recipe")

			_, err = bi.recipeInstaller.executeAndValidateWithProgress(bi.ctx, bi.manifest, bundleRecipe.Recipe)

			if err != nil {
				log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", recipeName, err)
				return err
			}
			log.Debugf("Done executing and validating with progress for recipe name %s.", recipeName)
		}
	}

	for _, dr := range bundleRecipe.Dependencies {
		err = bi.InstallBundleRecipe(dr)
		if err != nil {
			return err
		}
	}

	//TODO: actual install here
	return nil
}

// Installer bundle no prompting
// Error handling with core bundle, addtional
// TODO: Log Match needs to be reviewed, needs to log match process on the host
// TODO: maybe log match dont need detection, just look for all logs
