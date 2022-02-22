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
				Recipe: recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []execution.RecipeStatusType{
					execution.RecipeStatusTypes.AVAILABLE,
				},
			},
			{
				Recipe: recipes.NewRecipeBuilder().ID("ID2").Name("recipe2").Build(),
				DetectedStatuses: []execution.RecipeStatusType{
					execution.RecipeStatusTypes.AVAILABLE,
				},
			},
		},
	}

	actualError := bundleInstallerTestImpl.bundleInstaller.InstallStopOnError(&bundle, true)

	//Should stop on first recipe
	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
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
				Recipe: recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []execution.RecipeStatusType{
					execution.RecipeStatusTypes.AVAILABLE,
				},
			},
			{
				Recipe: recipes.NewRecipeBuilder().ID("ID2").Name("recipe2").Build(),
				DetectedStatuses: []execution.RecipeStatusType{
					execution.RecipeStatusTypes.AVAILABLE,
				},
			},
		},
	}

	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, true)

	//Should try both recipes
	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 2)
}

func TestReportsStatusHasSingleStatusWhenStatusNotAvailable(t *testing.T) {
	setup()
	expectedStatus := execution.RecipeStatusTypes.RECOMMENDED
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddStatus(expectedStatus)

	bundleInstallerTestImpl.bundleInstaller.reportStatus(bundle)

	assert.Equal(t, expectedStatus, bundle.BundleRecipes[0].DetectedStatuses[0])
	assert.Equal(t, 1, len(bundle.BundleRecipes[0].DetectedStatuses))
}

func TestReportsStatusHasDetectedAndAvailableWhenStatusIsAvailable(t *testing.T) {
	setup()
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddStatus(execution.RecipeStatusTypes.AVAILABLE)

	bundleInstallerTestImpl.bundleInstaller.reportStatus(bundle)

	assert.True(t, bundle.BundleRecipes[0].HasStatus(execution.RecipeStatusTypes.AVAILABLE))
	assert.True(t, bundle.BundleRecipes[0].HasStatus(execution.RecipeStatusTypes.DETECTED))
	assert.Equal(t, 2, len(bundle.BundleRecipes[0].DetectedStatuses))
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
