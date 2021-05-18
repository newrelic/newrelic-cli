// +build unit

package install

import (
	"errors"
	"net/url"
	"os"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/diagnose"
	"github.com/newrelic/newrelic-cli/internal/install/discovery"
	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/install/validation"
)

var (
	testRecipeName        = "Test Recipe"
	anotherTestRecipeName = "Another Test Recipe"
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
	status          = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	p               = ux.NewMockPrompter()
	pi              = ux.NewMockProgressIndicator()
	lkf             = NewMockLicenseKeyFetcher()
	cv              = diagnose.NewMockConfigValidator()
)

func TestInstall(t *testing.T) {
	assert.True(t, true)
}

func TestNewRecipeInstaller_InstallerContextFields(t *testing.T) {
	ic := InstallerContext{
		RecipePaths:        []string{"testRecipePath"},
		RecipeNames:        []string{"testRecipeName"},
		SkipDiscovery:      true,
		SkipIntegrations:   true,
		SkipLoggingInstall: true,
		SkipApm:            true,
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}

	require.True(t, reflect.DeepEqual(ic, i.InstallerContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.FetchRecipeFileFunc = fetchRecipeFileFunc
	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}

	recipe, err := i.recipeFromPath("http://recipe/URL")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.LoadRecipeFileFunc = loadRecipeFileFunc
	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}

	recipe, err := i.recipeFromPath("file.txt")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestInstall_Basic(t *testing.T) {
	ic := InstallerContext{}
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{Name: types.InfraAgentRecipeName},
		{Name: types.LoggingRecipeName},
	}
	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, f.FetchRecipeNameCount[types.InfraAgentRecipeName], 1)
	require.Equal(t, f.FetchRecipeNameCount[types.LoggingRecipeName], 1)
}

func TestInstall_DiscoveryComplete(t *testing.T) {
	ic := InstallerContext{}
	statusReporter := execution.NewMockStatusReporter()
	statusReporters = []execution.StatusSubscriber{statusReporter}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}

	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporter.DiscoveryCompleteCallCount)
}

func TestInstall_FailsOnInvalidOs(t *testing.T) {
	ic := InstallerContext{}
	discover := discovery.NewMockDiscoverer()
	discover.SetOs("darwin")
	mv = discovery.NewManifestValidator()
	statusReporter := execution.NewMockStatusReporter()
	statusReporters = []execution.StatusSubscriber{statusReporter}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, discover, l, mv, f, e, v, ff, status, p, pi, lkf, cv}

	err := i.Install()
	require.Error(t, err)
}

func TestInstall_RecipesAvailable(t *testing.T) {
	ic := InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
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
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipesAvailableCallCount)
}

func TestInstall_RecipeInstalled(t *testing.T) {
	ic := InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f2 := recipes.NewMockRecipeFetcher()
	f2.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f2.FetchRecipeVals = []types.OpenInstallationRecipe{
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
	}

	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f2, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeFailed(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
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
	}

	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateErr = errors.New("validationErr")

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 1, v.ValidateCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
}

func TestInstall_InstallComplete(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).InstallCanceledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
}

func TestInstall_InstallCanceled(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipeErr = types.ErrInterrupt

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCanceledCallCount)
}

func TestInstall_InstallCompleteError(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	p = &ux.MockPrompter{
		PromptYesNoVal: true,
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateErr = errors.New("test error")

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

// TestInstall_InstallCompleteError_guidedRecipeFail ensures that when we
// perform a guided install, that we only exit with an error if the infra-agent
// recipe fails.  If a recommended recipe fails, a log is emittd, but we do not
// want the error being returned.
// https://newrelic.atlassian.net/browse/VIRTUOSO-454
func TestInstall_InstallCompleteError_guidedRecipeFail(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           "badRecipe",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateErrs = []error{
		nil,
		errors.New("testing error"),
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_RecipeSkipped(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    "test displayName",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkippedApm(t *testing.T) {
	ic := InstallerContext{
		SkipApm: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    "test displayName",
		ValidationNRQL: "testNrql",
		Keywords:       []string{"apm"},
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkippedApmAnyKeyword(t *testing.T) {
	ic := InstallerContext{
		SkipApm: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    "test displayName",
		ValidationNRQL: "testNrql",
		Keywords:       []string{"xy", "apm", "z"},
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_SkipAll(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
		SkipIntegrations:   true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           "test-recipe",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       false,
		PromptMultiSelectVal: []string{},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_MultiSelect(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeRecommended(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{
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
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[testRecipeName])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[types.InfraAgentRecipeName])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeRecommendedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeAvailableCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipesAvailableCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportRecommended["java-java-java"])
}

func TestInstall_RecipeSkipped_AssumeYes(t *testing.T) {
	ic := InstallerContext{
		AssumeYes: true,
	}

	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    "test displayName",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:        types.InfraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        types.LoggingRecipeName,
			DisplayName: "Logging Recipe",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal: true,
	}

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_TargetedInstall_InstallsInfraAgent(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ic := InstallerContext{
		RecipeNames: []string{types.InfraAgentRecipeName},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_InstallsInfraAgentDependency(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ic := InstallerContext{
		RecipeNames: []string{"testRecipe"},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
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

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstallInfraAgent_NoInfraAgentDuplicate(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ic := InstallerContext{
		RecipeNames: []string{types.InfraAgentRecipeName},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_SkipInfra(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ic := InstallerContext{
		RecipeNames: []string{types.InfraAgentRecipeName},
		SkipInfra:   true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_TargetedInstall_SkipInfraDependency(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ic := InstallerContext{
		RecipeNames: []string{"testRecipe"},
		SkipInfra:   true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
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

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func TestInstall_GuidReport(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters, execution.NewConcreteSuccessLinkGenerator())
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.OpenInstallationRecipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.OpenInstallationRecipe{
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
	}

	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateVals = []string{
		"INFRAGUID",
		"TESTRECIPEGUID",
	}

	// Test for NEW_RELIC_CLI_VERSION
	os.Setenv("NEW_RELIC_CLI_VERSION", "testversion0.0.1")

	i := RecipeInstaller{ic, d, l, mv, f, e, v, ff, status, p, pi, lkf, cv}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, v.ValidateCallCount)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, v.ValidateVals[0], statusReporters[0].(*execution.MockStatusReporter).RecipeGUID[types.InfraAgentRecipeName])
	require.Equal(t, v.ValidateVals[1], statusReporters[0].(*execution.MockStatusReporter).RecipeGUID[testRecipeName])
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
