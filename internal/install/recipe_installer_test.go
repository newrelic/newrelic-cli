package install

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/diagnose"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
)

func TestConnectToPlatformShouldSuccess(t *testing.T) {
	var expected error
	pi := ux.NewSpinnerProgressIndicator()

	recipeInstall := NewRecipeInstallBuilder().WithConfigValidatorError(expected).WithProgressIndicator(pi).Build()

	err := recipeInstall.connectToPlatform()
	assert.NoError(t, err)
}

func TestConnectToPlatformShouldReturnError(t *testing.T) {
	expected := errors.New("Failing to connect to platform")
	pi := ux.NewSpinnerProgressIndicator()

	recipeInstall := NewRecipeInstallBuilder().WithConfigValidatorError(expected).WithProgressIndicator(pi).Build()

	actual := recipeInstall.connectToPlatform()
	assert.Error(t, actual)
	assert.Equal(t, expected.Error(), actual.Error())
}

func TestConnectToPlatformShouldReturnPaymentRequiredError(t *testing.T) {
	expected := nrErrors.NewPaymentRequiredError()
	pi := ux.NewSpinnerProgressIndicator()

	recipeInstall := NewRecipeInstallBuilder().WithConfigValidatorError(expected).WithProgressIndicator(pi).Build()

	actual := recipeInstall.connectToPlatform()
	assert.Error(t, actual)
	assert.IsType(t, &nrErrors.PaymentRequiredError{}, actual)
}

func TestConnectToPlatformErrorShouldReportConnectionError(t *testing.T) {
	expected := diagnose.ConnectionError{
		Err: errors.New("Connection Failed"),
	}

	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithConfigValidatorError(expected).Build()

	actual := recipeInstall.Install()
	assert.Error(t, actual)
	assert.IsType(t, diagnose.ConnectionError{}, actual)
	assert.Equal(t, 1, statusReporter.InstallCompleteCallCount, "Install Completed")
	assert.True(t, strings.Contains(statusReporter.InstallCompleteErr.Error(), expected.Error()))
}

func TestInstallWithFailDiscoveryReturnsError(t *testing.T) {
	expected := errors.New("Some Discover error")
	recipeInstall := NewRecipeInstallBuilder().WithDiscovererError(expected).Build()

	actual := recipeInstall.Install()

	assert.Error(t, actual)
	assert.True(t, strings.Contains(actual.Error(), expected.Error()))
}

func TestInstallWithInvalidDiscoveryResultReturnsError(t *testing.T) {
	expected := errors.New("some discovery validation error")

	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithDiscovererValidatorError(expected).Build()
	actual := recipeInstall.Install()

	assert.Error(t, actual)
	assert.Equal(t, 1, statusReporter.DiscoveryCompleteCallCount)
	assert.True(t, strings.Contains(actual.Error(), expected.Error()))
}

func TestInstallGuidedShouldSkipCoreInstall(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).withShouldInstallCore(func() bool { return false }).WithStatusReporter(statusReporter).Build()
	err := recipeInstall.Install()

	assert.Equal(t, "no recipes were installed", err.Error(), "no recipe installed")
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 0, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 0, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 0, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe Installed")
}
func TestInstallGuidedShouldSkipCoreWhileInstallOthers(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}

	r2 := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name("Other").Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).
		WithRecipeDetectionResult(r2).withShouldInstallCore(func() bool { return false }).WithStatusReporter(statusReporter).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.Install()

	assert.NoError(t, err, "No error during install")
	assert.Equal(t, 2, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r2.Recipe.Name], "Recipe Installed")
	assert.Equal(t, 0, statusReporter.ReportInstalled[r.Recipe.Name], "Core Not Installed")
}

func TestInstallGuidedShouldNotSkipCoreInstall(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	r2 := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name("Other").Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).
		WithRecipeDetectionResult(r2).WithStatusReporter(statusReporter).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.Install()

	assert.NoError(t, err)
	assert.Equal(t, 2, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 2, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 2, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 2, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe1 Installed")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r2.Recipe.Name], "Recipe2 Installed")
}

func TestInstallGuidedCoreShouldStopOnError(t *testing.T) {
	installErr := errors.New("Install Error")
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	r2 := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.LoggingRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).WithRecipeDetectionResult(r2).
		WithStatusReporter(statusReporter).WithRecipeExecutionError(installErr).Build()
	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "Install Error"))
	assert.Equal(t, 2, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 2, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 1, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 0, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe1 Installed")
	assert.Equal(t, 0, statusReporter.ReportRecommended[r2.Recipe.Name], "Recipe2 Recommended")
}

