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
)

func TestConnectToPlatformShouldSuccess(t *testing.T) {
	var expected error
	pi := ux.NewSpinnerProgressIndicator()

	recipeInstall := NewRecipeInstallBuilder().WithConfigValidatorError(expected).WithProgressIndicator(pi).Build()

	stdOut := captureLoggingOutput(func() {
		err := recipeInstall.connectToPlatform()
		assert.NoError(t, err)
	})
	assert.True(t, strings.Contains(stdOut, "Connecting"))
	assert.True(t, strings.Contains(stdOut, "Connected"))
}
func TestConnectToPlatformShouldRetrunError(t *testing.T) {
	expected := errors.New("Failing to connect to platform")
	pi := ux.NewSpinnerProgressIndicator()

	recipeInstall := NewRecipeInstallBuilder().WithConfigValidatorError(expected).WithProgressIndicator(pi).Build()

	stdOut := captureLoggingOutput(func() {
		actual := recipeInstall.connectToPlatform()
		assert.Error(t, actual)
		assert.Equal(t, expected.Error(), actual.Error())
	})

	assert.True(t, strings.Contains(stdOut, "Connecting"))
	assert.True(t, strings.Contains(stdOut, "Fail"))
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

func TestInstallTargetInstallWithoutRecipeShouldNotInstall(t *testing.T) {

	additionRecipeName := "additional"
	bundler := NewBundlerBuilder().Build()
	bundleInstaller := NewMockBundleInstaller()
	recipeInstall := NewRecipeInstallBuilder().WithTargetRecipeName(additionRecipeName).WithBundler(bundler).WithBundleInstaller(bundleInstaller).Build()
	additionalBundle := bundler.CreateAdditionalTargetedBundle([]string{additionRecipeName})
	_ = recipeInstall.install(context.TODO())

	assert.Equal(t, 0, len(additionalBundle.BundleRecipes))
	assert.Equal(t, 0, len(bundleInstaller.installedRecipes))
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
	recipeInstall := NewRecipeInstallBuilder().WithRecipeVarValues(vars, nil).WithRecipeExecutionResult(expected).Build()

	_, err := recipeInstall.executeAndValidateWithProgress(context.TODO(), &types.DiscoveryManifest{}, recipes.NewRecipeBuilder().Name("").Build(), true)

	assert.Error(t, err)
	assert.True(t, vars["assumeYes"] == "true")
	assert.True(t, strings.Contains(err.Error(), expected.Error()))
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

func captureLoggingOutput(f func()) string {
	var buf bytes.Buffer
	existingLogger := config.Logger
	existingLogger.SetOutput(&buf)
	existingLogger.SetLevel(logrus.DebugLevel)
	f()
	existingLogger.SetOutput(os.Stderr)
	return buf.String()
}

// import (
// 	// "errors"
// 	// "net/url"
// 	// "os"
// 	// "reflect"
// 	// "testing"

// 	// "github.com/stretchr/testify/require"

// 	// "github.com/newrelic/newrelic-cli/internal/cli"
// 	"github.com/newrelic/newrelic-cli/internal/diagnose"
// 	"github.com/newrelic/newrelic-cli/internal/install/discovery"
// 	"github.com/newrelic/newrelic-cli/internal/install/execution"
// 	"github.com/newrelic/newrelic-cli/internal/install/recipes"
// 	// "github.com/newrelic/newrelic-cli/internal/install/types"
// 	"github.com/newrelic/newrelic-cli/internal/install/ux"
// 	"github.com/newrelic/newrelic-cli/internal/install/validation"
// )

// var (
// 	//testRecipeName = "test-recipe"
// 	//anotherTestRecipeName = "another-test-recipe"
// 	// testRecipeFile        = &types.OpenInstallationRecipe{
// 	// 	Name: testRecipeName,
// 	// }

// 	d               = discovery.NewMockDiscoverer()
// 	mv              = discovery.NewEmptyManifestValidator()
// 	f               = recipes.NewMockRecipeFetcher()
// 	e               = execution.NewMockRecipeExecutor()
// 	v               = validation.NewMockRecipeValidator()
// 	ff              = recipes.NewMockRecipeFileFetcher()
// 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// 	status          = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// 	p               = ux.NewMockPrompter()
// 	pi              = ux.NewMockProgressIndicator()
// 	sp              = ux.NewMockProgressIndicator()
// 	lkf             = NewMockLicenseKeyFetcher()
// 	cv              = diagnose.NewMockConfigValidator()
// 	rvp             = execution.NewRecipeVarProvider()
// 	av              = validation.NewAgentValidator()
// )

// // func TestInstall_DiscoveryComplete(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	statusReporter := execution.NewMockStatusReporter()
// // 	statusReporters = []execution.StatusSubscriber{statusReporter}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			DisplayName:    types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}

// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 1, statusReporter.DiscoveryCompleteCallCount)
// // }

// // func TestInstall_RecipeAvailable(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			DisplayName:    types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           types.LoggingRecipeName,
// // 			DisplayName:    types.LoggingRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeAvailableCallCount)
// // }

// // func TestInstall_RecipeInstalled(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			DisplayName:    types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           types.LoggingRecipeName,
// // 			DisplayName:    types.LoggingRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 			LogMatch: []types.OpenInstallationLogMatch{
// // 				{
// // 					Name: "docker log",
// // 					File: "/var/lib/docker/containers/*/*.log",
// // 				},
// // 			},
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // }

// // func TestInstall_RecipeFailed(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}

// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())

// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			DisplayName:    types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	rv := validation.NewMockRecipeValidator()
// // 	rv.ValidateErr = errors.New("validationErr")

// // 	i := RecipeInstall{ic, d, l, mv, f, e, rv, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.Error(t, err)
// // 	require.Equal(t, 1, rv.ValidateCallCount)
// // 	// Infra fails fast
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
// // }

// // func TestInstall_NonInfraRecipeFailed(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}

// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())

// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			DisplayName:    types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	rv := validation.NewMockRecipeValidator()
// // 	rv.ValidateErrs = []error{nil, errors.New("validationErr")}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, rv, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.Error(t, err)
// // 	require.Equal(t, 2, rv.ValidateCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
// // }

// // func TestInstall_AllRecipesFailed(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}

// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())

// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           anotherTestRecipeName,
// // 			DisplayName:    anotherTestRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	rv := validation.NewMockRecipeValidator()
// // 	rv.ValidateErr = errors.New("validationErr")

// // 	i := RecipeInstall{ic, d, l, mv, f, e, rv, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.Error(t, err)
// // 	require.Equal(t, 2, rv.ValidateCallCount)
// // 	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
// // }

// // func TestInstall_InstallStarted(t *testing.T) {
// // 	ic := types.InstallerContext{}

// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	_ = i.Install()
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallStartedCallCount)
// // }

// // func TestInstall_InstallComplete(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           types.LoggingRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // 	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).InstallCanceledCallCount)
// // 	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
// // }

// // func TestInstall_InstallCanceled(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesErr = types.ErrInterrupt

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.Error(t, err)
// // 	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCanceledCallCount)
// // }

// // func TestInstall_InstallCompleteError(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	rv := validation.NewMockRecipeValidator()
// // 	rv.ValidateErr = errors.New("test error")

// // 	i := RecipeInstall{ic, d, l, mv, f, e, rv, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.Error(t, err)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
// // }

// // func TestInstall_InstallCompleteError_NoFailureWhenAnyRecipeSucceeds(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           "badRecipe",
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	rv := validation.NewMockRecipeValidator()
// // 	rv.ValidateErrs = []error{
// // 		nil,
// // 		errors.New("testing error"),
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, rv, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.Error(t, err)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
// // }

// // func TestInstall_RecipeSkipped_MultiSelect(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.LoggingRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	mp := &ux.MockPrompter{
// // 		PromptYesNoVal:       true,
// // 		PromptMultiSelectVal: []string{testRecipeName},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, mp, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // }

// // func TestInstall_RecipeSkipped_AssumeYes(t *testing.T) {
// // 	ic := types.InstallerContext{
// // 		AssumeYes: true,
// // 	}

// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:        types.InfraAgentRecipeName,
// // 			DisplayName: "Infra Recipe",
// // 		},
// // 		{
// // 			Name:        types.LoggingRecipeName,
// // 			DisplayName: "Logging Recipe",
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    "test displayName",
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
// // 	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
// // 	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // }

// // func TestInstall_TargetedInstall_InstallsInfraAgent(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	v = validation.NewMockRecipeValidator()

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // }

// // func TestInstall_TargetedInstall_FilterAllButProvided(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{
// // 		RecipeNames: []string{testRecipeName},
// // 	}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           anotherTestRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	v = validation.NewMockRecipeValidator()

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.True(t, status.IsTargetedInstall())
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // }

// // func TestInstall_TargetedInstall_InstallsInfraAgentDependency(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{
// // 		RecipeNames: []string{testRecipeName},
// // 	}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 			Dependencies:   []string{types.InfraAgentRecipeName},
// // 		},
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.True(t, status.IsTargetedInstall())
// // 	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // }

// // func TestInstall_TargetedInstallInfraAgent_NoInfraAgentDuplicate(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{
// // 		RecipeNames: []string{types.InfraAgentRecipeName},
// // 	}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.True(t, status.IsTargetedInstall())
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // }

// // func TestInstall_TargetedInstall_SkipInfraDependency(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           "testRecipe",
// // 			ValidationNRQL: "testNrql",
// // 			Dependencies:   []string{types.InfraAgentRecipeName},
// // 		},
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	i := RecipeInstall{ic, d, l, mv, f, e, v, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // }

// // func TestInstall_GuidReport(t *testing.T) {
// // 	ic := types.InstallerContext{}
// // 	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
// // 	rf := recipes.NewRecipeFilterRunner(ic, status)
// // 	f = recipes.NewMockRecipeFetcher()
// // 	f.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		{
// // 			Name:           types.InfraAgentRecipeName,
// // 			DisplayName:    types.InfraAgentRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 		{
// // 			Name:           types.LoggingRecipeName,
// // 			DisplayName:    types.LoggingRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 			Dependencies:   []string{types.InfraAgentRecipeName},
// // 		},
// // 		{
// // 			Name:           testRecipeName,
// // 			DisplayName:    testRecipeName,
// // 			ValidationNRQL: "testNrql",
// // 		},
// // 	}

// // 	rv := validation.NewMockRecipeValidator()
// // 	rv.ValidateVal = "GUID"

// // 	i := RecipeInstall{ic, d, l, mv, f, e, rv, ff, status, p, pi, sp, lkf, cv, rvp, rf, av}
// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 3, rv.ValidateCallCount)
// // 	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
// // 	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
// // 	require.Equal(t, rv.ValidateVal, statusReporters[0].(*execution.MockStatusReporter).RecipeGUID[types.InfraAgentRecipeName])
// // 	require.Equal(t, rv.ValidateVal, statusReporters[0].(*execution.MockStatusReporter).RecipeGUID[testRecipeName])
// // 	require.Equal(t, status.CLIVersion, cli.Version())
// // 	require.Equal(t, 6, len(statusReporters[0].(*execution.MockStatusReporter).Durations))
// // 	for _, duration := range statusReporters[0].(*execution.MockStatusReporter).Durations {
// // 		require.Less(t, int64(0), duration)
// // 	}
// // }

// // func TestInstall_ShouldDetect_PreInstallDetected(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	reporters := []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	mockDiscoverer := discovery.NewMockDiscoverer()
// // 	installStatus := execution.NewInstallStatus(reporters, execution.NewPlatformLinkGenerator())
// // 	mockOsValidator := discovery.NewMockOsValidator()
// // 	mValidator := discovery.NewMockManifestValidator(mockOsValidator)

// // 	matchedProcess := mockProcess{
// // 		cmdline: "apache2",
// // 		name:    `apache2`,
// // 		pid:     int32(1234),
// // 	}

// // 	dm := types.DiscoveryManifest{
// // 		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
// // 	}

// // 	mockDiscoverer.DiscoveryManifest = &dm

// // 	infraRecipe := types.OpenInstallationRecipe{
// // 		Name:           "infrastructure-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 	}

// // 	testRecipe := types.OpenInstallationRecipe{
// // 		Name:           "php-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 		ProcessMatch:   []string{"apache2"},
// // 		PreInstall: types.OpenInstallationPreInstallConfiguration{
// // 			RequireAtDiscovery: `exit 132`,
// // 		},
// // 	}

// // 	rf := recipes.NewRecipeFilterRunner(ic, installStatus)
// // 	rFetcher := recipes.NewMockRecipeFetcher()

// // 	rFetcher.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		// Should be detected and installed
// // 		infraRecipe,

// // 		// Should be detected (but not installed) due to preinstall check exiting with 132 status code
// // 		testRecipe,
// // 	}

// // 	mrv := validation.NewMockRecipeValidator()
// // 	i := RecipeInstall{ic, mockDiscoverer, l, mValidator, rFetcher, e, mrv, ff, installStatus, p, pi, sp, lkf, cv, rvp, rf, av}

// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 2, reporters[0].(*execution.MockStatusReporter).RecipeDetectedCallCount)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeAvailableCallCount)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
// // }

// // func TestInstall_ShouldDetect_PreInstallOk(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	reporters := []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	mockDiscoverer := discovery.NewMockDiscoverer()
// // 	installStatus := execution.NewInstallStatus(reporters, execution.NewPlatformLinkGenerator())
// // 	mockOsValidator := discovery.NewMockOsValidator()
// // 	mValidator := discovery.NewMockManifestValidator(mockOsValidator)

// // 	matchedProcess := mockProcess{
// // 		cmdline: "apache2",
// // 		name:    `apache2`,
// // 		pid:     int32(1234),
// // 	}

// // 	dm := types.DiscoveryManifest{
// // 		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
// // 	}

// // 	mockDiscoverer.DiscoveryManifest = &dm

// // 	infraRecipe := types.OpenInstallationRecipe{
// // 		Name:           "infrastructure-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 	}

// // 	testRecipe := types.OpenInstallationRecipe{
// // 		Name:           "php-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 		ProcessMatch:   []string{"apache2"},
// // 		PreInstall: types.OpenInstallationPreInstallConfiguration{
// // 			RequireAtDiscovery: `exit 0`, // simulate successful preinstall check
// // 		},
// // 	}

// // 	rf := recipes.NewRecipeFilterRunner(ic, installStatus)
// // 	rFetcher := recipes.NewMockRecipeFetcher()

// // 	rFetcher.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		// Should be detected and installed
// // 		infraRecipe,

// // 		// Should be detected and installed
// // 		testRecipe,
// // 	}

// // 	mrv := validation.NewMockRecipeValidator()
// // 	i := RecipeInstall{ic, mockDiscoverer, l, mValidator, rFetcher, e, mrv, ff, installStatus, p, pi, sp, lkf, cv, rvp, rf, av}

// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 2, reporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 2, reporters[0].(*execution.MockStatusReporter).RecipeDetectedCallCount)
// // }

// // func TestInstall_ShouldDetect_ProcessMatch_NoScript(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	reporters := []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	mockDiscoverer := discovery.NewMockDiscoverer()
// // 	installStatus := execution.NewInstallStatus(reporters, execution.NewPlatformLinkGenerator())
// // 	mockOsValidator := discovery.NewMockOsValidator()
// // 	mValidator := discovery.NewMockManifestValidator(mockOsValidator)

// // 	matchedProcess := mockProcess{
// // 		cmdline: "apache2",
// // 		name:    `apache2`,
// // 		pid:     int32(1234),
// // 	}

// // 	dm := types.DiscoveryManifest{
// // 		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
// // 	}

// // 	mockDiscoverer.DiscoveryManifest = &dm

// // 	infraRecipe := types.OpenInstallationRecipe{
// // 		Name:           "infrastructure-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 	}

// // 	testRecipe := types.OpenInstallationRecipe{
// // 		Name:           "test-recipe",
// // 		ValidationNRQL: "testNrql",
// // 		ProcessMatch:   []string{"apache2"},
// // 	}

// // 	rf := recipes.NewRecipeFilterRunner(ic, installStatus)
// // 	rFetcher := recipes.NewMockRecipeFetcher()

// // 	rFetcher.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		// Should be detected and installed
// // 		infraRecipe,

// // 		// Should be detected and installed
// // 		testRecipe,
// // 	}

// // 	mrv := validation.NewMockRecipeValidator()
// // 	i := RecipeInstall{ic, mockDiscoverer, l, mValidator, rFetcher, e, mrv, ff, installStatus, p, pi, sp, lkf, cv, rvp, rf, av}

// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 2, reporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 2, reporters[0].(*execution.MockStatusReporter).RecipeDetectedCallCount)
// // }

// // func TestInstall_ShouldNotDetect_NoProcessMatch(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	reporters := []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	mockDiscoverer := discovery.NewMockDiscoverer()
// // 	installStatus := execution.NewInstallStatus(reporters, execution.NewPlatformLinkGenerator())
// // 	mockOsValidator := discovery.NewMockOsValidator()
// // 	mValidator := discovery.NewMockManifestValidator(mockOsValidator)

// // 	matchedProcess := mockProcess{
// // 		cmdline: "node",
// // 		name:    `node`,
// // 		pid:     int32(1234),
// // 	}

// // 	dm := types.DiscoveryManifest{
// // 		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
// // 	}

// // 	mockDiscoverer.DiscoveryManifest = &dm

// // 	infraRecipe := types.OpenInstallationRecipe{
// // 		Name:           "infrastructure-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 	}

// // 	testRecipe := types.OpenInstallationRecipe{
// // 		Name:           "test-recipe",
// // 		ValidationNRQL: "testNrql",
// // 		ProcessMatch:   []string{"apache2"}, // does not match mocked `node` process
// // 	}

// // 	rf := recipes.NewRecipeFilterRunner(ic, installStatus)
// // 	rFetcher := recipes.NewMockRecipeFetcher()

// // 	rFetcher.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		// Should be detected and installed
// // 		infraRecipe,

// // 		// Should NOT be detected and installed
// // 		testRecipe,
// // 	}

// // 	mrv := validation.NewMockRecipeValidator()
// // 	i := RecipeInstall{ic, mockDiscoverer, l, mValidator, rFetcher, e, mrv, ff, installStatus, p, pi, sp, lkf, cv, rvp, rf, av}

// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeDetectedCallCount)
// // }

// // func TestInstall_ShouldNotDetect_PreInstallError(t *testing.T) {
// // 	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
// // 	ic := types.InstallerContext{}
// // 	reporters := []execution.StatusSubscriber{execution.NewMockStatusReporter()}
// // 	mockDiscoverer := discovery.NewMockDiscoverer()
// // 	installStatus := execution.NewInstallStatus(reporters, execution.NewPlatformLinkGenerator())
// // 	mockOsValidator := discovery.NewMockOsValidator()
// // 	mValidator := discovery.NewMockManifestValidator(mockOsValidator)

// // 	matchedProcess := mockProcess{
// // 		cmdline: "apache2",
// // 		name:    `apache2`,
// // 		pid:     int32(1234),
// // 	}

// // 	dm := types.DiscoveryManifest{
// // 		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
// // 	}

// // 	mockDiscoverer.DiscoveryManifest = &dm

// // 	infraRecipe := types.OpenInstallationRecipe{
// // 		Name:           "infrastructure-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 	}

// // 	testRecipe := types.OpenInstallationRecipe{
// // 		Name:           "php-agent-installer",
// // 		ValidationNRQL: "testNrql",
// // 		ProcessMatch:   []string{"apache2"},
// // 		PreInstall: types.OpenInstallationPreInstallConfiguration{
// // 			RequireAtDiscovery: `exit 1`, // simulate misc error in preinstall check
// // 		},
// // 	}

// // 	rf := recipes.NewRecipeFilterRunner(ic, installStatus)
// // 	rFetcher := recipes.NewMockRecipeFetcher()

// // 	rFetcher.FetchRecipesVal = []types.OpenInstallationRecipe{
// // 		// Should be detected and installed
// // 		infraRecipe,

// // 		// Should NOT be detected and should NOT be installed due to error in preinstall check
// // 		testRecipe,
// // 	}

// // 	mrv := validation.NewMockRecipeValidator()
// // 	i := RecipeInstall{ic, mockDiscoverer, l, mValidator, rFetcher, e, mrv, ff, installStatus, p, pi, sp, lkf, cv, rvp, rf, av}

// // 	err := i.Install()
// // 	require.NoError(t, err)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
// // 	require.Equal(t, 1, reporters[0].(*execution.MockStatusReporter).RecipeDetectedCallCount)
// // }

// // func fetchRecipeFileFunc(recipeURL *url.URL) (*types.OpenInstallationRecipe, error) {
// // 	return testRecipeFile, nil
// // }

// // func loadRecipeFileFunc(filename string) (*types.OpenInstallationRecipe, error) {
// // 	return testRecipeFile, nil
// // }

// // type mockProcess struct {
// // 	cmdline string
// // 	name    string
// // 	pid     int32
// // }

// // func (p mockProcess) Name() (string, error) {
// // 	return p.name, nil
// // }

// // func (p mockProcess) Cmd() (string, error) {
// // 	return p.cmdline, nil
// // }

// // func (p mockProcess) PID() int32 {
// // 	return p.pid
// // }
