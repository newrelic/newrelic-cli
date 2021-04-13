package discovery

import (
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ManifestValidator struct {
	validators []Validator
}

type Validator interface {
	Execute(m *types.DiscoveryManifest) string
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

func (mv *ManifestValidator) Execute(m *types.DiscoveryManifest) string {
	accumulator := ""
	errors := mv.FindAllValidationErrors(m)
	for _, single := range errors {
		if accumulator == "" {
			accumulator = single
		} else {
			accumulator = fmt.Sprintf("%s, %s", accumulator, single)
		}
	}
	if accumulator != "" {
		return fmt.Sprintf(errorPrefixFormat, accumulator)
	}
	return ""
}

func (mv *ManifestValidator) FindAllValidationErrors(m *types.DiscoveryManifest) []string {
	errors := []string{}

	for _, validator := range mv.validators {
		result := validator.Execute(m)
		if result != "" {
			errors = append(errors, result)
		}
	}

	return errors
}
