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
	"github.com/newrelic/newrelic-cli/internal/install/validation"
)

var (
	testRecipeName = "Test Recipe"
	testRecipeFile = &recipes.RecipeFile{
		Name: testRecipeName,
	}

	d  = discovery.NewMockDiscoverer()
	l  = discovery.NewMockFileFilterer()
	f  = recipes.NewMockRecipeFetcher()
	e  = execution.NewMockRecipeExecutor()
	v  = validation.NewMockRecipeValidator()
	ff = recipes.NewMockRecipeFileFetcher()
	sr = execution.NewMockStatusReporter()
)

func TestInstall(t *testing.T) {
	assert.True(t, true)
}

func TestNewRecipeInstaller_InstallerContextFields(t *testing.T) {
	ic := InstallerContext{
		RecipePaths:        []string{"testRecipePath"},
		RecipeNames:        []string{"testRecipeName"},
		SkipDiscovery:      true,
		SkipInfraInstall:   true,
		SkipIntegrations:   true,
		SkipLoggingInstall: true,
	}

	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}

	require.True(t, reflect.DeepEqual(ic, i.InstallerContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.FetchRecipeFileFunc = fetchRecipeFileFunc
	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}

	recipe, err := i.recipeFromPath("http://recipe/URL")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := InstallerContext{}
	ff = recipes.NewMockRecipeFileFetcher()
	ff.LoadRecipeFileFunc = loadRecipeFileFunc
	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}

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
	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}
	err := i.Install()
	require.NoError(t, err)
}

func TestInstall_ReportRecipesAvailable(t *testing.T) {
	ic := InstallerContext{}
	sr = execution.NewMockStatusReporter()
	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.ReportRecipesAvailableCallCount)
}

func TestInstall_ReportRecipeInstalled(t *testing.T) {
	ic := InstallerContext{}
	sr = execution.NewMockStatusReporter()
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
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

	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 3, sr.ReportRecipeInstalledCallCount)
}

func TestInstall_ReportRecipeFailed(t *testing.T) {
	ic := InstallerContext{
		SkipInfraInstall:   true,
		SkipLoggingInstall: true,
	}
	sr = execution.NewMockStatusReporter()
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{{
		ValidationNRQL: "testNrql",
	}}

	v = validation.NewMockRecipeValidator()
	v.ValidateErr = errors.New("testError")

	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.ReportRecipeFailedCallCount)
}

func TestInstall_ReportComplete(t *testing.T) {
	ic := InstallerContext{
		SkipInfraInstall:   true,
		SkipLoggingInstall: true,
	}
	sr = execution.NewMockStatusReporter()
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{}

	v = validation.NewMockRecipeValidator()

	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.ReportCompleteCallCount)
}

func TestInstall_ReportCompleteError(t *testing.T) {
	ic := InstallerContext{
		SkipLoggingInstall: true,
	}
	sr = execution.NewMockStatusReporter()
	f = recipes.NewMockRecipeFetcher()
	f.FetchRecommendationsVal = []types.Recipe{}
	f.FetchRecipeVals = []types.Recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = validation.NewMockRecipeValidator()
	v.ValidateErr = errors.New("test error")

	i := RecipeInstaller{ic, d, l, f, e, v, ff, sr}
	err := i.Install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.ReportCompleteCallCount)
}

func fetchRecipeFileFunc(recipeURL *url.URL) (*recipes.RecipeFile, error) {
	return testRecipeFile, nil
}

func loadRecipeFileFunc(filename string) (*recipes.RecipeFile, error) {
	return testRecipeFile, nil
}
