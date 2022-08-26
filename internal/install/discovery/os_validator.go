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

func (v *OsValidator) Validate(m *types.DiscoveryManifest) error {
	fmt.Print("\n****************************\n")
	fmt.Printf("\n THING:  %+v \n", m.OS)
	fmt.Print("\n****************************\n")
	if m.OS == "" {
		return fmt.Errorf(noOperatingSystemDetected)
	}
	if !(strings.ToLower(m.OS) == "linux" ||
		strings.ToLower(m.OS) == "windows" ||
		strings.ToLower(m.OS) == "darwin") {
		return fmt.Errorf(operatingSystemNotSupportedFormat, m.OS)
	}
	return nil
}
