// +build unit

package install

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReportRecipesAvailable_Basic(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	recipes := []recipe{{}}

	err := r.reportRecipesAvailable(recipes)
	require.NoError(t, err)
}

func TestReportRecipesAvailable_UserScopeError(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	c.userScopeError = errors.New("error")

	recipes := []recipe{{}}

	err := r.reportRecipesAvailable(recipes)
	require.Error(t, err)
}

func TestReportRecipeInstalled_Basic(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	e := recipeStatusEvent{
		entityGUID: "testGuid",
	}

	err := r.reportRecipeInstalled(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeInstalled_UserScopeOnly(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	e := recipeStatusEvent{}

	err := r.reportRecipeInstalled(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeInstalled_UserScopeError(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	c.userScopeError = errors.New("error")

	e := recipeStatusEvent{
		entityGUID: "testGuid",
	}

	err := r.reportRecipeInstalled(e)
	require.Error(t, err)
}

func TestReportRecipeInstalled_EntityScopeError(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	c.entityScopeError = errors.New("error")

	e := recipeStatusEvent{
		entityGUID: "testGuid",
	}

	err := r.reportRecipeInstalled(e)
	require.Error(t, err)
}

func TestReportRecipeFailed_Basic(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	e := recipeStatusEvent{
		entityGUID: "testGuid",
	}

	err := r.reportRecipeFailed(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeFailed_UserScopeOnly(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	e := recipeStatusEvent{}

	err := r.reportRecipeFailed(e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestReportRecipeFailed_UserScopeError(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	c.userScopeError = errors.New("error")

	e := recipeStatusEvent{
		entityGUID: "testGuid",
	}

	err := r.reportRecipeFailed(e)
	require.Error(t, err)
}

func TestReportRecipeFailed_EntityScopeError(t *testing.T) {
	c := newMockNerdstorageClient()
	r := newNerdStorageExecutionStatusReporter(c)

	c.entityScopeError = errors.New("error")

	e := recipeStatusEvent{
		entityGUID: "testGuid",
	}

	err := r.reportRecipeFailed(e)
	require.Error(t, err)
}
