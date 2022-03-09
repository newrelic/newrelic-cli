package diagnose

import "context"

type MockConfigValidator struct {
	Error error
}

func NewMockConfigValidator() *MockConfigValidator {
	return &MockConfigValidator{}
}

func (v *MockConfigValidator) Validate(ctx context.Context) error {
	return v.Error
}
