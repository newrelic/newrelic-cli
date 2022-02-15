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
	//FIXME: we used to report selected at one time, now we have to do in each or bundles? might need to check if UI Can handle naturally
	//TODO: a genearl method for reporting status
	var err error

	if len(bundleRecipe.Dependencies) == 0 {
		if _, ok := bi.installedRecipes[bundleRecipe.Recipe.Name]; !ok {
			recipeName := bundleRecipe.Recipe.Name
			bi.installedRecipes[recipeName] = true

			log.WithFields(log.Fields{
				"name": recipeName,
			}).Debug("installing recipe")

			if f, ok := recipeInstallFuncs[recipeName]; ok {
				//FIXME: nil for recipes, how do we get all the recipe for manifest
				err = f(bi.ctx, bi.recipeInstaller, bi.manifest, bundleRecipe.Recipe, nil)
			} else {
				_, err = bi.recipeInstaller.executeAndValidateWithProgress(bi.ctx, bi.manifest, bundleRecipe.Recipe)
			}

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

// func (bi *BundleInstaller) InstallRecipe() {
// 	//TODO: we need to report count in piece meal? Not sure how this will work
// 	// log.WithFields(log.Fields{
// 	// 	"recipe_count": len(recipes),
// 	// }).Debug("installing recipes")
// 	var lastError error

// 	for _, r := range recipes {
// 		var err error

// 		log.WithFields(log.Fields{
// 			"name": r.Name,
// 		}).Debug("installing recipe")

// 		if f, ok := recipeInstallFuncs[r.Name]; ok {
// 			err = f(ctx, i, m, &r, recipes)
// 		} else {
// 			_, err = i.executeAndValidateWithProgress(ctx, m, &r)
// 		}

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

// Installer bundle no prompting
// Error handling with core bundle, addtional
// TODO: Log Match needs to be reviewed, needs to log match process on the host
// TODO: maybe log match dont need detection, just look for all logs
