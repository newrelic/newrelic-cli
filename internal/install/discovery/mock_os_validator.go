package discovery

import "github.com/newrelic/newrelic-cli/internal/install/types"

type MockOsValidator struct {
	Error error
}

func NewMockOsValidator() *MockOsValidator {
	validator := MockOsValidator{}
	return &validator
}

func (v *MockOsValidator) Validate(m *types.DiscoveryManifest) error {
	return v.Error
}
