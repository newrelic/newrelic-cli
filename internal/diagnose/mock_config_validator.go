package diagnose

import "context"

type MockConfigValidator struct {
	result error
}

func NewMockConfigValidator(result error) *MockConfigValidator {
	return &MockConfigValidator{
		result: result,
	}
}

func (v *MockConfigValidator) Validate(ctx context.Context) error {
	return v.result
}
