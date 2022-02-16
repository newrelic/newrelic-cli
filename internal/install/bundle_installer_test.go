package install

import (
	"context"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

type bundleInstallerTest struct {
	installedRecipes map[string]bool
	ctx              context.Context
	manifest         *types.DiscoveryManifest
	recipeInstaller  *RecipeInstaller
}

var (
	bundleInstallerTestImpl = bundleInstallerTest{
		installedRecipes: nil,
	}
)

type mockInstallBundleRecipe struct {
	mock.Mock
}

func createBundleInstaller() *BundleInstaller {
	return NewBundleInstaller(bundleInstallerTestImpl.ctx, bundleInstallerTestImpl.manifest, bundleInstallerTestImpl.recipeInstaller)
}

// public functions
func TestBundleInstallerStopsOnError(t *testing.T) {

	expectedError := "I am an error"
	mockInstallBundleRecipe := new(mockInstallBundleRecipe)
	mockInstallBundleRecipe.On("installBundleRecipe", mock.Anything, mock.AnythingOfType("bool")).Return(expectedError)
	bundleInstaller := createBundleInstaller()

	//FIXME pass in recipes.Bundle
	actualError := bundleInstaller.InstallStopOnError(nil, true)

	require.Equal(t, expectedError, actualError)
}

func TestBundleInstallerContinuesOnError(t *testing.T) {
	require.Fail(t, "Implement me")
}

// private function, no logic; delegates to internal/install/execution/install_status#ReportStatus
//func TestBundleInstallerReportsStatus(t *testing.T) {
//	setup()
//
//	require.Fail(t, "Implement me")
//}

func TestBundleInstallerInstallsBundleRecipes(t *testing.T) {
	require.Fail(t, "Implement me")
}

func TestBundleInstallerInstallsBundleRecipesWithDependencies(t *testing.T) {
	require.Fail(t, "Implement me")
}
