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

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	recipes := []types.Recipe{{}}

	err := r.ReportRecipesAvailable(recipes)
	require.Error(t, err)
}

func TestReportRecipeInstalled_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeInstalled(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeInstalled_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	e := RecipeStatusEvent{}

	err := r.ReportRecipeInstalled(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeInstalled_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeInstalled(e)
	require.Error(t, err)
}

func TestReportRecipeInstalled_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeInstalled(e)
	require.Error(t, err)
}

func TestReportRecipeFailed_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeFailed(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeFailed_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	e := RecipeStatusEvent{}

	err := r.ReportRecipeFailed(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeFailed_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeFailed(e)
	require.Error(t, err)
}

func TestReportRecipeFailed_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.ReportRecipeFailed(e)
	require.Error(t, err)
}

func TestReportComplete_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	err := r.ReportComplete()
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportComplete_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.ReportComplete()
	require.Error(t, err)
}
