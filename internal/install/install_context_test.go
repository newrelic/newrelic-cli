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

func TestShouldInstallInfraAgent_RecipePathsProvided(t *testing.T) {
	ic := installContext{
		recipePaths: []string{"testFilename"},
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

func TestShouldInstallLogging_RecipePathsProvided(t *testing.T) {
	ic := installContext{
		recipePaths: []string{"testFilename"},
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

func TestRecipePathsProvided(t *testing.T) {
	ic := installContext{}
	require.False(t, ic.RecipePathsProvided())

	ic.recipePaths = []string{"testFilename"}
	require.True(t, ic.RecipePathsProvided())
}
