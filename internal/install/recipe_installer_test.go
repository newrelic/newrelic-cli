// +build unit

package install

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

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
	testRecipeName = "Test Recipe"
	testRecipeFile = &recipes.RecipeFile{
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
	s               = ux.NewMockProgressIndicator()
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}

	require.True(t, reflect.DeepEqual(ic, i.InstallerContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.FetchRecipeFileFunc = fetchRecipeFileFunc
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}

	recipe, err := i.recipeFromPath("http://recipe/URL")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.LoadRecipeFileFunc = loadRecipeFileFunc
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}

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
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, f.FetchRecipeNameCount[infraAgentRecipeName], 1)
	require.Equal(t, f.FetchRecipeNameCount[loggingRecipeName], 1)
}

func TestInstall_RecipesAvailable(t *testing.T) {
	ic := InstallerContext{}
	statusReporters = []execution.StatusSubscriber{execution.NewMockStatusReporter()}
	status = execution.NewInstallStatus(statusReporters)
	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
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
	f.FetchRecipeVals = []types.Recipe{{}, {}}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_skipping_everything(t *testing.T) {
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
	f.FetchRecipeVals = []types.Recipe{{Name: "one"}, {Name: "two"}}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       false,
		PromptMultiSelectVal: []string{},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	// The infra agent always gets installed
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_multiselect(t *testing.T) {
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
	f.FetchRecipeVals = []types.Recipe{{Name: "one"}, {Name: "two"}}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	// The infra agent always gets installed
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeRecommended(t *testing.T) {
	ic := InstallerContext{}
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
			Name:           "java-java-java",
			ValidationNRQL: "testNrql",
			InstallTargets: []types.OpenInstallationRecipeInstallTarget{
				{
					Type: types.OpenInstallationTargetTypeTypes.APPLICATION,
				},
			},
		},
	}
	f.FetchRecipeVals = []types.Recipe{{Name: "one"}, {Name: "two"}}

	v = validation.NewMockRecipeValidator()
	p = &ux.MockPrompter{
		PromptYesNoVal:       true,
		PromptMultiSelectVal: []string{testRecipeName},
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	// The infra agent is always installed
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled[testRecipeName])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportInstalled["one"])
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeRecommendedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).ReportRecommended["java-java-java"])
}

func TestInstallAdvancedMode_bounce_on_enter(t *testing.T) {
	ic := InstallerContext{}
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
		PromptYesNoVal: true,
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 2, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 1, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func TestInstall_RecipeSkipped_assumeYes(t *testing.T) {
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, status, p, s}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 0, statusReporters[0].(*execution.MockStatusReporter).RecipeSkippedCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstallingCallCount)
	require.Equal(t, 3, statusReporters[0].(*execution.MockStatusReporter).RecipeInstalledCallCount)
}

func fetchRecipeFileFunc(recipeURL *url.URL) (*recipes.RecipeFile, error) {
	return testRecipeFile, nil
}

func loadRecipeFileFunc(filename string) (*recipes.RecipeFile, error) {
	return testRecipeFile, nil
}
