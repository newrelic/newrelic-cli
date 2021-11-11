//build +unit

package execution

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestInstallEventsReporter_InstallStarted(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")

	err := r.InstallStarted(status)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallStatusCallCount)
}

func TestInstallEventsReporter_InstallComplete(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")

	err := r.InstallComplete(status)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallStatusCallCount)
}

func TestInstallEventsReporter_InstallCanceled(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	status.Statuses = []*RecipeStatus{
		{
			Name:   "test-recipe1",
			Status: RecipeStatusTypes.AVAILABLE,
		},
		{
			Name:   "test-recipe2",
			Status: RecipeStatusTypes.AVAILABLE,
		},
		{
			Name:   "test-recipe3",
			Status: RecipeStatusTypes.AVAILABLE,
		},
	}

	err := r.InstallCanceled(status)
	require.NoError(t, err)
	require.Equal(t, 3, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_InstallCanceled_ShouldNotReportDetectedEvent(t *testing.T) {
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	status.Statuses = []*RecipeStatus{
		{
			Name:   "php-agent-installer",
			Status: RecipeStatusTypes.DETECTED, // Not reported when install canceled
		},
		{
			Name:   "aws-integration",
			Status: RecipeStatusTypes.DETECTED, // Not reported when install canceled
		},
		{
			Name:   "logs-integration",
			Status: RecipeStatusTypes.INSTALLING,
		},
		{
			Name:   "mysql-open-source-integration",
			Status: RecipeStatusTypes.CANCELED,
		},
	}

	status.HasCanceledRecipes = true

	err := r.InstallCanceled(status)
	require.NoError(t, err)
	require.Equal(t, 2, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeInstalling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalling(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeFailed(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeInstalled(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeSkipped(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeSkipped(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeRecommended(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeRecommended(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_writeStatus(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstallEventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")

	recipes := []types.OpenInstallationRecipe{
		{
			Name:           types.InfraAgentRecipeName,
			DisplayName:    types.InfraAgentRecipeName,
			ValidationNRQL: "testNrql",
		},
		{
			Name:           types.LoggingRecipeName,
			DisplayName:    types.LoggingRecipeName,
			ValidationNRQL: "testNrql",
		},
	}

	createInstallEventCallCount := 0

	err := r.RecipeAvailable(status, recipes[0])
	createInstallEventCallCount++
	require.NoError(t, err)
	require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)

	for _, testRecipe := range recipes {
		err = r.RecipeAvailable(status, testRecipe)
		createInstallEventCallCount++
		require.NoError(t, err)
		require.Equal(t, createInstallEventCallCount, c.CreateInstallEventCallCount)
	}
}
