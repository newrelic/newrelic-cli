package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockProcessFilterer struct {
	FilterVal       []types.MatchedProcess
	FilterCallCount int
	FilterErr       error
}

func NewMockProcessFilterer() *MockProcessFilterer {
	return &MockProcessFilterer{}
}

func (f *MockProcessFilterer) Filter(context.Context, []types.GenericProcess, []types.OpenInstallationRecipe) ([]types.MatchedProcess, error) {
	f.FilterCallCount++

	if f.FilterErr != nil {
		return nil, f.FilterErr
	}

	return f.FilterVal, nil
}