func TestInstallTargetedInstallShouldInstallWithRecomendataion(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name("Other").Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	r2 := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name("Other2").Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).WithRecipeDetectionResult(r2).
		WithTargetRecipeName("Other").WithStatusReporter(statusReporter).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.Install()

	assert.NoError(t, err)
	assert.Equal(t, 2, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 1, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe1 Installed")
	assert.Equal(t, 1, statusReporter.ReportRecommended[r2.Recipe.Name], "Recipe2 Recommended")
}

func TestInstallTargetedInstallShouldInstallCoreIfCoreWasSkipped(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).withShouldInstallCore(func() bool { return false }).
		WithTargetRecipeName(types.InfraAgentRecipeName).WithStatusReporter(statusReporter).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.Install()

	assert.NoError(t, err)
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe1 Installed")
}

func TestInstallTargetedInstallWithoutRecipeShouldNotInstall(t *testing.T) {
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName("Other").WithStatusReporter(statusReporter).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.Equal(t, "no recipes were installed", err.Error())
	assert.Equal(t, 0, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 0, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 0, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
}

func TestInstallTargetedInstallWithOneUnsupportedOneInstalledShouldError(t *testing.T) {
	additionRecipeName := "additional"
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).WithTargetRecipeName(additionRecipeName).
		WithStatusReporter(statusReporter).Build()
	err := recipeInstall.install(context.TODO())

	assert.Error(t, err)
	assert.Equal(t, "one or more selected recipes could not be installed", err.Error())
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe Installed")
}

func TestInstallGuidededInstallAdditionalShouldInstall(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name("Other").Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).WithStatusReporter(statusReporter).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.install(context.TODO())

	assert.NoError(t, err, "No error during install")
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
	assert.Equal(t, 1, statusReporter.ReportInstalled[r.Recipe.Name], "Recipe Installed")
}

func TestPromptIfNotLatestCliVersionDoesNotLogMessagesOrErrorWhenVersionsMatch(t *testing.T) {
	getLatestCliVersionReleased = func(ctx context.Context) (string, error) {
		return "latest-version", nil
	}

	isLatestCliVersionInstalled = func(ctx context.Context, latestCliVersion string) (bool, error) {
		return true, nil
	}

	stdOut := captureLoggingOutput(func() {
		error := NewRecipeInstallBuilder().Build().promptIfNotLatestCLIVersion(MockContext{})
		assert.Nil(t, error)
	})

	assert.True(t, stdOut == "")
}

func TestPromptIfNotLatestCliVersionDisplaysErrorWhenLatestCliReleaseCannotBeDetermined(t *testing.T) {
	getLatestCliVersionReleased = func(ctx context.Context) (string, error) {
		return "", errors.New("couldn't fetch latest cli release")
	}

	stdOut := captureLoggingOutput(func() {
		error := NewRecipeInstallBuilder().Build().promptIfNotLatestCLIVersion(MockContext{})
		assert.Nil(t, error)
	})

	assert.True(t, strings.Contains(stdOut, "couldn't fetch latest cli release"))
}

func TestPromptIfNotLatestCliVersionDisplaysErrorWhenMostRecentInstalledCliCannotBeDetermined(t *testing.T) {
	getLatestCliVersionReleased = func(ctx context.Context) (string, error) {
		return "some-version", nil
	}

	isLatestCliVersionInstalled = func(ctx context.Context, latestCliVersion string) (bool, error) {
		return false, errors.New("something bad happened when comparing local to latest cli version")
	}

	stdOut := captureLoggingOutput(func() {
		error := NewRecipeInstallBuilder().Build().promptIfNotLatestCLIVersion(MockContext{})
		assert.Nil(t, error)
	})

	assert.True(t, strings.Contains(stdOut, "something bad happened when comparing local to latest cli version"))
}

func TestPromptIfNotLatestCliVersionErrorsIfNotLatestVersion(t *testing.T) {
	getLatestCliVersionReleased = func(ctx context.Context) (string, error) {
		return "some-version", nil
	}

	isLatestCliVersionInstalled = func(ctx context.Context, latestCliVersion string) (bool, error) {
		return false, nil
	}

	ri := NewRecipeInstallBuilder().Build()
	error := ri.promptIfNotLatestCLIVersion(MockContext{})

	assert.NotNil(t, error)
	assert.True(t, strings.Contains(error.Error(), "We need to update your New Relic CLI version to continue."))
	assert.True(t, ri.status.UpdateRequired)
}

func TestInstallWhenKeyFetchError(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	expected := errors.New("Some error")
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithLicenseKeyFetchResult(expected).WithRecipeDetectionResult(r).Build()
	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, statusReporter.RecipeInstallingCallCount, "Installed Count")
	assert.Equal(t, 1, statusReporter.InstallCompleteCallCount, "Install Complete Call Count")
}

func TestInstallWhenRecipeVarProviderError(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	expected := errors.New("Some error")
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithRecipeVarValues(nil, expected).WithRecipeDetectionResult(r).Build()
	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.Equal(t, expected, err)
	assert.Equal(t, 0, statusReporter.RecipeInstallingCallCount, "Installed Count")
	assert.Equal(t, 1, statusReporter.InstallCompleteCallCount, "Install Complete Call Count")
}

