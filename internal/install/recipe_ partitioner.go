package install

import (
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
)

type recipePartition struct {
	name           string
	description    string
	recipeNames    []string
	recipes        []types.OpenInstallationRecipe
	requireConfirm bool
	prompter       ux.PromptUIPrompter
}

func (rp *recipePartition) partition(recipesForInstall []types.OpenInstallationRecipe) []types.OpenInstallationRecipe {

	for _, n := range rp.recipeNames {
		for i, r := range recipesForInstall {
			if strings.EqualFold(r.Name, n) {
				rp.recipes = append(rp.recipes, r)
				recipesForInstall = append(recipesForInstall[:i], recipesForInstall[i+1:]...)
				break
			}
		}
	}

	return recipesForInstall
}

func (rp *recipePartition) any() bool {
	return len(rp.recipes) > 0
}

func (rp recipePartition) String() string {

	var recipeNames string
	for _, recipe := range rp.recipes {
		recipeNames += fmt.Sprintf("\n%s", recipe.DisplayName)
	}

	return fmt.Sprintf("\nNew Relic installing %s recipes: %s", rp.name, recipeNames)
}

var coreRecipePartition = recipePartition{
	name:        "Core",
	description: "This is the core partition",
	recipeNames: []string{
		types.InfraAgentRecipeName,
		types.LoggingRecipeName,
	},
	recipes:        make([]types.OpenInstallationRecipe, 0),
	requireConfirm: false,
}

var otherRecipePartition = recipePartition{
	name:           "",
	description:    "This is the non-core partition",
	recipeNames:    make([]string, 0),
	recipes:        make([]types.OpenInstallationRecipe, 0),
	requireConfirm: true,
}

type recipePartitions []recipePartition

func newRecipePartitions(recipesForInstall []types.OpenInstallationRecipe) *recipePartitions {
	partions := &[]recipePartition{
		coreRecipePartition,
		otherRecipePartition,
	}

	for _, partition := range *partions {
		if partition.name == otherRecipePartition.name {
			partition.recipes = recipesForInstall
		} else {
			recipesForInstall = partition.partition(recipesForInstall)
			partition.prompter = *ux.NewPromptUIPrompter()
		}
	}

	return (*recipePartitions)(partions)
}
