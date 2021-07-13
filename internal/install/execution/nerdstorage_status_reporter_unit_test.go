// +build unit

package execution

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestRecipeAvailable_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	err := r.RecipeAvailable(status, types.OpenInstallationRecipe{})
	require.NoError(t, err)
}

func TestRecipeAvailable_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.RecipeAvailable(status, types.OpenInstallationRecipe{})
	require.Error(t, err)
}

func TestRecipeInstalled_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestRecipeInstalled_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestRecipeInstalled_MultipleEntityGUIDs(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	status.withEntityGUID("testGuid2")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 2, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestRecipeInstalled_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.RecipeInstalled(status, e)
	require.Error(t, err)
}

func TestRecipeInstalled_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	err := r.RecipeInstalled(status, e)
	require.Error(t, err)
}

func TestRecipeFailed_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestRecipeFailed_UserScopeOnly(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	e := RecipeStatusEvent{}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestRecipeFailed_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.RecipeFailed(status, e)
	require.Error(t, err)
}

func TestRecipeFailed_EntityScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	c.WriteDocumentWithEntityScopeErr = errors.New("error")

	err := r.RecipeFailed(status, e)
	require.Error(t, err)
}

func TestInstallComplete_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestInstallComplete_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.InstallComplete(status)
	require.Error(t, err)
}

func TestInstallCanceled_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	err := r.InstallCanceled(status)
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestInstallCanceled_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.InstallCanceled(status)
	require.Error(t, err)
}

func TestDiscoveryComplete_Basic(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	err := r.DiscoveryComplete(status, types.DiscoveryManifest{})
	require.NoError(t, err)
	require.Equal(t, 1, c.writeDocumentWithUserScopeCallCount)
	require.Equal(t, 0, c.writeDocumentWithEntityScopeCallCount)
	require.Equal(t, 1, c.writeDocumentWithAccountScopeCallCount)
}

func TestDiscoveryComplete_UserScopeError(t *testing.T) {
	c := NewMockNerdStorageClient()
	r := NewNerdStorageStatusReporter(c)
	slg := NewPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)

	c.WriteDocumentWithUserScopeErr = errors.New("error")

	err := r.DiscoveryComplete(status, types.DiscoveryManifest{})
	require.Error(t, err)
}
