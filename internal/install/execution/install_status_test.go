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
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	require.NotEmpty(t, s.Timestamp)
	require.NotEmpty(t, s.DocumentID)
}

func TestStatusWithAvailableRecipes_Basic(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	r := []types.OpenInstallationRecipe{{
		Name: "testRecipe1",
	}, {
		Name: "testRecipe2",
	}}

	s.withAvailableRecipes(r)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, len(r), len(s.Statuses))
	for _, recipeStatus := range s.Statuses {
		require.Equal(t, RecipeStatusTypes.AVAILABLE, recipeStatus.Status)
	}
}

func TestStatusWithRecipeEvent_Basic(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
}

func TestStatusWithRecipeEvent_ErrorMessages(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)
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
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r}

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.AVAILABLE)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.AVAILABLE, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)

	s.Timestamp = 0
	s.withRecipeEvent(e, RecipeStatusTypes.INSTALLED)

	require.NotEmpty(t, s.Statuses)
	require.Equal(t, 1, len(s.Statuses))
	require.Equal(t, RecipeStatusTypes.INSTALLED, s.Statuses[0].Status)
	require.NotEmpty(t, s.Timestamp)
}

func TestStatusWithRecipeEvent_EntityGUID(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)
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
	s := NewInstallStatus([]StatusSubscriber{}, slg)
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
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}
	e := RecipeStatusEvent{Recipe: r, EntityGUID: "testGUID"}

	s.InstallStarted()
	require.False(t, s.Complete)
	require.NotNil(t, s.Timestamp)

	s.RecipeAvailable(r)

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

func TestInstallStatus_observabilityPackStatusUpdateMethods(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)

	eventFetchPending := ObservabilityPackStatusEvent{
		Name: "test-pack",
	}

	eventInstallPending := ObservabilityPackStatusEvent{
		ObservabilityPack: types.OpenInstallationObservabilityPack{
			Name: "test-pack",
		},
	}

	// Success event is the same for Fetch and Install
	eventSuccess := ObservabilityPackStatusEvent{
		ObservabilityPack: types.OpenInstallationObservabilityPack{
			Name: "test-pack",
		},
	}

	// Failure event is the same for Fetch and Install
	eventFailure := ObservabilityPackStatusEvent{
		ObservabilityPack: types.OpenInstallationObservabilityPack{
			Name: "test-pack",
		},
		Msg: "failure",
	}

	s.ObservabilityPackFetchPending(eventFetchPending)
	result := s.getObservabilityPackStatusByPackStatusType(ObservabilityPackStatusTypes.FetchPending)
	require.NotNil(t, result)
	require.Equal(t, result.Status, ObservabilityPackStatusTypes.FetchPending)

	s.ObservabilityPackFetchSuccess(eventSuccess)
	result = s.getObservabilityPackStatusByPackStatusType(ObservabilityPackStatusTypes.FetchSuccess)
	require.NotNil(t, result)
	require.Equal(t, result.Status, ObservabilityPackStatusTypes.FetchSuccess)

	s.ObservabilityPackFetchFailed(eventFailure)
	result = s.getObservabilityPackStatusByPackStatusType(ObservabilityPackStatusTypes.FetchFailed)
	require.NotNil(t, result)
	require.Equal(t, result.Status, ObservabilityPackStatusTypes.FetchFailed)
	require.Equal(t, result.Error.Message, "failure")

	s.ObservabilityPackInstallPending(eventInstallPending)
	result = s.getObservabilityPackStatusByPackStatusType(ObservabilityPackStatusTypes.InstallPending)
	require.NotNil(t, result)
	require.Equal(t, result.Status, ObservabilityPackStatusTypes.InstallPending)

	s.ObservabilityPackInstallSuccess(eventSuccess)
	result = s.getObservabilityPackStatusByPackStatusType(ObservabilityPackStatusTypes.InstallSuccess)
	require.NotNil(t, result)
	require.Equal(t, result.Status, ObservabilityPackStatusTypes.InstallSuccess)

	s.ObservabilityPackInstallFailed(eventFailure)
	result = s.getObservabilityPackStatusByPackStatusType(ObservabilityPackStatusTypes.InstallFailed)
	require.NotNil(t, result)
	require.Equal(t, result.Status, ObservabilityPackStatusTypes.InstallFailed)
	require.Equal(t, result.Error.Message, "failure")
}

func TestInstallStatus_shouldNotFailAvailableOnComplete(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}

	s.RecipeAvailable(r)

	s.InstallComplete(nil)
	require.Equal(t, RecipeStatusTypes.AVAILABLE, s.Statuses[0].Status)
}

func TestInstallStatus_shouldFailAvailableOnCancel(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{}, slg)
	r := types.OpenInstallationRecipe{Name: "testRecipe"}

	s.RecipeAvailable(r)

	s.InstallCanceled()
	require.Equal(t, RecipeStatusTypes.CANCELED, s.Statuses[0].Status)
}

func TestInstallStatus_multipleRecipeStatuses(t *testing.T) {
	slg := NewPlatformLinkGenerator()
	s := NewInstallStatus([]StatusSubscriber{NewMockStatusReporter()}, slg)
	recipeInstalled := types.OpenInstallationRecipe{Name: "installed"}
	installedRecipeEvent := RecipeStatusEvent{Recipe: recipeInstalled, EntityGUID: "installedGUID"}

	recipeSkipped := types.OpenInstallationRecipe{Name: "skipped"}
	skippedRecipeEvent := RecipeStatusEvent{Recipe: recipeSkipped, EntityGUID: "skippedGUID"}

	recipeErrored := types.OpenInstallationRecipe{Name: "errored"}
	erroredRecipeEvent := RecipeStatusEvent{Recipe: recipeErrored, EntityGUID: "erroredGUID"}

	recipeCanceled := types.OpenInstallationRecipe{Name: "installing"}
	canceledRecipeEvent := RecipeStatusEvent{Recipe: recipeCanceled, EntityGUID: "erroredGUID"}

	s.RecipeAvailable(recipeInstalled)
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
	s := NewInstallStatus([]StatusSubscriber{}, slg)

	require.Equal(t, "localhost:8888", s.HTTPSProxy)
}
