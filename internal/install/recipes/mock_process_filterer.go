package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockProcessFilterer struct {
	FilterVal       []types.MatchedProcess
	FilterCallCount int
}

func NewMockProcessFilterer() *MockProcessFilterer {
	return &MockProcessFilterer{}
}

func (f *MockProcessFilterer) Filter(context.Context, []types.GenericProcess, []types.OpenInstallationRecipe) []types.MatchedProcess {
	f.FilterCallCount++

	return f.FilterVal
}
