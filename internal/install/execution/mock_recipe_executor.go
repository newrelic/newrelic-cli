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

func (m *MockRecipeExecutor) Execute(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe) error {
	return nil
}
