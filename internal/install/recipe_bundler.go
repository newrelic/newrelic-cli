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
		if bundle.name == additionalBundle.name {
			// for additional bundle, we take either platform for guided, or target for targed
			if len(targetRecipes) == 0 {
				bundle.create(platformRecipes)
			} else {
				bundle.create(targetRecipes)
			}
			break // additional should always be last
		} else {
			platformRecipes = bundle.create(platformRecipes)
			// if there are targeted recipes, we override platform recipes
			for _, targetRecipe := range targetRecipes {
				bundle.recipes[targetRecipe.Name] = &targetRecipe
			}

		}
	}
}
