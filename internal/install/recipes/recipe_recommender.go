package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type RecipeRecommender struct {
	recipeFetcher   RecipeFetcher
	processFilterer ProcessFilterer
	recipeExecutor  execution.RecipeExecutor
	allRecipes      []types.OpenInstallationRecipe
}

func NewRecipeRecommender(rf RecipeFetcher, pf ProcessFilterer, re execution.RecipeExecutor) *RecipeRecommender {
	return &RecipeRecommender{
		recipeFetcher:   rf,
		processFilterer: pf,
		recipeExecutor:  re,
	}
}

func (r *RecipeRecommender) Recommend(ctx context.Context, m *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	if r.allRecipes == nil {
		if err := r.fetchAllRecipes(ctx, m); err != nil {
			return nil, err
		}
	}

	recCtx := newRecommendationContext(m)

	if err := r.findProcessMatches(ctx, recCtx); err != nil {
		return nil, err
	}

	if err := r.findCustomScriptedMatches(ctx, recCtx); err != nil {
		return nil, err
	}

	return recCtx.matches, nil
}

func (r *RecipeRecommender) fetchAllRecipes(ctx context.Context, m *types.DiscoveryManifest) error {
	allRecipes, err := r.recipeFetcher.FetchRecipes(ctx, m)
	if err != nil {
		return err
	}

	r.allRecipes = allRecipes
	return nil
}

func (r *RecipeRecommender) findProcessMatches(ctx context.Context, recCtx *recommendationContext) error {
	matches, err := r.processFilterer.Filter(ctx, recCtx.discoveryManifest.DiscoveredProcesses, r.allRecipes)
	if err != nil {
		return err
	}

	for _, match := range matches {
		recCtx.addMatch(&match.MatchingRecipe)
	}

	return nil
}

func (r *RecipeRecommender) findCustomScriptedMatches(ctx context.Context, recCtx *recommendationContext) error {
	for _, recipe := range r.allRecipes {
		if recipe.PreInstall.RequireAtDiscovery != "" {
			if err := r.recipeExecutor.ExecuteDiscovery(ctx, recipe, types.RecipeVars{}); err != nil {
				if err != nil {
					continue
				}
			}

			recCtx.addMatch(&recipe)
		}
	}

	return nil
}

type recommendationContext struct {
	discoveryManifest *types.DiscoveryManifest
	matches           []types.OpenInstallationRecipe
	matchMap          map[*types.OpenInstallationRecipe]bool
}

func newRecommendationContext(m *types.DiscoveryManifest) *recommendationContext {
	return &recommendationContext{
		discoveryManifest: m,
		matches:           []types.OpenInstallationRecipe{},
		matchMap:          map[*types.OpenInstallationRecipe]bool{},
	}
}

func (c *recommendationContext) addMatch(r *types.OpenInstallationRecipe) {
	if !c.matchMap[r] {
		c.matches = append(c.matches, *r)
		c.matchMap[r] = true
	}
}
