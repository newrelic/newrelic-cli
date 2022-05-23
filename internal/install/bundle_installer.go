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
	installedRecipes     map[string]bool
	installFailedRecipes map[string]error
	ctx                  context.Context
	manifest             *types.DiscoveryManifest
	statusReporter       StatusReporter
	recipeInstaller      RecipeInstaller
	prompter             Prompter
}

func NewBundleInstaller(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) *BundleInstaller {

	return &BundleInstaller{
		ctx:                  ctx,
		manifest:             manifest,
		recipeInstaller:      recipeInstallerInterface,
		statusReporter:       statusReporter,
		installedRecipes:     make(map[string]bool),
		installFailedRecipes: make(map[string]error),
		prompter:             NewPrompter(),
	}
}

func NewPrompter() *ux.PromptUIPrompter {
	return ux.NewPromptUIPrompter()
}

func (bi *BundleInstaller) InstallStopOnError(bundle *recipes.Bundle, assumeYes bool) error {

	bi.reportBundleStatus(bundle)

	installableBundleRecipes := bi.getInstallableBundleRecipes(bundle)
	if len(installableBundleRecipes) == 0 {
		return nil
	}

	for _, br := range installableBundleRecipes {
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
		fmt.Println("\nWe've detected additional monitoring that can be configured by installing the following:")

		for _, bundleRecipe := range installableBundleRecipes {
			fmt.Printf("  %s\n", bundleRecipe.Recipe.DisplayName)
		}

		fmt.Println()
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
			msg := fmt.Sprintf("execution failed for %s: %s", additionalRecipe.Recipe.Name, err)
			log.Debugf(msg)
		}
	}
}

func (bi *BundleInstaller) reportBundleStatus(bundle *recipes.Bundle) {
	for _, recipe := range bundle.BundleRecipes {
		if bi.installedRecipes[recipe.Recipe.Name] {
			continue
		}
		for _, ds := range recipe.DetectedStatuses {
			e := execution.RecipeStatusEvent{Recipe: *recipe.Recipe, ValidationDurationMs: ds.DurationMs}
			bi.statusReporter.ReportStatus(ds.Status, e)
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
			// When dependency fail, mark parent fail as well
			bi.installFailedRecipes[bundleRecipe.Recipe.Name] = err
			// For dependency failure, it's not handle with normal flow of error reporting, we will manuall report error
			msg := fmt.Sprintf("execution failed for %s: %s", bundleRecipe.Recipe.Name, err)
			e := execution.RecipeStatusEvent{
				Recipe: *bundleRecipe.Recipe,
				Msg:    msg,
			}
			bi.statusReporter.ReportStatus(execution.RecipeStatusTypes.FAILED, e)
			log.Debugf(msg)
			return err
		}
	}

	recipeName := bundleRecipe.Recipe.Name
	if bi.installedRecipes[bundleRecipe.Recipe.Name] {
		return nil
	}

	var ok bool
	if err, ok = bi.installFailedRecipes[bundleRecipe.Recipe.Name]; ok {
		return err
	}

	log.WithFields(log.Fields{
		"name": recipeName,
	}).Debug("installing recipe")

	_, err = bi.recipeInstaller.executeAndValidateWithProgress(bi.ctx, bi.manifest, bundleRecipe.Recipe, assumeYes)
	if err != nil {
		bi.installFailedRecipes[recipeName] = err
		log.Debugf("Failed while executing and validating with progress for recipe name %s, detail:%s", recipeName, err)
		return err
	}

	bi.installedRecipes[recipeName] = true
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
