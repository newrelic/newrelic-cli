package execution

import (
	"context"
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockFailingRecipeExecutor struct {
	result bool
}

func NewMockFailingRecipeExecutor() *MockFailingRecipeExecutor {
	return &MockFailingRecipeExecutor{
		result: false,
	}
}

func (m *MockFailingRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return fmt.Errorf("something went wrong")
}

func (m *MockFailingRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return fmt.Errorf("something went wrong")
}
