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
	test := createBundleInstallerTest().withRecipeInstallerError()

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
	}

	test.BundleInstaller.InstallContinueOnError(&bundle, true)
	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
}

func TestInstallContinueOnErrorReturnsImmediatelyWhenNoIsEntered(t *testing.T) {
	test := createBundleInstallerTest().withPrompterYesNoVal(false).withRecipeInstallerError()

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
	}

	test.BundleInstaller.InstallContinueOnError(&bundle, false)

	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
	test.mockStatusReporter.AssertNumberOfCalls(t, "ReportStatus", 1)
}

func TestInstallContinueOnErrorIgnoresUxPromptIfBundleIsAdditionalTargeted(t *testing.T) {
	test := createBundleInstallerTest().withRecipeInstallerError()

	bundle := recipes.Bundle{
		BundleRecipes: []*recipes.BundleRecipe{
			{
				Recipe:           recipes.NewRecipeBuilder().Name("recipe1").Build(),
				DetectedStatuses: []*recipes.DetectedStatusType{{Status: execution.RecipeStatusTypes.UNSUPPORTED}},
			},
		},
		Type: recipes.BundleTypes.ADDITIONALTARGETED,
	}

	test.BundleInstaller.InstallContinueOnError(&bundle, true)

	assert.Equal(t, 0, test.mockPrompter.PromptMultiSelectCallCount)
}

func TestInstallContinueOnErrorInstallAllWhenErroring(t *testing.T) {
	test := createBundleInstallerTest().withPrompterYesNoVal(true).withRecipeInstallerError()

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

	test.BundleInstaller.InstallContinueOnError(&bundle, false)

	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 2)
}

func TestInstallStopsOnErrorActuallyErrors(t *testing.T) {
	expectedError := errors.New("Kaboom " + time.Now().String())
	test := createBundleInstallerTest().withRecipeInstallerErrorWithMessage(expectedError)

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

	actualError := test.BundleInstaller.InstallStopOnError(&bundle, true)

	//Should stop on first recipe
	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func TestInstallContinueOnErrorOnlyInstallsAvailableRecipesInBundle(t *testing.T) {
	test := createBundleInstallerTest().withRecipeInstallerSuccess()

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

	test.BundleInstaller.InstallContinueOnError(&bundle, true)

	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.True(t, test.BundleInstaller.installedRecipes["recipe1"])
	assert.False(t, test.BundleInstaller.installedRecipes["recipe2"])
	assert.False(t, test.BundleInstaller.installedRecipes["recipe3"])
}

func TestInstallContinueOnErrorKeepsInstallingWhenNotErroring(t *testing.T) {
	test := createBundleInstallerTest().withRecipeInstallerSuccess()
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

	test.BundleInstaller.InstallContinueOnError(&bundle, true)

	//Should try both recipes
	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 2)
}

func TestReportsStatusHasSingleStatusWhenStatusNotAvailable(t *testing.T) {
	test := createBundleInstallerTest()
	bundle := givenBundle(types.InfraAgentRecipeName)
	expectedStatus := execution.RecipeStatusTypes.RECOMMENDED
	bundle.BundleRecipes[0].AddDetectionStatus(expectedStatus, 0)

	test.BundleInstaller.reportBundleStatus(bundle)

	assert.Equal(t, expectedStatus, bundle.BundleRecipes[0].DetectedStatuses[0].Status)
	assert.Equal(t, 1, len(bundle.BundleRecipes[0].DetectedStatuses))
}

func TestReportsStatusHasDetectedAndAvailableWhenStatusIsAvailable(t *testing.T) {
	test := createBundleInstallerTest()
	bundle := givenBundle(types.InfraAgentRecipeName)
	bundle.BundleRecipes[0].AddDetectionStatus(execution.RecipeStatusTypes.AVAILABLE, 0)

	test.BundleInstaller.reportBundleStatus(bundle)

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

type BundleInstallerTest struct {
	BundleInstaller     *BundleInstaller
	mockStatusReporter  *mockStatusReporter
	mockRecipeInstaller *mockRecipeInstaller
	mockPrompter        *ux.MockPrompter
}

func createBundleInstallerTest() *BundleInstallerTest {
	i := &BundleInstallerTest{
		mockStatusReporter:  new(mockStatusReporter),
		mockRecipeInstaller: new(mockRecipeInstaller),
		mockPrompter:        ux.NewMockPrompter(),
	}
	i.BundleInstaller = &BundleInstaller{
		ctx:              context.Background(),
		manifest:         &types.DiscoveryManifest{},
		recipeInstaller:  i.mockRecipeInstaller,
		statusReporter:   i.mockStatusReporter,
		installedRecipes: make(map[string]bool),
		prompter:         i.mockPrompter,
	}
	// Always stub status reporter usages
	i.withStatusReporter()
	return i
}

func (bi *BundleInstallerTest) withStatusReporter() *BundleInstallerTest {
	bi.mockStatusReporter.On("ReportStatus", mock.Anything, mock.Anything).Return()
	return bi
}

func (bi *BundleInstallerTest) withRecipeInstallerError() *BundleInstallerTest {
	return bi.withRecipeInstallerErrorWithMessage(errors.New("Nope, this is an error generated by a test"))
}

func (bi *BundleInstallerTest) withRecipeInstallerErrorWithMessage(e error) *BundleInstallerTest {
	bi.mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("A specific test error", e)
	return bi
}

func (bi *BundleInstallerTest) withRecipeInstallerSuccess() *BundleInstallerTest {
	bi.mockRecipeInstaller.On("executeAndValidateWithProgress", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("All good", nil)
	return bi
}

func (bi *BundleInstallerTest) withPrompterYesNoVal(val bool) *BundleInstallerTest {
	bi.mockPrompter.PromptYesNoVal = val
	return bi
}
