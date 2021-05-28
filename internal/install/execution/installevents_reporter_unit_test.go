// +build unit

package execution

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestInstallEventsReporter_RecipeFailed(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewConcreteSuccessLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeFailed(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)

}

func TestInstallEventsReporter_RecipeInstalling(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewConcreteSuccessLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalling(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeInstalled(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewConcreteSuccessLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeInstalled(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeSkipped(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewConcreteSuccessLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeSkipped(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}

func TestInstallEventsReporter_RecipeRecommended(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	c := NewMockInstalleventsClient()
	r := NewInstallEventsReporter(c)
	require.NotNil(t, r)

	slg := NewConcreteSuccessLinkGenerator()
	status := NewInstallStatus([]StatusSubscriber{}, slg)
	status.withEntityGUID("testGuid")
	e := RecipeStatusEvent{}

	err := r.RecipeRecommended(status, e)
	require.NoError(t, err)
	require.Equal(t, 1, c.CreateInstallEventCallCount)
}
