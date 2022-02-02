package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type bundler struct {
	bundles []*bundle
}

func newBundler() *bundler {
	b := &bundler{
		bundles: []*bundle{
			&coreBundle,
			&additionalBundle,
		}}

	return b
}

func (bl *bundler) create(platformRecipes []types.OpenInstallationRecipe,
	targetRecipes []types.OpenInstallationRecipe) {

	for _, bundle := range bl.bundles {
		platformRecipes = bundle.create(platformRecipes)
		// if there are targeted recipes, we override platform recipes
		for _, targetRecipe := range targetRecipes {
			bundle.recipes[targetRecipe.Name] = &targetRecipe
		}
	}
}
