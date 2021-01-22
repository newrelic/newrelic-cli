// +build unit

package execution

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestNewStatus(t *testing.T) {
	s := NewStatusRollup([]StatusReporter{})
	require.NotEmpty(t, s.Timestamp)
	require.NotEmpty(t, s.DocumentID)
}

func TestStatusWithAvailableRecipes_Basic(t *testing.T) {
	s := NewStatusRollup([]StatusReporter{})
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
	s := NewStatusRollup([]StatusReporter{})
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
	s := NewStatusRollup([]StatusReporter{})
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

func TestStatusWithRecipeEvent_EntityGUID(t *testing.T) {
	s := NewStatusRollup([]StatusReporter{})
	r := types.Recipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.Timestamp = 0
	s.withRecipeEvent(e, StatusTypes.INSTALLED)

	require.NotEmpty(t, s.EntityGUIDs)
	require.Equal(t, 1, len(s.EntityGUIDs))
	require.Equal(t, "testGUID", s.EntityGUIDs[0])
}

func TestStatusWithRecipeEvent_EntityGUIDExists(t *testing.T) {
	s := NewStatusRollup([]StatusReporter{})
	s.withEntityGUID("testGUID")
	r := types.Recipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.Timestamp = 0
	s.withRecipeEvent(e, StatusTypes.INSTALLED)

	require.NotEmpty(t, s.EntityGUIDs)
	require.Equal(t, 1, len(s.EntityGUIDs))
	require.Equal(t, "testGUID", s.EntityGUIDs[0])
}

func TestStatusRollup_statusUpdateMethods(t *testing.T) {
	s := NewStatusRollup([]StatusReporter{})
	r := types.Recipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.ReportRecipesAvailable([]types.Recipe{r})

	result := s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, StatusTypes.AVAILABLE)

	s.ReportRecipeInstalling(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, StatusTypes.INSTALLING)

	s.ReportRecipeInstalled(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, StatusTypes.INSTALLED)

	s.ReportRecipeFailed(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, StatusTypes.FAILED)
	require.True(t, s.hasFailed())

	s.ReportRecipeSkipped(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, StatusTypes.SKIPPED)
	require.False(t, s.hasFailed())

	s.ReportComplete(RecipeStatusEvent{EntityGUID: "123"})
	require.Equal(t, s.EntityGUIDs[0], "testGUID")
	require.Equal(t, s.EntityGUIDs[1], "123")
	require.True(t, s.Complete)
	require.NotNil(t, s.Timestamp)
}
