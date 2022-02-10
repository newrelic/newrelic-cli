package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockProcessMatchFinder struct {
	matchedProcesses []types.MatchedProcess
}

func NewMockProcessMatchFinder() *MockProcessMatchFinder {
	return &MockProcessMatchFinder{
		matchedProcesses: []types.MatchedProcess{},
	}
}

func (f *MockProcessMatchFinder) FindMatchesMultiple(ctx context.Context, processes []types.GenericProcess, recipes []types.OpenInstallationRecipe) []types.MatchedProcess {
	return f.matchedProcesses
}

func (f *MockProcessMatchFinder) FindMatches(ctx context.Context, processes []types.GenericProcess, recipe types.OpenInstallationRecipe) []types.MatchedProcess {
	return f.matchedProcesses
}
