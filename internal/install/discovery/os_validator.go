package discovery

import (
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type OsValidator struct{}

var (
	NoOperatingSystemDetected         = "failed to identify a valid operating system"
	OperatingSystemNotSupportedPrefix = "operating system"
	OperatingSystemNotSupportedSuffix = "is not supported"
	OperatingSystemNotSupportedFormat = OperatingSystemNotSupportedPrefix + " %s " + OperatingSystemNotSupportedSuffix
)

func NewOsValidator() *OsValidator {
	validator := OsValidator{}

	return &validator
}

func (v *OsValidator) Execute(m *types.DiscoveryManifest) string {
	if m.OS == "" {
		return NoOperatingSystemDetected
	}
	if !(strings.ToLower(m.OS) == "linux" || strings.ToLower(m.OS) == "windows") {
		return fmt.Sprintf(OperatingSystemNotSupportedFormat, m.OS)
	}
	return ""
}
