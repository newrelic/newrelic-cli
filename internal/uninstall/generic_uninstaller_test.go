package uninstall

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type fakePackageManager struct {
	stoppedServices   []string
	disabledServices  []string
	removedPackages   []string
	removedPaths      []string
	stopServiceErr    map[string]error
	disableServiceErr map[string]error
	removePackageErr  map[string]error
	removePathErr     map[string]error
}

func (f *fakePackageManager) StopService(name string) error {
	f.stoppedServices = append(f.stoppedServices, name)
	return f.stopServiceErr[name]
}

func (f *fakePackageManager) DisableService(name string) error {
	f.disabledServices = append(f.disabledServices, name)
	return f.disableServiceErr[name]
}

func (f *fakePackageManager) RemovePackage(name string) error {
	f.removedPackages = append(f.removedPackages, name)
	return f.removePackageErr[name]
}

func (f *fakePackageManager) RemovePath(path string) error {
	f.removedPaths = append(f.removedPaths, path)
	return f.removePathErr[path]
}

func TestGenericUninstaller_RemovesEverythingWhenNoErrors(t *testing.T) {
	pm := &fakePackageManager{}
	g := NewGenericUninstaller(pm)

	warnings, fatal := g.Uninstall(types.OpenInstallationUninstallMeta{
		Services: []string{"newrelic-infra"},
		Packages: []string{"newrelic-infra"},
		Paths:    []string{"/etc/newrelic-infra"},
	})

	require.NoError(t, fatal)
	require.Empty(t, warnings)
	require.Equal(t, []string{"newrelic-infra"}, pm.stoppedServices)
	require.Equal(t, []string{"newrelic-infra"}, pm.disabledServices)
	require.Equal(t, []string{"newrelic-infra"}, pm.removedPackages)
	require.Equal(t, []string{"/etc/newrelic-infra"}, pm.removedPaths)
}

func TestGenericUninstaller_BenignItemFailureContinuesToTheRest(t *testing.T) {
	pm := &fakePackageManager{
		removePackageErr: map[string]error{"already-gone": errors.New("package not installed")},
	}
	g := NewGenericUninstaller(pm)

	warnings, fatal := g.Uninstall(types.OpenInstallationUninstallMeta{
		Packages: []string{"already-gone"},
		Paths:    []string{"/etc/newrelic-infra"},
	})

	require.NoError(t, fatal)
	require.Len(t, warnings, 1)
	require.Equal(t, []string{"/etc/newrelic-infra"}, pm.removedPaths)
}

func TestGenericUninstaller_PermissionErrorAbortsRemainingItems(t *testing.T) {
	pm := &fakePackageManager{
		removePackageErr: map[string]error{"newrelic-infra": ErrPermissionDenied},
	}
	g := NewGenericUninstaller(pm)

	warnings, fatal := g.Uninstall(types.OpenInstallationUninstallMeta{
		Packages: []string{"newrelic-infra"},
		Paths:    []string{"/etc/newrelic-infra"},
	})

	require.Error(t, fatal)
	require.True(t, errors.Is(fatal, ErrPermissionDenied))
	require.Empty(t, warnings)
	require.Empty(t, pm.removedPaths, "paths should not be touched after a fatal permission error")
}
