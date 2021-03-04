// +build unit

package install

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	testRecipeFile        = &recipes.RecipeFile{
		Name: testRecipeName,
	}

	d               = discovery.NewMockDiscoverer()
	l               = discovery.NewMockFileFilterer()
	f               = recipes.NewMockRecipeFetcher()
	e               = execution.NewMockRecipeExecutor()
	v               = validation.NewMockRecipeValidator()
	ff              = recipes.NewMockRecipeFileFetcher()
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status          = execution.NewInstallStatus(statusReporters)
	p               = ux.NewMockPrompter()
	pi              = ux.NewMockProgressIndicator()
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
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}

	require.True(t, reflect.DeepEqual(ic, i.InstallerContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.FetchRecipeFileFunc = fetchRecipeFileFunc
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}

	recipe, err := i.recipeFromPath("http://recipe/URL")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.LoadRecipeFileFunc = loadRecipeFileFunc
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}

	recipe, err := i.recipeFromPath("file.txt")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestInstall_Basic(t *testing.T) {
	ic := InstallerContext{}
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecipeVals = []types.Recipe{
		{Name: infraAgentRecipeName},
		{Name: loggingRecipeName},
	}
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, f.FetchRecipeNameCount[infraAgentRecipeName], 1)
	require.Equal(t, f.FetchRecipeNameCount[loggingRecipeName], 1)
}

func TestInstall_DiscoveryComplete(t *testing.T) {
	ic := InstallerContext{}
	statusReporter := execution.NewMockStatusReporter()
	statusReporters = []execution.StatusSubscriber{statusReporter}
	status = execution.NewInstallStatus(statusReporters)
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			DisplayName:    infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}

	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporter.DiscoveryCompleteCallCount)
}

func TestInstall_RecipesAvailable(t *testing.T) {
	ic := InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			DisplayName:    infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			DisplayName:    loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipesAvailableCallCount)
}

func TestInstall_RecipeInstalled(t *testing.T) {
	ic := InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			DisplayName:    infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			DisplayName:    loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeFailed(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			DisplayName:    infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			DisplayName:    loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateErr = errors.New("validationErr")

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
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
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
}

func TestInstall_InstallCompleteError(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	p = &ux.MockPrompter{
		PromptYesNoVal: true,
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateErr = errors.New("test error")

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.Error(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeFailedCallCount)
}

func TestInstall_RecipeSkipped(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           testRecipeName,
		DisplayName:    "test displayName",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:        infraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        loggingRecipeName,
			DisplayName: "Logging Recipe",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectAll: true,
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
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
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           "test-recipe",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       false,
		PromptMultiSelectVal: []string{},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
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
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           testRecipeName,
		DisplayName:    testRecipeName,
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
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
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{
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
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[testRecipeName])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[infraAgentRecipeName])
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
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		Name:           testRecipeName,
		DisplayName:    "test displayName",
		ValidationNRQL: "testNrql",
	}}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:        infraAgentRecipeName,
			DisplayName: "Infra Recipe",
		},
		{
			Name:        loggingRecipeName,
			DisplayName: "Logging Recipe",
		},
	}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal: true,
		// PromptMultiSelectAll: true,
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_TargetedInstall_InstallsInfraAgent(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	ic := InstallerContext{
		RecipeNames: []string{"testRecipe"},
	}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           "testRecipe",
			ValidationNRQL: "testNrql",
		},
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, pi}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).InstallCompleteCallCount)
}

func fetchRecipeFileFunc(recipeURL *url.URL) (*recipes.RecipeFile, error) {
	return testRecipeFile, nil
}

func loadRecipeFileFunc(filename string) (*recipes.RecipeFile, error) {
	return testRecipeFile, nil
}
