package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecipeVars_ToSlice(t *testing.T) {
	r := RecipeVars{
		"testKey":        "testValue",
		"anotherTestKey": "anotherTestValue",
	}

	require.ElementsMatch(t, []string{"testKey=testValue", "anotherTestKey=anotherTestValue"}, r.ToSlice())
}
