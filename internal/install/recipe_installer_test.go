package install

import (
	"bytes"
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/config"

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

func TestInstallWithFailDiscoveryReturnsError(t *testing.T) {
	expected := errors.New("Some Discover error")
	recipeInstall := NewRecipeInstallBuilder().WithDiscovererError(expected).Build()

	actual := recipeInstall.install(context.TODO())

	assert.Error(t, actual)
	assert.True(t, strings.Contains(actual.Error(), expected.Error()))
}

func TestInstallWithInvalidDiscoveryResultReturnsError(t *testing.T) {
	expected := errors.New("some discovery validation error")

	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithDiscovererValidatorError(expected).Build()
	actual := recipeInstall.install(context.TODO())

	assert.Error(t, actual)
	assert.Equal(t, 1, statusReporter.DiscoveryCompleteCallCount)
	assert.True(t, strings.Contains(actual.Error(), expected.Error()))
}

func TestInstallShouldSkipCoreInstall(t *testing.T) {
	bundler := NewBundlerBuilder().WithCoreRecipe("Core").Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithBundler(bundler).withShouldInstallCore(func() bool { return false }).WithBundleInstaller(bundleInstaller).Build()
	coreBundle := bundler.CreateCoreBundle()

	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 1, len(coreBundle.BundleRecipes))
	assert.False(t, bundleInstaller.installedRecipes[coreBundle.BundleRecipes[0].Recipe.Name])
}

func TestInstallShouldNotSkipCoreInstall(t *testing.T) {
	bundler := NewBundlerBuilder().WithCoreRecipe("Core").Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	coreBundle := bundler.CreateCoreBundle()
	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 1, len(coreBundle.BundleRecipes))
	assert.True(t, bundleInstaller.installedRecipes[coreBundle.BundleRecipes[0].Recipe.Name])
}

func TestInstallCoreShouldStopOnError(t *testing.T) {
	bundler := NewBundlerBuilder().WithCoreRecipe("Core").Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	coreBundle := bundler.CreateCoreBundle()
	bundleInstaller.Error = errors.New("Install Error")
	err := recipeInstall.install(context.TODO())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "Install Error")
	assert.Equal(t, 1, len(coreBundle.BundleRecipes))
	assert.True(t, len(bundleInstaller.installedRecipes) == 0)
}

func TestInstallTargetInstallShouldInstall(t *testing.T) {
	additionRecipeName := "additional"
	bundler := NewBundlerBuilder().WithAdditionalRecipe(additionRecipeName).Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(additionRecipeName).WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	additionalBundle := bundler.CreateAdditionalTargetedBundle([]string{additionRecipeName})
	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 1, len(additionalBundle.BundleRecipes))
	assert.True(t, bundleInstaller.installedRecipes[additionalBundle.BundleRecipes[0].Recipe.Name])
}
func TestInstallTargetInstallShouldNotInstallCoreIfCoreWasNotSkipped(t *testing.T) {
	additionRecipeName := types.InfraAgentRecipeName
	bundler := NewBundlerBuilder().WithAdditionalRecipe(additionRecipeName).Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(additionRecipeName).WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 0, len(bundleInstaller.installedRecipes))
}

func TestInstallTargetInstallShouldInstallCoreIfCoreWasSkipped(t *testing.T) {
	additionRecipeName := types.InfraAgentRecipeName
	bundler := NewBundlerBuilder().WithAdditionalRecipe(additionRecipeName).Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithBundler(bundler).WithTargetRecipeName(additionRecipeName).withShouldInstallCore(func() bool { return false }).WithBundleInstaller(bundleInstaller).Build()
	additionalBundle := bundler.CreateAdditionalTargetedBundle([]string{additionRecipeName})
	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 1, len(additionalBundle.BundleRecipes))
	assert.True(t, bundleInstaller.installedRecipes[additionalBundle.BundleRecipes[0].Recipe.Name])
}

