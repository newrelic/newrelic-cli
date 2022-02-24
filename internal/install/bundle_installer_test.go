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

func newBundleInstaller(ctx context.Context, manifest *types.DiscoveryManifest, recipeInstallerInterface RecipeInstaller, statusReporter StatusReporter) *BundleInstaller {
	return &BundleInstaller{
		ctx:              ctx,
		manifest:         manifest,
		recipeInstaller:  recipeInstallerInterface,
		statusReporter:   statusReporter,
		installedRecipes: make(map[string]bool),
		prompter:         NewPrompter(),
	}
}

func createBundleInstaller() *BundleInstaller {
	mockStatusReporter := new(mockStatusReporter)
	mockRecipeInstaller := new(mockRecipeInstaller)

	return newBundleInstaller(context.Background(), &types.DiscoveryManifest{}, mockRecipeInstaller, mockStatusReporter)
}

func (bi *BundleInstaller) withStatusReporter(sr StatusReporter) *BundleInstaller {
	bi.statusReporter = sr
	return bi
}

func (bi *BundleInstaller) withRecipeInstaller(ri RecipeInstaller) *BundleInstaller {
	bi.recipeInstaller = ri
	return bi
}

func (bi *BundleInstaller) withPrompter(p Prompter) *BundleInstaller {
	bi.prompter = p
	return bi
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
	mockRecipeInstaller := new(mockRecipeInstaller)
	mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))

	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()

	bundleInstaller := createBundleInstaller().withStatusReporter(mockStatusReporter).withRecipeInstaller(mockRecipeInstaller)

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
	}

	bundleInstaller.InstallContinueOnError(&bundle, true)
	mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
}

func TestInstallContinueOnErrorReturnsImmediatelyWhenNoIsEntered(t *testing.T) {
	mockPrompter := ux.NewMockPrompter()
	mockPrompter.PromptYesNoVal = false
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	mockRecipeInstaller := new(mockRecipeInstaller)
	mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))
	bundleInstaller := createBundleInstaller().withPrompter(mockPrompter).withRecipeInstaller(mockRecipeInstaller).withStatusReporter(mockStatusReporter)

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
	}

	bundleInstaller.InstallContinueOnError(&bundle, false)

	mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
	mockStatusReporter.AssertNumberOfCalls(t, "ReportStatus", 1)
}

func TestInstallContinueOnErrorIgnoresUxPromptIfBundleIsAdditionalTargeted(t *testing.T) {
	mockPrompter := ux.NewMockPrompter()
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	mockRecipeInstaller := new(mockRecipeInstaller)
	mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(mock.AnythingOfType("String"), mock.AnythingOfType("error"))
	bundleInstaller := createBundleInstaller().withPrompter(mockPrompter).withRecipeInstaller(mockRecipeInstaller).withStatusReporter(mockStatusReporter)

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
		Type: recipes.BundleTypes.ADDITIONALTARGETED,
	}

	bundleInstaller.InstallContinueOnError(&bundle, true)

	assert.Equal(t, 0, mockPrompter.PromptMultiSelectCallCount)
}

// TODO come back to this, not sure if the test makes sense
//func TestInstallContinueOnErrorReturnsInstallsWhenYesIsEntered(t *testing.T) {
//	createBundleInstaller()
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
	expectedError := errors.New("Kaboom " + time.Now().String())
	mockRecipeInstaller := new(mockRecipeInstaller)
	mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", expectedError)
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	bundleInstaller := createBundleInstaller().withRecipeInstaller(mockRecipeInstaller).withStatusReporter(mockStatusReporter)

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

	actualError := bundleInstaller.InstallStopOnError(&bundle, true)

	//Should stop on first recipe
	mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func TestInstallContinueOnErrorOnlyInstallsAvailableRecipesInBundle(t *testing.T) {
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	mockRecipeInstaller := new(mockRecipeInstaller)
	mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("great success", nil)

	bundleInstaller := createBundleInstaller().withRecipeInstaller(mockRecipeInstaller).withStatusReporter(mockStatusReporter)

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

	bundleInstaller.InstallContinueOnError(&bundle, true)

	mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.True(t, bundleInstaller.installedRecipes["recipe1"])
	assert.False(t, bundleInstaller.installedRecipes["recipe2"])
	assert.False(t, bundleInstaller.installedRecipes["recipe3"])
}

func TestInstallContinueOnErrorKeepsInstalling(t *testing.T) {
	mockedRecipeInstaller := new(mockRecipeInstaller)
	mockedRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("great success", nil)
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()

	bundleInstaller := createBundleInstaller().withRecipeInstaller(mockedRecipeInstaller).withStatusReporter(mockStatusReporter)
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

	bundleInstaller.InstallContinueOnError(&bundle, true)

	//Should try both recipes
	mockedRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 2)
}

func TestReportsStatusHasSingleStatusWhenStatusNotAvailable(t *testing.T) {
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	bundleInstaller := createBundleInstaller().withStatusReporter(mockStatusReporter)
	bundle := givenBundle(types.InfraAgentRecipeName)
	expectedStatus := execution.RecipeStatusTypes.RECOMMENDED
	bundle.BundleRecipes[0].AddDetectionStatus(expectedStatus, 0)

	bundleInstaller.reportBundleStatus(bundle)

	assert.Equal(t, expectedStatus, bundle.BundleRecipes[0].DetectedStatuses[0].Status)
	assert.Equal(t, 1, len(bundle.BundleRecipes[0].DetectedStatuses))
}

func TestReportsStatusHasDetectedAndAvailableWhenStatusIsAvailable(t *testing.T) {
	mockStatusReporter := new(mockStatusReporter)
	mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	bundleInstaller := createBundleInstaller().withStatusReporter(mockStatusReporter)
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, 0)

	bundleInstaller.reportBundleStatus(bundle)

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
