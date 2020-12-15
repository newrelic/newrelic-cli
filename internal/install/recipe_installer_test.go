// +build unit

package install

import (
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

	d := newMockDiscoverer()
	l := newMockFileFilterer()
	f := newMockRecipeFetcher()
	e := newMockRecipeExecutor()
	v := newMockRecipeValidator()
	ff := newMockRecipeFileFetcher()
	sr := newMockExecutionStatusReporter()

	i := newRecipeInstaller(ic, d, l, f, e, v, ff, sr)

	require.True(t, reflect.DeepEqual(ic, i.installContext))
}

func TestShouldGetRecipeFromURL(t *testing.T) {
	ic := installContext{}
	ff := newMockRecipeFileFetcher()
	ff.fetchRecipeFileFunc = fetchRecipeFileFunc
	i := newRecipeInstaller(ic, nil, nil, nil, nil, nil, ff, nil)

	recipe := i.recipeFromPathFatal("http://recipe/URL")
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func TestShouldGetRecipeFromFile(t *testing.T) {
	ic := installContext{}
	ff := newMockRecipeFileFetcher()
	ff.loadRecipeFileFunc = loadRecipeFileFunc
	i := newRecipeInstaller(ic, nil, nil, nil, nil, nil, ff, nil)

	recipe := i.recipeFromPathFatal("file.txt")
	require.NotNil(t, recipe)
	require.Equal(t, recipe.Name, testRecipeName)
}

func fetchRecipeFileFunc(recipeURL *url.URL) (*recipeFile, error) {
	return testRecipeFile, nil
}

func loadRecipeFileFunc(filename string) (*recipeFile, error) {
	return testRecipeFile, nil
}
