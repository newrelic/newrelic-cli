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

func (m *MockRecipeExecutor) Prepare(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe, y bool) (types.RecipeVars, error) {
	return types.RecipeVars{}, nil
}

func (m *MockRecipeExecutor) Execute(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe, v types.RecipeVars) error {
	return nil
}
