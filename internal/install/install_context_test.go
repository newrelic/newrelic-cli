package install

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldRunDiscovery_Default(t *testing.T) {
	ic := InstallerContext{}
	require.True(t, ic.ShouldRunDiscovery())

	ic.SkipDiscovery = true
	require.False(t, ic.ShouldRunDiscovery())
}

func TestShouldInstallInfraAgent_Default(t *testing.T) {
	ic := InstallerContext{}
	require.True(t, ic.ShouldInstallInfraAgent())

	ic.SkipInfraInstall = true
	require.False(t, ic.ShouldInstallInfraAgent())
}

func TestShouldInstallInfraAgent_RecipePathsProvided(t *testing.T) {
	ic := InstallerContext{
		RecipePaths: []string{"testPath"},
	}
	require.False(t, ic.ShouldInstallInfraAgent())
}

func TestShouldInstallLogging_Default(t *testing.T) {
	ic := InstallerContext{}
	require.True(t, ic.ShouldInstallLogging())

	ic.SkipLoggingInstall = true
	require.False(t, ic.ShouldInstallLogging())
}

func TestShouldInstallLogging_RecipesProvided(t *testing.T) {
	ic := InstallerContext{
		RecipePaths: []string{"testPath"},
	}
	require.False(t, ic.ShouldInstallLogging())
}

func TestShouldInstallIntegrations_Default(t *testing.T) {
	ic := InstallerContext{}
	require.True(t, ic.ShouldInstallIntegrations())

	ic.SkipIntegrations = true
	require.False(t, ic.ShouldInstallIntegrations())
}

func TestShouldInstallLogging_RecipePathsProvided(t *testing.T) {
	ic := InstallerContext{
		RecipePaths: []string{"testPath"},
	}
	require.True(t, ic.ShouldInstallIntegrations())
}

func TestRecipeNamesProvided(t *testing.T) {
	ic := InstallerContext{}

	require.False(t, ic.RecipeNamesProvided())

	ic.RecipeNames = []string{"testName"}
	require.True(t, ic.RecipeNamesProvided())
}

func TestRecipePathsProvided(t *testing.T) {
	ic := InstallerContext{}
	require.False(t, ic.RecipePathsProvided())

	ic.RecipePaths = []string{"testPath"}
	require.True(t, ic.RecipePathsProvided())
}

func TestRecipesProvided(t *testing.T) {
	ic := InstallerContext{}
	require.False(t, ic.RecipesProvided())

	ic.RecipePaths = []string{"testPath"}
	require.True(t, ic.RecipesProvided())

	ic.RecipePaths = []string{}
	ic.RecipeNames = []string{"testName"}
	require.True(t, ic.RecipesProvided())
}
