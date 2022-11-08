package discovery

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func NewEmptyManifestValidator() *ManifestValidator {
	return &ManifestValidator{}
}

type MockManifestValidator struct {
	validators []Validator
}

func NewMockManifestValidator(mockValidator Validator) *ManifestValidator {
	mv := ManifestValidator{
		validators: []Validator{mockValidator},
	}

	return &mv
}

func (mv *MockManifestValidator) Validate(m *types.DiscoveryManifest) error {
	var accumulator error
	errors := mv.FindAllValidationErrors(m)
	for _, single := range errors {
		if accumulator == nil {
			accumulator = single
		} else {
			accumulator = fmt.Errorf("%s, %s", accumulator, single.Error())
		}
	}

	if accumulator != nil {
		// Flag as unsupported OS
		m.IsUnsupported = true
		return &types.UnsupportedOperatingSystemError{
			Err: fmt.Errorf(errorPrefixFormat, accumulator),
		}
	}
	return nil
}

func (mv *MockManifestValidator) FindAllValidationErrors(m *types.DiscoveryManifest) []error {
	errors := []error{}

	for _, validator := range mv.validators {
		err := validator.Validate(m)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