func TestInstallGuidedWhenInstallFails(t *testing.T) {
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	expected := errors.New("Some error")
	vars := map[string]string{}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).
		WithRecipeDetectionResult(r).WithRecipeVarValues(vars, nil).WithRecipeExecutionError(expected).Build()
	recipeInstall.AssumeYes = true

	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.True(t, vars["assumeYes"] == "true")
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 1, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
}

func TestInstallGuidedWhenGoTaskFails(t *testing.T) {
	expected := types.NewGoTaskGeneralError(errors.New("Some error"))
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	vars := map[string]string{}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).
		WithRecipeDetectionResult(r).WithRecipeVarValues(vars, nil).WithRecipeExecutionError(expected).Build()
	recipeInstall.AssumeYes = true

	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.True(t, vars["assumeYes"] == "true")
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 1, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")
}

func TestInstallWhenInstallIsCancelled(t *testing.T) {
	expected := types.ErrInterrupt
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name(types.InfraAgentRecipeName).Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	vars := map[string]string{}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).
		WithRecipeDetectionResult(r).WithRecipeVarValues(vars, nil).WithRecipeExecutionError(expected).Build()
	recipeInstall.AssumeYes = true

	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.True(t, vars["assumeYes"] == "true")
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 1, statusReporter.InstallCanceledCallCount, "Cancelled Count")
}

func TestInstallWhenInstallIsUnsupported(t *testing.T) {
	expected := &types.UnsupportedOperatingSystemError{Err: errors.New("Unsupported")}
	r := &recipes.RecipeDetectionResult{
		Recipe: recipes.NewRecipeBuilder().Name("Other").Build(),
		Status: execution.RecipeStatusTypes.AVAILABLE,
	}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithRecipeDetectionResult(r).WithStatusReporter(statusReporter).WithRecipeExecutionError(expected).Build()
	recipeInstall.AssumeYes = true
	err := recipeInstall.Install()

	assert.Error(t, err)
	assert.Equal(t, 1, statusReporter.RecipeDetectedCallCount, "Detection Count")
	assert.Equal(t, 1, statusReporter.RecipeAvailableCallCount, "Available Count")
	assert.Equal(t, 1, statusReporter.RecipeInstallingCallCount, "Installing Count")
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount, "Failed Count")
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount, "Unsupported Count")
	assert.Equal(t, 0, statusReporter.RecipeInstalledCallCount, "InstalledCount")
	assert.Equal(t, 0, statusReporter.RecipeRecommendedCallCount, "Recommendation Count")
	assert.Equal(t, 0, statusReporter.RecipeSkippedCallCount, "Skipped Count")
	assert.Equal(t, 0, statusReporter.RecipeCanceledCallCount, "Cancelled Count")

}

func TestExecuteAndValidateWithProgressWhenInstallWithNoValidationMethod(t *testing.T) {
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).Build()

	entityGUID, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.NoError(t, err)
	assert.Equal(t, "", entityGUID)
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount)
}

func TestExecuteAndValidateRecipeWithAllMethodWithNoValidationMethods(t *testing.T) {
	recipeInstall := NewRecipeInstallBuilder().Build()

	entityGUID, err := recipeInstall.validateRecipeViaAllMethods(context.TODO(), recipes.NewRecipeBuilder().Name("").Build(), &types.DiscoveryManifest{}, nil, false)

	assert.NoError(t, err)
	assert.Equal(t, "", entityGUID)
}

func TestExecuteAndValidateRecipeWithAllMethodWithValidationMethods(t *testing.T) {
	recipeInstall := NewRecipeInstallBuilder().WithOutput("{\"EntityGuid\":\"abcd\"}").Build()

	entityGUID, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.NoError(t, err)
	assert.Equal(t, "abcd", entityGUID)
}

func TestExecuteAndValidateRecipeWithAllMethodWithAgentValidatorError(t *testing.T) {
	expected := errors.New("Some error")
	recipeInstall := NewRecipeInstallBuilder().WithAgentValidationError(expected).Build()
	recipe := recipes.NewRecipeBuilder().Name("").Build()
	recipe.ValidationURL = "http://url.com"

	_, err := recipeInstall.validateRecipeViaAllMethods(context.TODO(), recipe, &types.DiscoveryManifest{}, nil, false)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.True(t, strings.Contains(err.Error(), "no validation was successful.  most recent validation error"))
}

