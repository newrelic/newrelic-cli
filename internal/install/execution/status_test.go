//+ build unit
package execution

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestNewStatus(t *testing.T) {
	s := NewStatusRollup()
	require.NotEmpty(t, s.Timestamp)
	require.NotEmpty(t, s.DocumentID)
}

func TestStatusWithAvailableRecipes_Basic(t *testing.T) {
	s := NewStatusRollup()
	r := []types.Recipe{{
		Name: "testRecipe1",
	}, {
		Name: "testRecipe2",
	}}

	s.withAvailableRecipes(r)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, len(r), len(s.Statuses))
	for _, recipeStatus := range s.Statuses {
		require.Equal(t, StatusTypes.AVAILABLE, recipeStatus.Status)
	}
}

func TestStatusWithRecipeEvent_Basic(t *testing.T) {
	s := NewStatusRollup()
	r := types.Recipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, StatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, StatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
}

func TestExecutionStatusWithRecipeEvent_RecipeExists(t *testing.T) {
	s := NewStatusRollup()
	r := types.Recipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, StatusTypes.AVAILABLE)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, StatusTypes.AVAILABLE, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)

	s.Timestamp = 0
	s.withRecipeEvent(e, StatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, StatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
}
