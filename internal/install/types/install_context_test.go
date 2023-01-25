package types

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestSetTags(t *testing.T) {
	t.Setenv(EnvInstallCustomAttributes, "")
	ic := InstallerContext{}
	args := []string{"tag1:test", "tag2:test"}
	ic.SetTags(args)

	require.Equal(t, "tag1:test,tag2:test,"+BuiltinTags, os.Getenv(EnvInstallCustomAttributes))
}

func TestSetTagsShouldSkipIncorrectSyntax(t *testing.T) {
	t.Setenv(EnvInstallCustomAttributes, "")
	ic := InstallerContext{}
	args := []string{"tag1:test", "notvalidtag"}
	ic.SetTags(args)

	require.Equal(t, "tag1:test,"+BuiltinTags, os.Getenv(EnvInstallCustomAttributes))
}

func TestSetTagsShouldAddDeployedBy(t *testing.T) {
	t.Setenv(EnvInstallCustomAttributes, "")
	ic := InstallerContext{}
	ic.SetTags([]string{})

	require.Equal(t, BuiltinTags, os.Getenv(EnvInstallCustomAttributes))
}

func TestSetTagsShouldNotReplaceDeployedBy(t *testing.T) {
	t.Setenv(EnvInstallCustomAttributes, "")
	ic := InstallerContext{}
	args := []string{"nr_deployed_by:Me", "tag1:test", "tag2:test"}
	ic.SetTags(args)

	require.Equal(t, "nr_deployed_by:Me,tag1:test,tag2:test", os.Getenv(EnvInstallCustomAttributes))
}

func TestShouldGetDeployedBy(t *testing.T) {
	t.Setenv(EnvInstallCustomAttributes, "")
	ic := InstallerContext{}
	args := []string{"nr_deployed_by:SomeoneElse", "tag2:test"}
	ic.SetTags(args)

	require.Equal(t, "SomeoneElse", ic.GetDeployedBy())
}
