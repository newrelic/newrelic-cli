package install

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

type bundleInstaller struct {
	context            context.Context
	recipeInstaller    *RecipeInstaller
	recipeForPlatoform []types.OpenInstallationRecipe
	manifest           *types.DiscoveryManifest
	prompter           ux.PromptUIPrompter
}

func newBundleInstaller(c context.Context, i *RecipeInstaller, m *types.DiscoveryManifest) *bundleInstaller {
	return &bundleInstaller{
		context:         c,
		recipeInstaller: i,
		manifest:        m,
		prompter:        *ux.NewPromptUIPrompter(),
	}
}

// for core bundle, we install regardless of targeted or guided
func (bi *bundleInstaller) installBundle(bundle *bundle) error {

	var bundleRecipes []types.OpenInstallationRecipe

	for _, value := range bundle.recipes {
		bundleRecipes = append(bundleRecipes, *value)
	}

	bundleRecipes = bi.recipeInstaller.recipeFilterer.RunFilterAll(bi.context, bundleRecipes, &bi.recipeInstaller.status.DiscoveryManifest)

	if !bundle.any() {
		return nil
	}

	if bundle.shouldPrompt {
		confirm, err := bi.prompter.PromptYesNo(bundle.promptMessage())
		if err != nil {
			panic(err)
		}

		if !confirm {
			for _, r := range bundleRecipes {
				bi.recipeInstaller.status.RecipeSkipped(execution.RecipeStatusEvent{Recipe: r})
				return nil
			}
		}
	}

	dependencies := resolveDependencies(bundleRecipes, bi.recipeForPlatoform)
	bundleRecipes = addIfMissing(bundleRecipes, dependencies)
	logRecipes(bundleRecipes)

	if err := bi.recipeInstaller.installRecipes(bi.context, bi.manifest, bundleRecipes); err != nil {
		return err
	}
	bundle.recipes = map[string]*types.OpenInstallationRecipe{}

	for _, r := range bundleRecipes {
		bundle.recipes[r.Name] = &r
	}

	return nil
}
