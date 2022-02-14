package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type BundleInstaller struct {
	installedRecipes []*types.OpenInstallationRecipe
}

func (bi *BundleInstaller) InstallNonPromptBundle(bundle *recipes.Bundle) {

	for _, br := range bundle.BundleRecipes {
		bi.InstallBundleRecipe(br)
	}
}

func (bi *BundleInstaller) InstallBundleRecipe(bundleRecipe *recipes.BundleRecipe) {

	for _, dr := range bundleRecipe.Dependencies {
		bi.InstallBundleRecipe(dr)
	}
	//TODO: actual install here
}
