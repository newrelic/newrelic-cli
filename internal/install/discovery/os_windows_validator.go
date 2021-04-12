package discovery

import (
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type OsWindowsValidator struct{}

var (
	WindowsVersionNoLongerSupported = "This version of Windows is no longer supported"
	WindowsNoVersionMessage         = "Failed to identified a valid version of Windows"
)

func NewOsWindowsValidator() *OsWindowsValidator {
	validator := OsWindowsValidator{}

	return &validator
}

func (v *OsWindowsValidator) Execute(m *types.DiscoveryManifest) string {
	if m.OS != "windows" {
		return ""
	}

	versions := strings.Split(m.PlatformVersion, ".")

	switch len(versions) {
	case 0:
		return WindowsNoVersionMessage
	case 1:
		major, err := strconv.Atoi(versions[0])
		if err == nil {
			return ensureMinimumVersion(major, 0)
		}
		return WindowsNoVersionMessage
	default:
		major, aerr := strconv.Atoi(versions[0])
		if aerr == nil {
			minor, ierr := strconv.Atoi(versions[1])
			if ierr == nil {
				return ensureMinimumVersion(major, minor)
			}
		}
	}

	return WindowsNoVersionMessage
}

func ensureMinimumVersion(major int, minor int) string {
	if major < 6 {
		return WindowsVersionNoLongerSupported
	}
	if major == 6 {
		if minor == 0 {
			return WindowsVersionNoLongerSupported
		}
	}
	return ""
}
