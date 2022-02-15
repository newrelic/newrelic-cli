package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type MockLogMatchFinder struct {
	Matches []types.OpenInstallationLogMatch
}

// NewMockFileFilterer creates a new instance of MockFileFilterer.
func NewMockLogMatchFinder() LogMatchFinderDefinition {
	return &MockLogMatchFinder{}
}

func (m *MockLogMatchFinder) GetPaths(ctx context.Context, recipes []types.OpenInstallationRecipe) []types.OpenInstallationLogMatch {
	return m.Matches
}
