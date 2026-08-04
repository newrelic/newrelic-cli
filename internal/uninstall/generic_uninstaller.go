package uninstall

import (
	"errors"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// ErrPermissionDenied is returned (or wrapped) by a PackageManager when an
// operation fails because the process lacks the privileges to perform it.
// GenericUninstaller treats this as fatal for the whole recipe, unlike other
// per-item failures which are logged as warnings and don't stop the rest.
var ErrPermissionDenied = errors.New("permission denied")

// PackageManager abstracts the OS-specific mechanics of removing a package,
// stopping/disabling a service, and deleting a filesystem path.
type PackageManager interface {
	StopService(name string) error
	DisableService(name string) error
	RemovePackage(name string) error
	RemovePath(path string) error
}

// GenericUninstaller performs Tier 2 removal, driven purely by a recipe's
// declarative OpenInstallationUninstallMeta, with no recipe-authored script.
type GenericUninstaller struct {
	pm PackageManager
}

// NewGenericUninstaller returns a GenericUninstaller backed by the given PackageManager.
func NewGenericUninstaller(pm PackageManager) *GenericUninstaller {
	return &GenericUninstaller{pm: pm}
}

// Uninstall stops/disables the declared services, removes the declared packages,
// then removes the declared paths. A benign per-item error is collected as a
// warning and processing continues; an error wrapping ErrPermissionDenied aborts
// immediately and is returned as the fatal error.
func (g *GenericUninstaller) Uninstall(meta types.OpenInstallationUninstallMeta) (warnings []error, fatal error) {
	for _, name := range meta.Services {
		if err := g.pm.StopService(name); err != nil {
			if errors.Is(err, ErrPermissionDenied) {
				return warnings, err
			}
			warnings = append(warnings, err)
		}
		if err := g.pm.DisableService(name); err != nil {
			if errors.Is(err, ErrPermissionDenied) {
				return warnings, err
			}
			warnings = append(warnings, err)
		}
	}

	for _, name := range meta.Packages {
		if err := g.pm.RemovePackage(name); err != nil {
			if errors.Is(err, ErrPermissionDenied) {
				return warnings, err
			}
			warnings = append(warnings, err)
		}
	}

	for _, path := range meta.Paths {
		if err := g.pm.RemovePath(path); err != nil {
			if errors.Is(err, ErrPermissionDenied) {
				return warnings, err
			}
			warnings = append(warnings, err)
		}
	}

	return warnings, nil
}
