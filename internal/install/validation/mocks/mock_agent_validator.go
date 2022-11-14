package mocks

import "context"

type MockAgentValidator struct {
	Error error
}

func (m *MockAgentValidator) Validate(ctx context.Context, url string) (string, error) {
	return "", m.Error
}
