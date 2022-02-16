package install

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type bundleInstallerTest struct {
	installedRecipes map[string]bool
	ctx              context.Context
	manifest         *types.DiscoveryManifest
	recipeInstaller  *RecipeInstaller
	statusReporter   *mockStatusReporter
	bundleInstaller  *BundleInstaller
}

var (
	bundleInstallerTestImpl *bundleInstallerTest
)

type mockInstallBundleRecipe struct {
	mock.Mock
}

func setup() {
	bundleInstallerTestImpl = &bundleInstallerTest{
		statusReporter: &mockStatusReporter{},
	}

	bundleInstallerTestImpl.bundleInstaller = NewBundleInstaller(
		bundleInstallerTestImpl.ctx,
		bundleInstallerTestImpl.manifest,
		bundleInstallerTestImpl.recipeInstaller,
		bundleInstallerTestImpl.statusReporter)
}

// public functions
// func TestBundleInstallerStopsOnError(t *testing.T) {
// 	setup()

// 	expectedError := "I am an error"
// 	mockInstallBundleRecipe := new(mockInstallBundleRecipe)
// 	mockInstallBundleRecipe.On("installBundleRecipe", mock.Anything, mock.AnythingOfType("bool")).Return(expectedError)

// 	//FIXME pass in recipes.Bundle
// 	actualError := bundleInstaller.InstallStopOnError(nil, true)

// 	require.Equal(t, expectedError, actualError)
// }

func TestBundleInstallerContinuesOnError(t *testing.T) {
	require.Fail(t, "Implement me")
}

func TestBundleInstallerReportsStatus(t *testing.T) {
	setup()
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddStatus(execution.RecipeStatusTypes.AVAILABLE)

	bundleInstallerTestImpl.bundleInstaller.reportStatus(bundle)

	actual := bundleInstallerTestImpl.statusReporter.counter
	expected := len(bundle.BundleRecipes[0].Statuses)
	require.Equal(t, expected, actual)

}

func givenBundle(recipeName string) *recipes.Bundle {
	bundle := &recipes.Bundle{}
	r := &types.OpenInstallationRecipe{
		Name: recipeName,
	}
	br := &recipes.BundleRecipe{
		Recipe: r,
	}
	bundle.AddRecipe(br)
	return bundle
}

func TestBundleInstallerInstallsBundleRecipes(t *testing.T) {
	require.Fail(t, "Implement me")
}

func TestBundleInstallerInstallsBundleRecipesWithDependencies(t *testing.T) {
	require.Fail(t, "Implement me")
}

type mockStatusReporter struct {
	counter int
}

func (sr *mockStatusReporter) ReportStatus(status execution.RecipeStatusType, recipe types.OpenInstallationRecipe) {
	sr.counter++
}
