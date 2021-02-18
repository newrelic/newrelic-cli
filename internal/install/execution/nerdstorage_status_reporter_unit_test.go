// +build unit

package execution

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecipesAvailable_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus(nil)

	recipes := []types.Recipe{{}}

	err := r.RecipesAvailable(status, recipes)
	require.NoError(t, err)
}

func TestRecipesAvailable_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	recipes := []types.Recipe{{}}

	err := r.RecipesAvailable(status, recipes)
	require.Error(t, err)
}

func TestRecipeInstalled_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestRecipeInstalled_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestRecipeInstalled_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.RecipeInstalled(status, e)
	require.Error(t, err)
}

func TestRecipeInstalled_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.RecipeInstalled(status, e)
	require.Error(t, err)
}

func TestRecipeFailed_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
}

func TestRecipeFailed_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	e := RecipeStatusEvent{}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestRecipeFailed_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.RecipeFailed(status, e)
	require.Error(t, err)
}

func TestRecipeFailed_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	e := RecipeStatusEvent{
		EntityGUID: "testGuid",
	}

	err := r.RecipeFailed(status, e)
	require.Error(t, err)
}

func TestInstallComplete_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
}

func TestInstallComplete_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	status := NewInstallStatus([]StatusSubscriber{r})

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.InstallComplete(status)
	require.Error(t, err)
}
