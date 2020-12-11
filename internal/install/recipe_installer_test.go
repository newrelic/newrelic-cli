// +build unit

package install

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstall(t *testing.T) {
	assert.True(t, true)
}

func TestNewRecipeInstaller_InstallContextFields(t *testing.T) {
	ic := installContext{
		specifyActions:    true,
		installInfraAgent: true,
		installLogging:    true,
		recipeFilenames:   []string{"testRecipeFilename"},
		recipeNames:       []string{"testRecipeName"},
	}

	d := newMockDiscoverer()
	l := newMockFileFilterer()
	f := newMockRecipeFetcher()
	e := newMockRecipeExecutor()
	v := newMockRecipeValidator()
	ff := newMockRecipeFileFetcher()

	i := newRecipeInstaller(ic, d, l, f, e, v, ff)

	require.True(t, reflect.DeepEqual(ic, i.installContext))
}
