package validation

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockRecipeValidator struct {
	Error error
}

func NewMockRecipeValidator() *MockRecipeValidator {
	return &MockRecipeValidator{}
}

func (m *MockRecipeValidator) ValidateRecipe(ctx context.Context, dm types.DiscoveryManifest, r types.OpenInstallationRecipe, vars types.RecipeVars) (string, error) {
	return "", m.Error
}
