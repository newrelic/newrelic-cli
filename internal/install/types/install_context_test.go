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
	ic := InstallerContext{}
	args := []string{"tag1:test", "tag2:test"}
	ic.SetTags(args)
	args = append(args, BuiltinTags)

	require.Equal(t, args, ic.Tags)
	require.Equal(t, "tag1:test,tag2:test,"+BuiltinTags, os.Getenv(EnvInstallCustomAttributes))
}
