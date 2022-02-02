package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type bundle struct {
	name string
	// use map of recipe names to combine recipe
	recipes      map[string]*types.OpenInstallationRecipe
	shouldPrompt bool
}

//type coreBundle []types.OpenInstallationRecipe
//type additionalBundle []types.OpenInstallationRecipe

var coreBundle = bundle{
	name: "Core",
	recipes: map[string]*types.OpenInstallationRecipe{
		types.InfraAgentRecipeName: nil,
		types.LoggingRecipeName:    nil,
		types.GoldenRecipeName:     nil,
	},
}

var additionalBundle = bundle{
	name:         "Additional",
	recipes:      map[string]*types.OpenInstallationRecipe{},
	shouldPrompt: true,
}

// create recipe for bundle, and return the recipes without the extracted recipes
func (b *bundle) create(recipesForInstall []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {

	// if bundle has recipe definition, we extract, otherwise, we assume it's additional bundle take all
	if b.any() {
		for i, r := range recipesForInstall {
			if _, ok := b.recipes[r.Name]; ok {
				b.recipes[r.Name] = &r
				recipesForInstall = append(recipesForInstall[:i], recipesForInstall[i+1:]...)
				break
			}
		}

	} else {
		for i, r := range recipesForInstall {
			b.recipes[r.Name] = &r
			recipesForInstall = append(recipesForInstall[:i], recipesForInstall[i+1:]...)
			break
		}
	}

	return recipesForInstall
}

func (b *bundle) any() bool {
	return len(b.recipes) > 0
}

func (b *bundle) promptMessage() string {
	return "Would you like to enable monitoring for the following?"
}
