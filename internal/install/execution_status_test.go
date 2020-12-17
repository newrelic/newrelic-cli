//+ build unit
package install

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewExecutionStatus(t *testing.T) {
	s := newExecutionStatusRollup()
	require.NotEmpty(t, s.Timestamp)
	require.NotEmpty(t, s.DocumentID)
}

func TestExecutionStatusWithAvailableRecipes_Basic(t *testing.T) {
	s := newExecutionStatusRollup()
	r := []recipe{{
		Name: "testRecipe1",
	}, {
		Name: "testRecipe2",
	}}

	s.withAvailableRecipes(r)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, len(r), len(s.Statuses))
	for _, recipeStatus := range s.Statuses {
		require.Equal(t, executionStatusTypes.AVAILABLE, recipeStatus.Status)
	}
}

func TestExecutionStatusWithRecipeEvent_Basic(t *testing.T) {
	s := newExecutionStatusRollup()
	r := recipe{Name: "testRecipe"}
	e := recipeStatusEvent{recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, executionStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, executionStatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
}

func TestExecutionStatusWithRecipeEvent_RecipeExists(t *testing.T) {
	s := newExecutionStatusRollup()
	r := recipe{Name: "testRecipe"}
	e := recipeStatusEvent{recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, executionStatusTypes.AVAILABLE)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, executionStatusTypes.AVAILABLE, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)

	s.Timestamp = 0
	s.withRecipeEvent(e, executionStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, executionStatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
}
