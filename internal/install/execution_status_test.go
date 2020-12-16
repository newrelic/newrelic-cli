//+ build unit
package install

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewExecutionStatus(t *testing.T) {
	s := newExecutionStatus()
	require.NotEmpty(t, s.Timestamp)
	require.NotEmpty(t, s.DocumentID)
}

func TestExecutionStatusWithAvailableRecipes_Basic(t *testing.T) {
	s := newExecutionStatus()
	r := []recipe{{
		Name: "testRecipe1",
	}, {
		Name: "testRecipe2",
	}}

	s.withAvailableRecipes(r)

	require.NotEmpty(t, s.Recipes)
	require.Equal(t, len(r), len(s.Recipes))
	for _, recipeStatus := range s.Recipes {
		require.Equal(t, executionStatusRecipeStatusTypes.AVAILABLE, recipeStatus.Status)
	}
}

func TestExecutionStatusWithRecipeEvent_Basic(t *testing.T) {
	s := newExecutionStatus()
	r := recipe{Name: "testRecipe"}
	e := recipeStatusEvent{recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, executionStatusRecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Recipes)
	require.Equal(t, 1, len(s.Recipes))
	require.Equal(t, executionStatusRecipeStatusTypes.INSTALLED, s.Recipes[0].Status)
	require.NotEmpty(t, s.Timestamp)
}

func TestExecutionStatusWithRecipeEvent_RecipeExists(t *testing.T) {
	s := newExecutionStatus()
	r := recipe{Name: "testRecipe"}
	e := recipeStatusEvent{recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, executionStatusRecipeStatusTypes.AVAILABLE)

	require.NotEmpty(t, s.Recipes)
	require.Equal(t, 1, len(s.Recipes))
	require.Equal(t, executionStatusRecipeStatusTypes.AVAILABLE, s.Recipes[0].Status)
	require.NotEmpty(t, s.Timestamp)

	s.Timestamp = 0
	s.withRecipeEvent(e, executionStatusRecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Recipes)
	require.Equal(t, 1, len(s.Recipes))
	require.Equal(t, executionStatusRecipeStatusTypes.INSTALLED, s.Recipes[0].Status)
	require.NotEmpty(t, s.Timestamp)
}
