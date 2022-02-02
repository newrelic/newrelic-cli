package recipes

import (
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type bundle struct {
	name        string
	recipeNames []string
	recipes     []types.OpenInstallationRecipe
}

var coreBundle = bundle{
	name: "Core",
	recipeNames: []string{
		types.InfraAgentRecipeName,
		types.LoggingRecipeName,
		types.GoldenRecipeName,
	},
	recipes: []types.OpenInstallationRecipe{},
}

var additionalBundle = bundle{
	name:    "Other",
	recipes: []types.OpenInstallationRecipe{},
}

// extracts recipe for bundle, and return the recipes without the extracted recipes
func (b *bundle) create(recipesForInstall []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {
	for _, n := range b.recipeNames {
		for i, r := range recipesForInstall {
			if strings.EqualFold(r.Name, n) {
				b.recipes = append(b.recipes, r)
				recipesForInstall = append(recipesForInstall[:i], recipesForInstall[i+1:]...)
				break
			}
		}
	}

	return recipesForInstall
}

func (b *bundle) any() bool {
	return len(b.recipes) > 0
}

type bundler struct {
	bundles []*bundle
}

func (b *bundler) install() error {

	return nil
}

func newBundler() *bundler {

	b := &bundler{
		bundles: []*bundle{
			&coreBundle,
			&additionalBundle,
		}}

	return b
}

/*

Bundler b = bundler {}
bundler.install(){

}
*/
