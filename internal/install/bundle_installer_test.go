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
	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.UNSUPPORTED)

	test.BundleInstaller.InstallContinueOnError(test.bundle, true)
	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
}

func TestInstallContinueOnErrorReturnsImmediatelyWhenNoIsEntered(t *testing.T) {
	test := createBundleInstallerTest().withPrompterYesNoVal(false).withRecipeInstallerError()
	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.UNSUPPORTED)

	test.BundleInstaller.InstallContinueOnError(test.bundle, false)

	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 0)
	test.mockStatusReporter.AssertNumberOfCalls(t, "ReportStatus", 1)
}

func TestInstallContinueOnErrorIgnoresUxPromptIfBundleIsAdditionalTargeted(t *testing.T) {
	test := createBundleInstallerTest().withRecipeInstallerError()
	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.UNSUPPORTED)

	test.BundleInstaller.InstallContinueOnError(test.bundle, true)

	assert.Equal(t, 0, test.mockPrompter.PromptMultiSelectCallCount)
}

func TestInstallContinueOnErrorInstallAllWhenErroring(t *testing.T) {
	test := createBundleInstallerTest().withPrompterYesNoVal(true).withRecipeInstallerError()
	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.AVAILABLE)
	test.addRecipeToBundle("recipe2", execution.RecipeStatusTypes.AVAILABLE)

	test.BundleInstaller.InstallContinueOnError(test.bundle, false)

	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 2)
}

func TestInstallStopsOnErrorActuallyErrors(t *testing.T) {
	expectedError := errors.New("Kaboom " + time.Now().String())
	test := createBundleInstallerTest().withRecipeInstallerErrorWithMessage(expectedError)

	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.AVAILABLE)
	test.addRecipeToBundle("recipe2", execution.RecipeStatusTypes.AVAILABLE)

	actualError := test.BundleInstaller.InstallStopOnError(test.bundle, true)

	//Should stop on first recipe
	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.Equal(t, expectedError.Error(), actualError.Error())
}

func TestInstallContinueOnErrorOnlyInstallsAvailableRecipesInBundle(t *testing.T) {
	test := createBundleInstallerTest().withRecipeInstallerSuccess()

	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.AVAILABLE)
	test.addRecipeToBundle("recipe2", execution.RecipeStatusTypes.UNSUPPORTED)
	test.addRecipeToBundle("recipe3", execution.RecipeStatusTypes.DETECTED)

	test.BundleInstaller.InstallContinueOnError(test.bundle, true)

	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 1)
	assert.True(t, test.BundleInstaller.installedRecipes["recipe1"])
	assert.False(t, test.BundleInstaller.installedRecipes["recipe2"])
	assert.False(t, test.BundleInstaller.installedRecipes["recipe3"])
}

func TestInstallContinueOnErrorKeepsInstallingWhenNotErroring(t *testing.T) {
	test := createBundleInstallerTest().withRecipeInstallerSuccess()

	test.addRecipeToBundle("recipe1", execution.RecipeStatusTypes.AVAILABLE)
	test.addRecipeToBundle("recipe2", execution.RecipeStatusTypes.AVAILABLE)

	test.BundleInstaller.InstallContinueOnError(test.bundle, true)
	//Should try both recipes
	test.mockRecipeInstaller.AssertNumberOfCalls(t, "executeAndValidateWithProgress", 2)
}

func TestReportsStatusHasSingleStatusWhenStatusNotAvailable(t *testing.T) {
	test := createBundleInstallerTest()
	test.addRecipeToBundle(types.InfraAgentRecipeName, execution.RecipeStatusTypes.RECOMMENDED)
	expectedStatus := execution.RecipeStatusTypes.RECOMMENDED

	test.BundleInstaller.reportBundleStatus(test.bundle)

	assert.Equal(t, expectedStatus, test.bundle.BundleRecipes[0].DetectedStatuses[0].Status)
	assert.Equal(t, 1, len(test.bundle.BundleRecipes[0].DetectedStatuses))
}

func TestReportsStatusHasDetectedAndAvailableWhenStatusIsAvailable(t *testing.T) {
	test := createBundleInstallerTest()
	test.addRecipeToBundle(types.InfraAgentRecipeName, execution.RecipeStatusTypes.AVAILABLE)

	test.BundleInstaller.reportBundleStatus(test.bundle)

	assert.True(t, test.bundle.BundleRecipes[0].HasStatus(execution.RecipeStatusTypes.AVAILABLE))
	assert.True(t, test.bundle.BundleRecipes[0].HasStatus(execution.RecipeStatusTypes.DETECTED))
	assert.Equal(t, 2, len(test.bundle.BundleRecipes[0].DetectedStatuses))
}

type BundleInstallerTest struct {
	BundleInstaller     *BundleInstaller
	mockStatusReporter  *mockStatusReporter
	mockRecipeInstaller *mockRecipeInstaller
	mockPrompter        *ux.MockPrompter
	bundle              *recipes.Bundle
}

func createBundleInstallerTest() *BundleInstallerTest {
	i := &BundleInstallerTest{
		mockStatusReporter:  new(mockStatusReporter),
		mockRecipeInstaller: new(mockRecipeInstaller),
		mockPrompter:        ux.NewMockPrompter(),
		bundle:              &recipes.Bundle{},
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

func (bi *BundleInstallerTest) addRecipeToBundle(name string, status execution.RecipeStatusType) *BundleInstallerTest {

	br := &recipes.BundleRecipe{
		Recipe: recipes.NewRecipeBuilder().Name(name).Build(),
	}

	br.AddDetectionStatus(status, 0)
	bi.bundle.AddRecipe(br)
	return bi
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
