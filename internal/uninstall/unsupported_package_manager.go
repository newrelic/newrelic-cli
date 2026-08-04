package uninstall

import "fmt"

// UnsupportedPackageManager is the default PackageManager until a real,
// OS-specific implementation (systemd/apt/yum/msiexec) is built. It fails
// every operation with a clear message instead of silently doing nothing.
type UnsupportedPackageManager struct{}

func (u *UnsupportedPackageManager) StopService(name string) error {
	return fmt.Errorf("stopping service %q is not supported on this platform yet", name)
}

func (u *UnsupportedPackageManager) DisableService(name string) error {
	return fmt.Errorf("disabling service %q is not supported on this platform yet", name)
}

func (u *UnsupportedPackageManager) RemovePackage(name string) error {
	return fmt.Errorf("removing package %q is not supported on this platform yet", name)
}

func (u *UnsupportedPackageManager) RemovePath(path string) error {
	return fmt.Errorf("removing path %q is not supported on this platform yet", path)
}
