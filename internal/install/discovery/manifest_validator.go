package discovery

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ManifestValidator struct {
	validators []Validator
}

type Validator interface {
	Execute(m *types.DiscoveryManifest) error
}

var (
	errorPrefix       = "Installation requirements error:"
	errorPrefixFormat = errorPrefix + " %s"
)

// NewManifestValidator returns a new instance of ManifestValidator.
func NewManifestValidator() *ManifestValidator {
	mv := ManifestValidator{
		validators: []Validator{},
	}

	mv.validators = append(mv.validators, NewOsValidator())
	mv.validators = append(mv.validators, NewOsWindowsValidator())

	return &mv
}

func (mv *ManifestValidator) Execute(m *types.DiscoveryManifest) error {
	var accumulator error = nil
	errors := mv.FindAllValidationErrors(m)
	for _, single := range errors {
		if accumulator == nil {
			accumulator = single
		} else {
			accumulator = fmt.Errorf("%s, %s", accumulator, single.Error())
		}
	}
	if accumulator != nil {
		return fmt.Errorf(errorPrefixFormat, accumulator)
	}
	return nil
}

func (mv *ManifestValidator) FindAllValidationErrors(m *types.DiscoveryManifest) []error {
	errors := []error{}

	for _, validator := range mv.validators {
		err := validator.Execute(m)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