func TestExecuteAndValidateRecipeWithAllMethodWithRecipeValidationError(t *testing.T) {
	expected := errors.New("Some error")
	recipeInstall := NewRecipeInstallBuilder().WithRecipeValidationError(expected).Build()
	recipe := recipes.NewRecipeBuilder().Name("").Build()
	recipe.ValidationNRQL = "FROM SOMETHING"

	_, err := recipeInstall.validateRecipeViaAllMethods(context.TODO(), recipe, &types.DiscoveryManifest{}, nil, false)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.True(t, strings.Contains(err.Error(), "no validation was successful.  most recent validation error"))
}

func TestExecuteAndValidateWithProgressWhenPostValidationFailed(t *testing.T) {
	expected := errors.New("Some error")
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithRecipeValidationError(expected).Build()
	recipe := recipes.NewRecipeBuilder().Name("").Build()
	recipe.ValidationNRQL = "FROM SOMETHING"

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipe, false)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.True(t, strings.Contains(err.Error(), "encountered an error while validating receipt of data for"))
	assert.Equal(t, 1, statusReporter.RecipeFailedCallCount)
	assert.Equal(t, 0, statusReporter.InstallCanceledCallCount)
}

func TestExecuteAndValidateWithProgressWhenSucceed(t *testing.T) {
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).Build()
	recipe := recipes.NewRecipeBuilder().Name("").Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipe, false)

	assert.NoError(t, err)
	assert.Equal(t, 0, statusReporter.RecipeFailedCallCount)
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount)
}

func TestReportUnSupportTargetRecipeWithBadRecipeName(t *testing.T) {
	targetRecipe := "target"
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(targetRecipe).WithStatusReporter(statusReporter).Build()
	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{}, nil
	}, &types.DiscoveryManifest{})
	bundle := &recipes.Bundle{}

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo)
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount)
}

func TestReportUnSupportTargetRecipeWithoutTarget(t *testing.T) {
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).Build()
	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{}, nil
	}, &types.DiscoveryManifest{})
	bundle := &recipes.Bundle{}

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo)
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount)
}

func TestReportUnSupportTargetRecipeWithBundleContainRecipe(t *testing.T) {
	targetRecipe := "target"
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(targetRecipe).WithStatusReporter(statusReporter).Build()
	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{}, nil
	}, &types.DiscoveryManifest{})
	bundle := &recipes.Bundle{}
	recipe := &recipes.BundleRecipe{Recipe: recipes.NewRecipeBuilder().Name(targetRecipe).Build()}
	bundle.AddRecipe(recipe)

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo)
	assert.Equal(t, 0, statusReporter.RecipeUnsupportedCallCount)
}

func TestReportUnSupportTargetRecipeWithUnsupportForPlatform(t *testing.T) {
	targetRecipe := "target"
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(targetRecipe).WithStatusReporter(statusReporter).Build()
	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{
			recipes.NewRecipeBuilder().Name(targetRecipe).Build(),
		}, nil
	}, &types.DiscoveryManifest{})
	bundle := &recipes.Bundle{}

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo)
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount)
}

func TestRecipeInstallerShouldGetEntityGuidFromRecipeExecution(t *testing.T) {
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithOutput("{\"EntityGuid\":\"abcd\"}").Build()
	recipe := recipes.NewRecipeBuilder().Name("").Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipe, false)

	assert.NoError(t, err)
	assert.Equal(t, 1, statusReporter.RecipeInstalledCallCount)
	assert.Equal(t, "abcd", statusReporter.GUIDs[0])
}

func TestIsTargetInstallRecipeShouldFindTarget(t *testing.T) {
	targetRecipe := "target"
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(targetRecipe).Build()

	actual := recipeInstall.isTargetInstallRecipe(targetRecipe)

	assert.True(t, actual)
}

func TestIsTargetInstallRecipeShouldNotFindTarget(t *testing.T) {
	recipeInstall := NewRecipeInstallBuilder().Build()

	actual := recipeInstall.isTargetInstallRecipe("target")

	assert.False(t, actual)
}

func TestWhenSingleInstallRunningErrorOnMultiple(t *testing.T) {
	recipeInstall := NewRecipeInstallBuilder().WithRunningProcess("env=123 newrelic install", "newrelic").WithRunningProcess("env=456 newrelic install", "newrelic").Build()
	err := recipeInstall.Install()
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "only 1 newrelic install command can run at one time"))
}

func TestWhenSingleInstallRunningNoError(t *testing.T) {
	recipeInstall := NewRecipeInstallBuilder().WithRunningProcess("env=123 newrelic install", "newrelic").Build()
	err := recipeInstall.Install()
	if err != nil {
		assert.False(t, strings.Contains(err.Error(), "only 1 newrelic install command can run at one time"))
	}
}

func captureLoggingOutput(f func()) string {
	var buf bytes.Buffer
	existingLogger := config.Logger
	existingLogger.SetOutput(&buf)
	existingLogger.SetLevel(logrus.DebugLevel)
	f()
	existingLogger.SetOutput(os.Stderr)
	return buf.String()
}
