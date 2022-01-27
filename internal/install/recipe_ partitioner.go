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

func (rp *recipePartition) getPromptMessage() string {

	rn := strings.Join(rp.recipeNames, ",")
	return "New Relic CLI has detected:  " + rn + ".\n Would you like to go ahead and install?"
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
		types.GoldenRecipeName,
	},
	recipes:        make([]types.OpenInstallationRecipe, 0),
	requireConfirm: false,
}

var otherRecipePartition = recipePartition{
	name:           "Other",
	description:    "This is the non-core partition",
	recipeNames:    make([]string, 0),
	recipes:        make([]types.OpenInstallationRecipe, 0),
	requireConfirm: true,
	prompter:       *ux.NewPromptUIPrompter(),
}

type recipePartitions []*recipePartition

func newRecipePartitions(recipesForInstall []types.OpenInstallationRecipe) *recipePartitions {
	partitions := []*recipePartition{
		&coreRecipePartition,
		&otherRecipePartition,
	}

	for _, partition := range partitions {
		if partition.name == otherRecipePartition.name {
			partition.recipes = append(partition.recipes, recipesForInstall...)
		} else {
			recipesForInstall = partition.partition(recipesForInstall)
		}
	}

	return (*recipePartitions)(&partitions)
}
