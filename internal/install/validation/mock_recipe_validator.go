package validation

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeValidator struct {
	ValidateErr       error
	ValidateCallCount int
	ValidateVal       string
}

func NewMockRecipeValidator() *MockRecipeValidator {
	return &MockRecipeValidator{}
}

func (m *MockRecipeValidator) Validate(ctx context.Context, dm types.DiscoveryManifest, r types.Recipe) (string, error) {
	m.ValidateCallCount++
	return m.ValidateVal, m.ValidateErr
}
