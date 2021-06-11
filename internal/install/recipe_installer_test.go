// +build unit

package install

import (
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
)

var (
	testRecipeName        = "test-recipe"
	anotherTestRecipeName = "another-test-recipe"
	testRecipeFile        = &types.OpenInstallationRecipe{
		Name: testRecipeName,
	}

	d               = discovery.NewMockDiscoverer()
	l               = discovery.NewMockFileFilterer()
	mv              = discovery.NewEmptyManifestValidator()
	f               = recipes.NewMockRecipeFetcher()
	e               = execution.NewMockRecipeExecutor()
	v               = validation.NewMockRecipeValidator()
	ff              = recipes.NewMockRecipeFileFetcher()
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status          = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	p               = ux.NewMockPrompter()
	pi              = ux.NewMockProgressIndicator()
	lkf             = NewMockLicenseKeyFetcher()
	cv              = diagnose.NewMockConfigValidator()
	rvp             = execution.NewRecipeVarProvider()
)

func TestNewRecipeInstaller_InstallerContextFields(t *testing.T) {
	ic := types.InstallerContext{
		RecipePaths:        []string{"testRecipePath"},
		RecipeNames:        []string{"testRecipeName"},
		SkipIntegrations:   true,
		SkipLoggingInstall: true,
		SkipApm:            true,
	}
	rf := recipes.NewRecipeFilterRunner(ic, status)

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}

	require.True(t, reflect.DeepEqual(ic, i.InstallerContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := types.InstallerContext{}
	rf := recipes.NewRecipeFilterRunner(ic, status)
	ff = recipes.NewMockRecipeFileFetcher()
	ff.FetchRecipeFileFunc = fetchRecipeFileFunc
	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}

	recipe, err := i.recipeFromPath("http://recipe/URL")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := types.InstallerContext{}
	rf := recipes.NewRecipeFilterRunner(ic, status)
	ff = recipes.NewMockRecipeFileFetcher()
	ff.LoadRecipeFileFunc = loadRecipeFileFunc
	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}

	recipe, err := i.recipeFromPath("file.txt")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestInstall_DiscoveryComplete(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{}
	statusReporter := execution.NewMockStatusReporter()
	statusReporters = []execution.StatusSubscriber{statusReporter}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}

	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporter.DiscoveryCompleteCallCount)
}

func TestInstall_FailsOnInvalidOs(t *testing.T) {
	ic := types.InstallerContext{}
	discover := discovery.NewMockDiscoverer()
	discover.SetOs("darwin")
	mv = discovery.NewManifestValidator()
	statusReporter := execution.NewMockStatusReporter()
	statusReporters = []execution.StatusSubscriber{statusReporter}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, discover, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}

	err := i.Install()
	require.Error(t, err)
}

func TestInstall_RecipeAvailable(t *testing.T) {
	ic := types.InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			DisplayName:    types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeAvailableCallCount)
}

