package ux

import "github.com/newrelic/newrelic-cli/internal/install/types"

type MockProgressIndicator struct{}

func NewMockProgressIndicator() *MockProgressIndicator {
	return &MockProgressIndicator{}
}

func (s *MockProgressIndicator) Fail(types.Recipe) {
}

func (s MockProgressIndicator) Success(types.Recipe) {
}

func (s *MockProgressIndicator) Start(types.Recipe) {
}

func (s MockProgressIndicator) Stop() {
}
