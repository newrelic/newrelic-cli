package execution

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeExecutor struct {
	result bool
}

func NewMockRecipeExecutor() *MockRecipeExecutor {
	return &MockRecipeExecutor{
		result: false,
	}
}

func (m *MockRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return nil
}

func (m *MockRecipeExecutor) ExecuteDiscovery(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return nil
}
