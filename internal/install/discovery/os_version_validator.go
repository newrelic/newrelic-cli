package discovery

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type OsVersionValidator struct {
	minMajor int
	minMinor int
	os       string
	platform string
}

var (
	versionNoLongerSupported = "This version of %s is no longer supported"
	noVersionMessage         = "Failed to identified a valid version of %s"
)

func NewOsVersionValidator(os string, platform string, minMajor int, minMinor int) *OsVersionValidator {
	validator := OsVersionValidator{
		minMajor: minMajor,
		minMinor: minMinor,
		os:       os,
		platform: platform,
	}

	return &validator
}

func (v *OsVersionValidator) Validate(m *types.DiscoveryManifest) error {
	if v.os != m.OS {
		return nil
	}
	if (m.Platform != v.platform) && v.platform != "" {
		return nil
	}

	versions := strings.Split(m.PlatformVersion, ".")

	switch len(versions) {
	case 0:
		return v.formatErrorMessage(noVersionMessage, m)
	case 1:
		major, err := strconv.Atoi(versions[0])
		if err == nil {
			return v.ensureMinimumVersion(major, v.minMinor-1, m)
		}
		return v.formatErrorMessage(noVersionMessage, m)
	default:
		major, aerr := strconv.Atoi(versions[0])
		if aerr == nil {
			minor, ierr := strconv.Atoi(versions[1])
			if ierr == nil {
				return v.ensureMinimumVersion(major, minor, m)
			}
		}
	}

	return v.formatErrorMessage(noVersionMessage, m)
}

func (v *OsVersionValidator) ensureMinimumVersion(major int, minor int, m *types.DiscoveryManifest) error {
	if major < v.minMajor {
		return v.formatErrorMessage(versionNoLongerSupported, m)
	}
	if major == v.minMajor {
		if minor < v.minMinor {
			return v.formatErrorMessage(versionNoLongerSupported, m)
		}
	}
	return nil
}

func (v *OsVersionValidator) formatErrorMessage(message string, m *types.DiscoveryManifest) error {
	targetSpecific := m.OS
	if m.Platform != "" {
		targetSpecific = fmt.Sprintf("%s/%s", targetSpecific, m.Platform)
	}
	return fmt.Errorf(message, targetSpecific)
}
