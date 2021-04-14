package discovery

import (
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type OsValidator struct{}

var (
	noOperatingSystemDetected         = "failed to identify a valid operating system"
	operatingSystemNotSupportedPrefix = "operating system"
	operatingSystemNotSupportedSuffix = "is not supported"
	operatingSystemNotSupportedFormat = operatingSystemNotSupportedPrefix + " %s " + operatingSystemNotSupportedSuffix
)

func NewOsValidator() *OsValidator {
	validator := OsValidator{}

	return &validator
}

func (v *OsValidator) Execute(m *types.DiscoveryManifest) error {
	if m.OS == "" {
		return fmt.Errorf(noOperatingSystemDetected)
	}
	if !(strings.ToLower(m.OS) == "linux" || strings.ToLower(m.OS) == "windows") {
		return fmt.Errorf(operatingSystemNotSupportedFormat, m.OS)
	}
	return nil
}
