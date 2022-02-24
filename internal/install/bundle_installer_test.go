package install

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/ux"

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
	manifest := types.DiscoveryManifest{}

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

/*
	New tests needed:
	1. getInstallabeBundleRecipes
	2. InstallContinueOnError
		a. Test no installable recipes returns immediately
		b. Test prompt of no returns immediately (how can we mocka  prompt?)
		c  Test prompt of yes calls install (didn't seem to pass/fail as expected)
        d. make sure ux/prompt doesnt get called if bundle is Additional/Targeted
		e. Test if a mix of installable and not installable recipes, only installable trigger


	for os-specific unsupported, consider these cases

	3. When no recipe for CoreBundle
	4. No recipe for AdditionalBundle
*/

func TestInstallContinueOnErrorReturnsImmediately(t *testing.T) {
	setup()
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
	}

	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, true)

	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
}

func TestInstallContinueOnErrorReturnsImmediatelyWhenNoIsEntered(t *testing.T) {
	setup()
	mockPrompter := ux.NewMockPrompter()
	mockPrompter.PromptYesNoVal = false
	bundleInstallerTestImpl.bundleInstaller.prompter = mockPrompter

	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
	}

	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, false)

	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
	assert.Equal(t, 1, bundleInstallerTestImpl.statusReporter.counter)
}

func TestInstallContinueOnErrorIgnoresUxPromptIfBundleIsAdditionalTargeted(t *testing.T) {
	setup()
	mockPrompter := ux.NewMockPrompter()
	bundleInstallerTestImpl.bundleInstaller.prompter = mockPrompter
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
		Type: recipes.BundleTypes.ADDITIONALTARGETED,
	}

	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, true)

	assert.Equal(t, 0, mockPrompter.PromptMultiSelectCallCount)
}

// TODO come back to this, not sure if the test makes sense
//func TestInstallContinueOnErrorReturnsInstallsWhenYesIsEntered(t *testing.T) {
//	setup()
//	mockPrompter := ux.NewMockPrompter()
//	mockPrompter.PromptYesNoVal = true
//	bundleInstallerTestImpl.bundleInstaller.prompter = mockPrompter
//
//	mockedRecipeInstaller := new(mockRecipeInstaller)
//	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))
//	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller
//
//	bundle := recipes.Bundle{
//		BundleRecipes: []*recipes.BundleRecipe{
//			{
//				Recipe: recipes.NewRecipeBuilder().Name("recipe1").Build(),
//				DetectedStatuses: []execution.RecipeStatusType{
//					execution.RecipeStatusTypes.UNSUPPORTED,
//				},
//			},
//		},
//	}
//
//	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, false)
//
//	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
//}

func TestInstallStopsOnErrorActuallyErrors(t *testing.T) {
	setup()
	expectedError := errors.New("Kaboom " + time.Now().String())
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", expectedError)
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller
	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.AVAILABLE}},
			},
			{
				Recipe:           recipes.NewRecipeBuilder().ID("ID2").Name("recipe2").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.AVAILABLE}},
			},
		},
	}

	actualError := bundleInstallerTestImpl.bundleInstaller.InstallStopOnError(&bundle, true)

	//Should stop on first recipe
	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func TestInstallContinueOnErrorOnlyInstallsAvailableRecipesInBundle(t *testing.T) {
	setup()
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("great success", nil)
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.AVAILABLE}},
			},
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe2").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe3").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.DETECTED}},
			},
		},
	}

	bundleInstallerTestImpl.bundleInstaller.InstallContinueOnError(&bundle, true)

	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.True(t, bundleInstallerTestImpl.bundleInstaller.installedRecipes["recipe1"])
	assert.False(t, bundleInstallerTestImpl.bundleInstaller.installedRecipes["recipe2"])
	assert.False(t, bundleInstallerTestImpl.bundleInstaller.installedRecipes["recipe3"])
}

func TestInstallContinueOnErrorKeepsInstalling(t *testing.T) {
	setup()
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("great success", nil)
	bundleInstallerTestImpl.bundleInstaller.recipeInstaller = mockedRecipeInstaller
	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.AVAILABLE}},
			},
			{
				Recipe:           recipes.NewRecipeBuilder().ID("ID2").Name("recipe2").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.AVAILABLE}},
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
	bundle.BundleRecipes[0].AddDetectionStatus(expectedStatus, 0)

	bundleInstallerTestImpl.bundleInstaller.reportBundleStatus(bundle)

	assert.Equal(t, expectedStatus, bundle.BundleRecipes[0].DetectedStatuses[0].Status)
	assert.Equal(t, 1, len(bundle.BundleRecipes[0].DetectedStatuses))
}

func TestReportsStatusHasDetectedAndAvailableWhenStatusIsAvailable(t *testing.T) {
	setup()
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, 0)

	bundleInstallerTestImpl.bundleInstaller.reportBundleStatus(bundle)

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
