package install

import (
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type recipeGroup int

const (
	core recipeGroup = iota
	others
)

type recipeSelector struct {
	recipeGroupMap map[recipeGroup][]string
}

func (r *recipeSelector) getRecipeGroup(rg recipeGroup, dr []types.OpenInstallationRecipe) []string {

	return r.recipeGroupMap[rg]
}

func newRecipeSelector() *recipeSelector {
	r := recipeSelector{
		recipeGroupMap: make(map[recipeGroup][]string),
	}
	r.recipeGroupMap[core] = append(r.recipeGroupMap[core], types.InfraAgentRecipeName)
	r.recipeGroupMap[core] = append(r.recipeGroupMap[core], types.LoggingRecipeName)

	return &r
}
