// +build unit

package execution

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestReportRecipesAvailable_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup(nil)

	recipes := []types.Recipe{{}}

	err := r.ReportRecipesAvailable(status, recipes)
	require.NoError(t, err)
}

func TestReportRecipesAvailable_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	recipes := []types.Recipe{{}}

	err := r.ReportRecipesAvailable(status, recipes)
	require.Error(t, err)
}

func TestReportRecipeInstalled_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeInstalled_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	e := RecipeStatusEvent{}

	err := r.ReportRecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeInstalled_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeInstalled(status, e)
	require.Error(t, err)
}

func TestReportRecipeInstalled_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeInstalled(status, e)
	require.Error(t, err)
}

func TestReportRecipeFailed_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeFailed_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	e := RecipeStatusEvent{}

	err := r.ReportRecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeFailed_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeFailed(status, e)
	require.Error(t, err)
}

func TestReportRecipeFailed_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeFailed(status, e)
	require.Error(t, err)
}

func TestReportComplete_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	err := r.ReportComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportComplete_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewStatusRollup([]StatusReporter{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.ReportComplete(status)
	require.Error(t, err)
}
