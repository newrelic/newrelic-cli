package install

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldRunDiscovery_Default(t *testing.T) {
	ic := installContext{}
	require.True(t, ic.ShouldRunDiscovery())

	ic.skipDiscovery = true
	require.False(t, ic.ShouldRunDiscovery())
}

func TestShouldInstallInfraAgent_Default(t *testing.T) {
	ic := installContext{}
	require.True(t, ic.ShouldInstallInfraAgent())

	ic.skipInfraInstall = true
	require.False(t, ic.ShouldInstallInfraAgent())
}

func TestShouldInstallInfraAgent_RecipePathsProvided(t *testing.T) {
	ic := installContext{
		recipePaths: []string{"testPath"},
	}
	require.False(t, ic.ShouldInstallInfraAgent())
}

func TestShouldInstallLogging_Default(t *testing.T) {
	ic := installContext{}
	require.True(t, ic.ShouldInstallLogging())

	ic.skipLoggingInstall = true
	require.False(t, ic.ShouldInstallLogging())
}

func TestShouldInstallLogging_RecipesProvided(t *testing.T) {
	ic := installContext{
		recipePaths: []string{"testPath"},
	}
	require.False(t, ic.ShouldInstallLogging())
}

func TestShouldInstallIntegrations_Default(t *testing.T) {
	ic := installContext{}
	require.True(t, ic.ShouldInstallIntegrations())

	ic.skipIntegrations = true
	require.False(t, ic.ShouldInstallIntegrations())
}

func TestShouldInstallLogging_RecipePathsProvided(t *testing.T) {
	ic := installContext{
		recipePaths: []string{"testPath"},
	}
	require.True(t, ic.ShouldInstallIntegrations())
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

	ic.recipePaths = []string{"testPath"}
	require.True(t, ic.RecipePathsProvided())
}

func TestRecipesProvided(t *testing.T) {
	ic := installContext{}
	require.False(t, ic.RecipesProvided())

	ic.recipePaths = []string{"testPath"}
	require.True(t, ic.RecipesProvided())

	ic.recipePaths = []string{}
	ic.recipeNames = []string{"testName"}
	require.True(t, ic.RecipesProvided())
}
