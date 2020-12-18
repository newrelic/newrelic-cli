// +build unit

package install

import (
	"errors"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testRecipeName = "Test Recipe"
	testRecipeFile = &recipeFile{
		Name: testRecipeName,
	}

	d  = newMockDiscoverer()
	l  = newMockFileFilterer()
	f  = newMockRecipeFetcher()
	e  = newMockRecipeExecutor()
	v  = newMockRecipeValidator()
	ff = newMockRecipeFileFetcher()
	sr = newMockExecutionStatusReporter()
)

func TestInstall(t *testing.T) {
	assert.True(t, true)
}

func TestNewRecipeInstaller_InstallContextFields(t *testing.T) {
	ic := installContext{
		recipePaths:        []string{"testRecipePath"},
		recipeNames:        []string{"testRecipeName"},
		skipDiscovery:      true,
		skipInfraInstall:   true,
		skipIntegrations:   true,
		skipLoggingInstall: true,
	}

	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)

	require.True(t, reflect.DeepEqual(ic, i.installContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := installContext{}
	ff = newMockRecipeFileFetcher()
	ff.fetchRecipeFileFunc = fetchRecipeFileFunc
	i := newRecipeInstaller(ic, nil, nil, nil, nil, nil, ff, nil)

	recipe, err := i.recipeFromPath("http://recipe/URL")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := installContext{}
	ff = newMockRecipeFileFetcher()
	ff.loadRecipeFileFunc = loadRecipeFileFunc
	i := newRecipeInstaller(ic, nil, nil, nil, nil, nil, ff, nil)

	recipe, err := i.recipeFromPath("file.txt")
	require.NoError(t, err)
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestInstall_Basic(t *testing.T) {
	ic := installContext{}
	f = newMockRecipeFetcher()
	f.fetchRecipeVals = []recipe{
		{Name: infraAgentRecipeName},
		{Name: loggingRecipeName},
	}
	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)
	err := i.install()
	require.NoError(t, err)
}

func TestInstall_ReportRecipesAvailable(t *testing.T) {
	ic := installContext{}
	sr = newMockExecutionStatusReporter()
	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)
	err := i.install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.reportRecipesAvailableCallCount)
}

func TestInstall_ReportRecipeInstalled(t *testing.T) {
	ic := installContext{}
	sr = newMockExecutionStatusReporter()
	f = newMockRecipeFetcher()
	f.fetchRecommendationsVal = []recipe{{
		ValidationNRQL: "testNrql",
	}}
	f.fetchRecipeVals = []recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           loggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}
	v = newMockRecipeValidator()

	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)
	err := i.install()
	require.NoError(t, err)
	require.Equal(t, 3, sr.reportRecipeInstalledCallCount)
}

func TestInstall_ReportRecipeFailed(t *testing.T) {
	ic := installContext{
		skipInfraInstall:   true,
		skipLoggingInstall: true,
	}
	sr = newMockExecutionStatusReporter()
	f = newMockRecipeFetcher()
	f.fetchRecommendationsVal = []recipe{{
		ValidationNRQL: "testNrql",
	}}

	v = newMockRecipeValidator()
	v.validateErr = errors.New("testError")

	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)
	err := i.install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.reportRecipeFailedCallCount)
}

func TestInstall_ReportComplete(t *testing.T) {
	ic := installContext{
		skipInfraInstall:   true,
		skipLoggingInstall: true,
	}
	sr = newMockExecutionStatusReporter()
	f = newMockRecipeFetcher()
	f.fetchRecommendationsVal = []recipe{}

	v = newMockRecipeValidator()

	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)
	err := i.install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.reportCompleteCallCount)
}

func TestInstall_ReportCompleteError(t *testing.T) {
	ic := installContext{
		skipLoggingInstall: true,
	}
	sr = newMockExecutionStatusReporter()
	f = newMockRecipeFetcher()
	f.fetchRecommendationsVal = []recipe{}
	f.fetchRecipeVals = []recipe{
		{
			Name:           infraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	v = newMockRecipeValidator()
	v.validateErr = errors.New("test error")

	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)
	err := i.install()
	require.NoError(t, err)
	require.Equal(t, 1, sr.reportCompleteCallCount)
}

func fetchRecipeFileFunc(recipeURL *url.URL) (*recipeFile, error) {
	return testRecipeFile, nil
}

func loadRecipeFileFunc(filename string) (*recipeFile, error) {
	return testRecipeFile, nil
}
