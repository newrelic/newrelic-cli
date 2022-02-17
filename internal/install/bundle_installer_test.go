package install

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type bundleInstallerTest struct {
	// FIXME currently unused...
	// installedRecipes map[string]bool
	ctx             context.Context
	manifest        *types.DiscoveryManifest
	recipeInstaller *mockRecipeInstaller
	statusReporter  *mockStatusReporter
	bundleInstaller *BundleInstaller
}

var (
	bundleInstallerTestImpl *bundleInstallerTest
)

func setup() {
	manifest := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{},
	}

	bundleInstallerTestImpl = &bundleInstallerTest{
		statusReporter:  &mockStatusReporter{},
		recipeInstaller: &mockRecipeInstaller{},
		manifest:        &manifest,
		ctx:             context.Background(),
	}

	bundleInstallerTestImpl.bundleInstaller = NewBundleInstaller(
		bundleInstallerTestImpl.ctx,
		bundleInstallerTestImpl.manifest,
		bundleInstallerTestImpl.recipeInstaller,
		bundleInstallerTestImpl.statusReporter)
}

func TestInstallStopsOnErrorActuallyErrors(t *testing.T) {
	setup()
	expectedError := errors.New("Kaboom " + time.Now().String())
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", expectedError)
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe: recipes.NewRecipeBuilder().ID("0").Name("recipe1").Build(),
				RecipeStatuses: []recipes.RecipeStatus{
					{
						Status:     execution.RecipeStatusTypes.AVAILABLE,
						StatusTime: time.Now(),
					},
				},
			},
		},
	}

	actualError := bundleInstallerTestImpl.bundleInstaller.InstallStopOnError(&bundle, true)

	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func TestInstallContinueOnErrorKeepsInstalling(t *testing.T) {
	setup()
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("great success", nil)
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe: recipes.NewRecipeBuilder().ID("0").Name("recipe1").Build(),
				RecipeStatuses: []recipes.RecipeStatus{
					{
						Status:     execution.RecipeStatusTypes.AVAILABLE,
						StatusTime: time.Now(),
					},
				},
			},
		},
	}

	//TODO: Need to find out how to verify error was thrown
	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, true)
}

func TestReportsStatusHasSingleStatusWhenStatusNotAvailable(t *testing.T) {
	setup()
	expectedStatus := execution.RecipeStatusTypes.RECOMMENDED
	expectedStatusTime := time.Now()
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddStatus(expectedStatus, expectedStatusTime)

	bundleInstallerTestImpl.bundleInstaller.reportStatus(bundle)

	assert.Equal(t, expectedStatus, bundle.BundleRecipes[0].RecipeStatuses[0].Status)
	assert.Equal(t, expectedStatusTime, bundle.BundleRecipes[0].RecipeStatuses[0].StatusTime)
	assert.Equal(t, 1, len(bundle.BundleRecipes[0].RecipeStatuses))
}

func TestReportsStatusHasDetectedAndAvailableWhenStatusIsAvailable(t *testing.T) {
	setup()
	statusTime := time.Now()
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddStatus(execution.RecipeStatusTypes.AVAILABLE, statusTime)

	bundleInstallerTestImpl.bundleInstaller.reportStatus(bundle)

	assert.True(t, bundle.BundleRecipes[0].HasStatus(execution.RecipeStatusTypes.AVAILABLE))
	assert.True(t, bundle.BundleRecipes[0].HasStatus(execution.RecipeStatusTypes.DETECTED))
	assert.Equal(t, 2, len(bundle.BundleRecipes[0].RecipeStatuses))
	assert.Equal(t, statusTime, bundle.BundleRecipes[0].RecipeStatuses[0].StatusTime)
	assert.Equal(t, statusTime, bundle.BundleRecipes[0].RecipeStatuses[1].StatusTime)

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
