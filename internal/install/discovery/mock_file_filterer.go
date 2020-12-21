package discovery

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// MockFileFilterer is a mock implementation of the FileFilterer interface that
// provides method spies for testing scenarios.
type MockFileFilterer struct {
	FilterCallCount int
	FilterErr       error
	FilterVal       []types.LogMatch
}

// NewMockFileFilterer creates a new instance of MockFileFilterer.
func NewMockFileFilterer() *MockFileFilterer {
	return &MockFileFilterer{}
}

func (m *MockFileFilterer) Filter(ctx context.Context, recipes []types.Recipe) ([]types.LogMatch, error) {
	m.FilterCallCount++
	return m.FilterVal, m.FilterErr
}
