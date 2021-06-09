package diagnose

import "context"

type MockConfigValidator struct{}

func NewMockConfigValidator() *MockConfigValidator {
	return &MockConfigValidator{}
}

func (v *MockConfigValidator) Validate(ctx context.Context) error {
	return nil
}
