package install

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldInstallInfraAgent_Default(t *testing.T) {
	ic := installContext{}
	require.True(t, ic.ShouldInstallInfraAgent())
}

func TestShouldInstallInfraAgent_SpecifyActions(t *testing.T) {
	ic := installContext{
		specifyActions: true,
	}
	require.False(t, ic.ShouldInstallInfraAgent())

	ic.installInfraAgent = true
	require.True(t, ic.ShouldInstallInfraAgent())
}

func TestShouldInstallInfraAgent_RecipeFilenamesProvided(t *testing.T) {
	ic := installContext{
		recipeFilenames: []string{"testFilename"},
	}
	require.False(t, ic.ShouldInstallInfraAgent())

	ic.installInfraAgent = true
	require.False(t, ic.ShouldInstallInfraAgent())
}
func TestShouldInstallLogging_Default(t *testing.T) {
	ic := installContext{}
	require.True(t, ic.ShouldInstallLogging())
}

func TestShouldInstallLogging_SpecifyActions(t *testing.T) {
	ic := installContext{
		specifyActions: true,
	}
	require.False(t, ic.ShouldInstallLogging())

	ic.installLogging = true
	require.True(t, ic.ShouldInstallLogging())
}

func TestShouldInstallLogging_RecipeFilenamesProvided(t *testing.T) {
	ic := installContext{
		recipeFilenames: []string{"testFilename"},
	}
	require.False(t, ic.ShouldInstallLogging())

	ic.installInfraAgent = true
	require.False(t, ic.ShouldInstallLogging())
}

func TestRecipeNamesProvided(t *testing.T) {
	ic := installContext{}
	require.False(t, ic.RecipeNamesProvided())

	ic.recipeNames = []string{"testName"}
	require.True(t, ic.RecipeNamesProvided())
}

func TestRecipeFilenamesProvided(t *testing.T) {
	ic := installContext{}
	require.False(t, ic.RecipeFilenamesProvided())

	ic.recipeFilenames = []string{"testFilename"}
	require.True(t, ic.RecipeFilenamesProvided())
}
