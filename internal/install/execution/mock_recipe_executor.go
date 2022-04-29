package execution

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeExecutor struct {
	ExecuteErr error
}

func NewMockRecipeExecutor() *MockRecipeExecutor {
	return &MockRecipeExecutor{}
}

func (m *MockRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return m.ExecuteErr
}

func (m *MockRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return m.ExecuteErr
}

func (m *MockRecipeExecutor) GetOutput() *OutputParser {
	return NewOutputParser(map[string]interface{}{})
}
