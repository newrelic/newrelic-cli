package recipes

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeFilterer struct {
	filters []RecipeFilter
}

func NewRecipeFilterer() *RecipeFilterer {
	return &RecipeFilterer{
		filters: []RecipeFilter{
			NewProcessMatchRecipeFilter(),
			NewScriptEvaluationRecipeFilter(),
		},
	}
}

func (rf *RecipeFilterer) Filter(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	for _, f := range rf.filters {
		filtered, err := f.Execute(ctx, r, m)
		if err != nil {
			return false, err
		}

		if filtered {
			log.Debugf("Filtering out recipe %s", r.Name)
			return true, nil
		}
	}

	return false, nil
}

func (rf *RecipeFilterer) FilterMultiple(ctx context.Context, r []types.OpenInstallationRecipe, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	results := []types.OpenInstallationRecipe{}
	for _, recipe := range r {
		filtered, err := rf.Filter(ctx, &recipe, m)
		if err != nil {
			return nil, err
		}

		if !filtered {
			results = append(results, recipe)
		}
	}

	return results, nil
}

type RecipeFilter interface {
	Execute(context.Context, *types.OpenInstallationRecipe, *types.DiscoveryManifest) (bool, error)
}

type ProcessMatchRecipeFilter struct {
	processMatchFinder ProcessMatchFinder
}

func NewProcessMatchRecipeFilter() *ProcessMatchRecipeFilter {
	return &ProcessMatchRecipeFilter{
		processMatchFinder: NewRegexProcessMatchFinder(),
	}
}

func (f *ProcessMatchRecipeFilter) Execute(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	matches, err := f.processMatchFinder.FindMatches(ctx, m.DiscoveredProcesses, *r)
	if err != nil {
		fmt.Println(3)
		return false, err
	}

	return len(r.ProcessMatch) > 0 && len(matches) == 0, nil
}

type ScriptEvaluationRecipeFilter struct {
	recipeExecutor execution.RecipeExecutor
}

func NewScriptEvaluationRecipeFilter() *ScriptEvaluationRecipeFilter {
	return &ScriptEvaluationRecipeFilter{
		recipeExecutor: execution.NewShRecipeExecutor(),
	}
}

func (f *ScriptEvaluationRecipeFilter) Execute(ctx context.Context, r *types.OpenInstallationRecipe, m *types.DiscoveryManifest) (bool, error) {
	if err := f.recipeExecutor.ExecuteDiscovery(ctx, *r, types.RecipeVars{}); err != nil {
		return true, nil
	}

	return false, nil
}
