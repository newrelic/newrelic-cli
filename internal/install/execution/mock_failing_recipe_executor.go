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

func (m *MockFailingRecipeExecutor) Prepare(ctx context.Context, dm types.DiscoveryManifest, r types.OpenInstallationRecipe, y bool, z string) (types.RecipeVars, error) {
	return types.RecipeVars{}, nil
}

func (m *MockFailingRecipeExecutor) Execute(ctx context.Context, dm types.DiscoveryManifest, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return fmt.Errorf("something went wrong")
}