func TestInstallTargetInstallWithoutRecipeShouldNotInstall(t *testing.T) {
	additionRecipeName := "additional"
	bundler := NewBundlerBuilder().Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(additionRecipeName).WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	additionalBundle := bundler.CreateAdditionalTargetedBundle([]string{additionRecipeName})
	err := recipeInstall.install(context.TODO())

	assert.Error(t, err)
	assert.Equal(t, "no recipes were installed", err.Error())
	assert.Equal(t, 0, len(additionalBundle.BundleRecipes))
	assert.Equal(t, 0, len(bundleInstaller.installedRecipes))
}

func TestInstallTargetInstallWithOneUnsupportedOneInstalledShouldError(t *testing.T) {
	additionRecipeName := "additional"
	bundler := NewBundlerBuilder().Build()
	bundleInstaller := NewMockBundleInstaller()
	bundleInstaller.installedRecipes["test"] = true

	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(additionRecipeName).WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	additionalBundle := bundler.CreateAdditionalTargetedBundle([]string{additionRecipeName})
	err := recipeInstall.install(context.TODO())

	assert.Error(t, err)
	assert.Equal(t, "one or more selected recipes could not be installed", err.Error())
	assert.Equal(t, 0, len(additionalBundle.BundleRecipes))
	assert.Equal(t, 1, len(bundleInstaller.installedRecipes))
}

func TestInstallGuidedInstallAdditionalShouldInstall(t *testing.T) {
	additionRecipeName := "additional"
	bundler := NewBundlerBuilder().WithAdditionalRecipe(additionRecipeName).Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	additionalBundle := bundler.CreateAdditionalGuidedBundle()
	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 1, len(additionalBundle.BundleRecipes))
	assert.Equal(t, 1, len(bundleInstaller.installedRecipes))
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

func TestExecuteAndValidateWithProgressWhenKeyFetchError(t *testing.T) {
	expected := errors.New("Some error")
	recipeInstall := NewRecipeInstallBuilder().WithLicenseKeyFetchResult(expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), nil, recipes.NewRecipeBuilder().Name("").Build(), false)

	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestExecuteAndValidateWithProgressWhenRecipeVarProviderError(t *testing.T) {
	expected := errors.New("Some error")
	recipeInstall := NewRecipeInstallBuilder().WithRecipeVarValues(nil, expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), false)

	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestExecuteAndValidateWithProgressWhenInstallFails(t *testing.T) {
	expected := errors.New("Some error")
	vars := map[string]string{}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithRecipeVarValues(vars, nil).WithRecipeExecutionResult(expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.Error(t, err)
	assert.True(t, vars["assumeYes"] == "true")
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.Equal(t, 1, statusReporter.RecipeFailedCallCount)
}

func TestExecuteAndValidateWithProgressWhenInstallGoTaskFails(t *testing.T) {
	expected := types.NewGoTaskGeneralError(errors.New("Some error"))
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithRecipeExecutionResult(expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.Equal(t, 1, statusReporter.RecipeFailedCallCount)
}

func TestExecuteAndValidateWithProgressWhenInstallCancelled(t *testing.T) {
	expected := types.ErrInterrupt
	recipeInstall := NewRecipeInstallBuilder().WithRecipeExecutionResult(expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
}

func TestExecuteAndValidateWithProgressWhenInstallUnsupported(t *testing.T) {
	expected := &types.UnsupportedOperatingSystemError{Err: errors.New("Unsupported")}
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).WithRecipeExecutionResult(expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount)
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

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo, recipeInstall.RecipeNames)
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount)
}

func TestReportUnSupportTargetRecipeWithoutTarget(t *testing.T) {
	statusReporter := execution.NewMockStatusReporter()
	recipeInstall := NewRecipeInstallBuilder().WithStatusReporter(statusReporter).Build()
	repo := recipes.NewRecipeRepository(func() ([]*types.OpenInstallationRecipe, error) {
		return []*types.OpenInstallationRecipe{}, nil
	}, &types.DiscoveryManifest{})
	bundle := &recipes.Bundle{}

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo, recipeInstall.RecipeNames)
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

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo, recipeInstall.RecipeNames)
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

	recipeInstall.reportUnsupportedTargetedRecipes(bundle, repo, recipeInstall.RecipeNames)
	assert.Equal(t, 1, statusReporter.RecipeUnsupportedCallCount)
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