func TestInstall_RecipeInstalled(t *testing.T) {
	ic := types.InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			DisplayName:    types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
			LogMatch: []types.OpenInstallationLogMatch{
				{
					Name: "docker log",
					File: "/var/lib/docker/containers/*/*.log",
				},
			},
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeFailed(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{}

	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())

	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	rv := validation.NewMockRecipeValidator()
	rv.ValidateErr = errors.New("validationErr")

	i := RecipeInstaller{ic, d, l, mv, f, e, rv, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 1, rv.ValidateCallCount)
	// Infra fails fast
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_NonInfraRecipeFailed(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{}

	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())

	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	rv := validation.NewMockRecipeValidator()
	rv.ValidateErrs = []error{nil, errors.New("validationErr")}

	i := RecipeInstaller{ic, d, l, mv, f, e, rv, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 2, rv.ValidateCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_AllRecipesFailed(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{}

	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())

	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           anotherTestRecipeName,
			DisplayName:    anotherTestRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	rv := validation.NewMockRecipeValidator()
	rv.ValidateErr = errors.New("validationErr")

	i := RecipeInstaller{ic, d, l, mv, f, e, rv, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 2, rv.ValidateCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_InstallComplete(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).InstallCanceledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
}

func TestInstall_InstallCanceled(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesErr = types.ErrInterrupt

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCanceledCallCount)
}

func TestInstall_InstallCompleteError(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	rv := validation.NewMockRecipeValidator()
	rv.ValidateErr = errors.New("test error")

	i := RecipeInstaller{ic, d, l, mv, f, e, rv, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_InstallCompleteError_NoFailureWhenAnyRecipeSucceeds(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           "badRecipe",
			ValidationNRQL: "testNrql",
		},
	}

	rv := validation.NewMockRecipeValidator()
	rv.ValidateErrs = []error{
		nil,
		errors.New("testing error"),
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, rv, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_RecipeSkipped(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
		{
			Name:           testRecipeName,
			DisplayName:    "test displayName",
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkippedApm(t *testing.T) {
	ic := types.InstallerContext{
		SkipApm: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
		{
			Name:           testRecipeName,
			DisplayName:    "test displayName",
			ValidationNRQL: "testNrql",
			Keywords:       []string{"apm"},
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkippedApmAnyKeyword(t *testing.T) {
	ic := types.InstallerContext{
		SkipApm: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
		{
			Name:           testRecipeName,
			DisplayName:    "test displayName",
			ValidationNRQL: "testNrql",
			Keywords:       []string{"xy", "apm", "z"},
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_SkipIntegrations(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	log.SetLevel(log.TraceLevel)
	ic := types.InstallerContext{
		AssumeYes:        true,
		SkipIntegrations: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name: "test-recipe",
			InstallTargets: []types.OpenInstallationRecipeInstallTarget{
				{Type: types.OpenInstallationTargetTypeTypes.HOST},
			},
			ValidationNRQL: "testNrql",
		},
	}

	mp := &ux.MockPrompter{
		PromptYesNoVal:       false,
		PromptMultiSelectVal: []string{},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, mp, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_MultiSelect(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	mp := &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, mp, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeRecommended(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
			InstallTargets: []types.OpenInstallationRecipeInstallTarget{
				{
					Type: types.OpenInstallationTargetTypeTypes.HOST,
				},
			},
		},
		{
			Name:           anotherTestRecipeName,
			DisplayName:    anotherTestRecipeName,
			ValidationNRQL: "testNrql",
			InstallTargets: []types.OpenInstallationRecipeInstallTarget{
				{
					Type: types.OpenInstallationTargetTypeTypes.HOST,
				},
			},
		},
		{
			Name:           "java-java-java",
			ValidationNRQL: "testNrql",
			InstallTargets: []types.OpenInstallationRecipeInstallTarget{
				{
					Type: types.OpenInstallationTargetTypeTypes.APPLICATION,
				},
			},
		},
	}

	mp := &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{types.InfraAgentRecipeName, testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, mp, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[testRecipeName])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[types.InfraAgentRecipeName])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeRecommendedCallCount)
	require.Equal(t, 4, statusReporters[0].(*execution.MockStatusReporter).RecipeAvailableCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportRecommended["java-java-java"])
}

func TestInstall_RecipeSkipped_AssumeYes(t *testing.T) {
	ic := types.InstallerContext{
		AssumeYes: true,
	}

	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
		{
			Name:           testRecipeName,
			DisplayName:    "test displayName",
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_TargetedInstall_InstallsInfraAgent(t *testing.T) {
	ic := types.InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_FilterAllButProvided(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{
		RecipeNames: []string{testRecipeName},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           testRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           anotherTestRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_InstallsInfraAgentDependency(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{
		RecipeNames: []string{testRecipeName},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           testRecipeName,
			ValidationNRQL: "testNrql",
			Dependencies:   []string{types.InfraAgentRecipeName},
		},
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstallInfraAgent_NoInfraAgentDuplicate(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{
		RecipeNames: []string{types.InfraAgentRecipeName},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_SkipInfra(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{
		SkipInfraInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_SkipInfraDependency(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	ic := types.InstallerContext{
		SkipInfraInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           "testRecipe",
			ValidationNRQL: "testNrql",
			Dependencies:   []string{types.InfraAgentRecipeName},
		},
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_GuidReport(t *testing.T) {
	ic := types.InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewPlatformLinkGenerator())
	rf := recipes.NewRecipeFilterRunner(ic, status)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipesVal = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			DisplayName:    types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
			Dependencies:   []string{types.InfraAgentRecipeName},
		},
		{
			Name:           testRecipeName,
			DisplayName:    testRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	rv := validation.NewMockRecipeValidator()
	rv.ValidateVals = []string{
		"INFRAGUID",
		"TESTRECIPEGUID",
	}

	// Test for NEW_RELIC_CLI_VERSION
	os.Setenv("NEW_RELIC_CLI_VERSION", "testversion0.0.1")

	i := RecipeInstaller{ic, d, l, mv, f, e, rv, ff, status, p, pi, lkf, cv, rvp, rf}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, rv.ValidateCallCount)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, rv.ValidateVals[0], statusReporters[0].(*execution.MockStatusReporter).RecipeGUID[types.InfraAgentRecipeName])
	require.Equal(t, rv.ValidateVals[1], statusReporters[0].(*execution.MockStatusReporter).RecipeGUID[testRecipeName])
	require.Equal(t, status.CLIVersion, "testversion0.0.1")
	require.Equal(t, 3, len(statusReporters[0].(*execution.MockStatusReporter).Durations))
	for _, duration := range statusReporters[0].(*execution.MockStatusReporter).Durations {
		require.Less(t, int64(0), duration)
	}
}

func fetchRecipeFileFunc(recipeURL *url.URL) (*types.OpenInstallationRecipe, error) {
	return testRecipeFile, nil
}

func loadRecipeFileFunc(filename string) (*types.OpenInstallationRecipe, error) {
	return testRecipeFile, nil
}
