package validation

import (
	"context"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type MockRecipeValidator struct {
	ValidateErrs      []error
	ValidateErr       error
	ValidateCallCount int
	ValidateVal       string
	ValidateVals      []string
}

func NewMockRecipeValidator() *MockRecipeValidator {
	return &MockRecipeValidator{}
}

func (m *MockRecipeValidator) ValidateQuery(ctx context.Context, query string) (string, error) {
	return m.validate(ctx)
}

func (m *MockRecipeValidator) ValidateRecipe(ctx context.Context, dm types.DiscoveryManifest, r types.OpenInstallationRecipe) (string, error) {
	return m.validate(ctx)
}

func (m *MockRecipeValidator) validate(ctx context.Context) (string, error) {
	m.ValidateCallCount++

	var err error
	var val string

	if len(m.ValidateErrs) > 0 {
		i := utils.MinOf(m.ValidateCallCount, len(m.ValidateErrs)) - 1
		err = m.ValidateErrs[i]
	} else {
		err = m.ValidateErr
	}

	if len(m.ValidateVals) > 0 {
		i := utils.MinOf(m.ValidateCallCount, len(m.ValidateVals)) - 1
		val = m.ValidateVals[i]
	} else {
		val = m.ValidateVal
	}

	time.Sleep(1 * time.Millisecond)

	return val, err
}
