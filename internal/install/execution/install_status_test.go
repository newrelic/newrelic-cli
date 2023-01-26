//go:build unit
// +build unit

package execution

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestNewInstallStatus(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	require.NotEmpty(t, s.Timestamp)
	require.NotEmpty(t, s.DocumentID)
}

func TestStatusWithRecipeEvent_Basic(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
	require.True(t, s.RecipeHasStatus(r.Name, RecipeStatusTypes.INSTALLED))
}

func TestStatusWithRecipeEvent_ErrorMessages(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{
		Recipe: r,
		Msg:    "thing failed for a reason",
	}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.FAILED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.FAILED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
	require.Equal(t, s.Statuses[0].Error, s.Error)
	require.Equal(t, e.Msg, s.Error.Message)
}

func TestExecutionStatusWithRecipeEvent_RecipeExists(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.AVAILABLE)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.AVAILABLE, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
	require.True(t, s.RecipeHasStatus(r.Name, RecipeStatusTypes.AVAILABLE))

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
	require.True(t, s.RecipeHasStatus(r.Name, RecipeStatusTypes.INSTALLED))
}

func TestStatusWithRecipeEvent_EntityGUID(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.EntityGUIDs)
	require.Equal(t, 1, len(s.EntityGUIDs))
	require.Equal(t, "testGUID", s.EntityGUIDs[0])
}

func TestStatusWithRecipeEvent_EntityGUIDExists(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	s.withEntityGUID("testGUID")
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.EntityGUIDs)
	require.Equal(t, 1, len(s.EntityGUIDs))
	require.Equal(t, "testGUID", s.EntityGUIDs[0])
}

func TestInstallStatus_statusUpdateMethods(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.InstallStarted()
	require.False(t, s.Complete)
	require.NotNil(t, s.Timestamp)

	s.RecipeAvailable(NewRecipeStatusEvent(&r))

	result := s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, RecipeStatusTypes.AVAILABLE)

	s.RecipeInstalling(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, RecipeStatusTypes.INSTALLING)

	s.RecipeInstalled(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, RecipeStatusTypes.INSTALLED)

	s.RecipeFailed(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, RecipeStatusTypes.FAILED)
	require.True(t, s.hasAnyRecipeStatus(RecipeStatusTypes.FAILED))

	s.RecipeSkipped(e)
	result = s.getStatus(r)
	require.NotNil(t, result)
	require.Equal(t, result.Status, RecipeStatusTypes.SKIPPED)
	require.False(t, s.hasAnyRecipeStatus(RecipeStatusTypes.FAILED))

	s.InstallComplete(nil)
	require.True(t, s.Complete)
	require.NotNil(t, s.Timestamp)
}

func TestInstallStatus_shouldNotFailAvailableOnComplete(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}

	s.RecipeAvailable(NewRecipeStatusEvent(&r))

	s.InstallComplete(nil)
	require.Equal(t, RecipeStatusTypes.AVAILABLE, s.Statuses[0].Status)
}

func TestInstallStatus_shouldFailAvailableOnCancel(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}

	s.RecipeAvailable(NewRecipeStatusEvent(&r))

	s.InstallCanceled()
	require.Equal(t, RecipeStatusTypes.CANCELED, s.Statuses[0].Status)
}

func TestInstallStatus_multipleRecipeStatuses(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	recipeInstalled := types.OpenInstallationRecipe{Name: "installed"}
	installedRecipeEvent := RecipeStatusEvent{Recipe: recipeInstalled, EntityGUID: "installedGUID"}

	recipeSkipped := types.OpenInstallationRecipe{Name: "skipped"}
	skippedRecipeEvent := RecipeStatusEvent{Recipe: recipeSkipped, EntityGUID: "skippedGUID"}

	recipeErrored := types.OpenInstallationRecipe{Name: "errored"}
	erroredRecipeEvent := RecipeStatusEvent{Recipe: recipeErrored, EntityGUID: "erroredGUID"}

	recipeCanceled := types.OpenInstallationRecipe{Name: "installing"}
	canceledRecipeEvent := RecipeStatusEvent{Recipe: recipeCanceled, EntityGUID: "erroredGUID"}

	s.RecipeAvailable(NewRecipeStatusEvent(&recipeInstalled))
	s.RecipeInstalling(canceledRecipeEvent)
	s.RecipeInstalled(installedRecipeEvent)
	s.RecipeSkipped(skippedRecipeEvent)
	s.RecipeFailed(erroredRecipeEvent)

	s.InstallCanceled()

	require.True(t, s.HasInstalledRecipes)
	require.True(t, s.HasSkippedRecipes)
	require.True(t, s.HasCanceledRecipes)
	require.True(t, s.HasFailedRecipes)
}

func TestStatus_HTTPSProxy(t *testing.T) {
	os.Setenv("HTTPS_PROXY", "localhost:8888")
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)

	require.Equal(t, "localhost:8888", s.HTTPSProxy)
}

func TestSetTargetInstallShouldSet(t *testing.T) {

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	recipeNames := []string{"infra", "logging"}
	status.SetTargetedInstall(recipeNames)

	require.True(t, status.targetedInstall)
	require.Equal(t, status.targetedInstallNames, recipeNames)
}

func TestSetTargetInstallShouldNotSet(t *testing.T) {

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)
	recipeNames := []string{}
	status.SetTargetedInstall(recipeNames)

	require.False(t, status.targetedInstall)
	require.Equal(t, len(recipeNames), len(status.targetedInstallNames))
}
