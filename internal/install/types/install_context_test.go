package types

import (
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
