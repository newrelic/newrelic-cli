package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeEvaluator struct {
	recipeStatus map[string]execution.RecipeStatusType
}

func NewMockRecipeEvaluator() *MockRecipeEvaluator {
	return &MockRecipeEvaluator{
		recipeStatus: make(map[string]execution.RecipeStatusType),
	}
}

func (mre *MockRecipeEvaluator) WithRecipeStatus(recipe *types.OpenInstallationRecipe, status execution.RecipeStatusType) {
	mre.recipeStatus[recipe.Name] = status
}

func (mre *MockRecipeEvaluator) DetectionStatus(ctx context.Context, recipe *types.OpenInstallationRecipe) execution.RecipeStatusType {
	if status, ok := mre.recipeStatus[recipe.Name]; ok {
		return status
	}
	return execution.RecipeStatusTypes.NULL
}
