package discovery

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ManifestValidator struct {
	validators []validator
}

type validator interface {
	Validate(m *types.DiscoveryManifest) error
}

var (
	errorPrefix       = "Installation requirements error:"
	errorPrefixFormat = errorPrefix + " %s"
)

// NewManifestValidator returns a new instance of ManifestValidator.
func NewManifestValidator() *ManifestValidator {
	mv := ManifestValidator{
		validators: []validator{},
	}

	mv.validators = append(mv.validators, NewOsValidator())
	mv.validators = append(mv.validators, NewOsVersionValidator("windows", "", 6, 2))
	mv.validators = append(mv.validators, NewOsVersionValidator("linux", "ubuntu", 16, 04))

	return &mv
}

func (mv *ManifestValidator) Validate(m *types.DiscoveryManifest) error {
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
		return &types.UnsupportedOperatingSytemError{
			Err: fmt.Errorf(errorPrefixFormat, accumulator),
		}
	}
	return nil
}

func (mv *ManifestValidator) FindAllValidationErrors(m *types.DiscoveryManifest) []error {
	errors := []error{}

	for _, validator := range mv.validators {
		err := validator.Validate(m)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
