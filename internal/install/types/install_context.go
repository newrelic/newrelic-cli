package types

import (
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

const (
	ApmKeyword = "Apm"
	DefaultTag = "nr_deployed_by:newrelic-cli"
)

// nolint: maligned
type InstallerContext struct {
	AssumeYes   bool
	RecipeNames []string
	RecipePaths []string
	// LocalRecipes is the path to a local recipe directory from which to load recipes.
	LocalRecipes string
	EntityTags   []string
}

func (i *InstallerContext) RecipePathsProvided() bool {
	return len(i.RecipePaths) > 0
}

func (i *InstallerContext) RecipeNamesProvided() bool {
	return len(i.RecipeNames) > 0
}

func (i *InstallerContext) GetEntityTags() ([]entities.TaggingTagInput, bool) {
	i.EntityTags = append(i.EntityTags, DefaultTag)
	t, err := utils.AssembleTagsInput(i.EntityTags)
	if err != nil || len(t) == 0 {
		return nil, false
	}

	return t, true
}
