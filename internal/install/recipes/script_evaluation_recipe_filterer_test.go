//go:build unit
// +build unit

package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestScriptEvaluationRecipeFilterer_ShouldPassPreInstall(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name:         "test-recipe",
		ProcessMatch: []string{"apache2"},
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "exit 0", // simulate failed preinstall check
		},
	}

	matchedProcess := mockProcess{
		cmdline: "apache2",
		name:    `apache2`,
		pid:     int32(1234),
	}

	m := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
	}

	installStatus := &execution.InstallStatus{
		DiscoveryManifest: m,
	}

	r := NewScriptEvaluationRecipeFilterer(installStatus)

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.NoError(t, err)
}

func TestScriptEvaluationRecipeFilterer_ShouldFailPreInstallWithDetectedEventCaptured(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name:         "test-recipe",
		ProcessMatch: []string{"apache2"},
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			// Using exit code 132 in means DETECTED
			RequireAtDiscovery: `echo '"{\"aUsefulKey\":\"a useful value\"}"' >&2; exit 132`, // simulate failed preinstall check with exit 132
		},
	}

	matchedProcess := mockProcess{
		cmdline: "apache2",
		name:    `apache2`,
		pid:     int32(1234),
	}

	m := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
	}

	mockInstallEventsClient := execution.NewMockInstallEventsClient()
	installEventsReporter := execution.NewInstallEventsReporter(mockInstallEventsClient)

	mockReporter := execution.NewMockStatusReporter()
	statusSubscribers := []execution.StatusSubscriber{installEventsReporter, mockReporter}
	platformLinkGenerator := execution.NewPlatformLinkGenerator()
	installStatus := execution.NewInstallStatus(statusSubscribers, platformLinkGenerator)

	r := NewScriptEvaluationRecipeFilterer(installStatus)

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.Error(t, err)

	require.Equal(t, 1, mockInstallEventsClient.CreateRecipeEventCallCount)
	require.Equal(t, 1, mockReporter.RecipeDetectedCallCount)
}

func TestScriptEvaluationRecipeFilterer_ShouldFailPreInstallWithUnsupportedEventCaptured(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name:         "test-recipe",
		ProcessMatch: []string{"apache2"},
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			// Using exit code 1 should trigger the UNSUPPORTED event
			RequireAtDiscovery: `echo '"{\"aUsefulKey\":\"a useful value\"}"' >&2; exit 1`,
		},
	}

	matchedProcess := mockProcess{
		cmdline: "apache2",
		name:    `apache2`,
		pid:     int32(1234),
	}

	m := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
	}

	mockInstallEventsClient := execution.NewMockInstallEventsClient()
	installEventsReporter := execution.NewInstallEventsReporter(mockInstallEventsClient)

	mockReporter := execution.NewMockStatusReporter()
	statusSubscribers := []execution.StatusSubscriber{installEventsReporter, mockReporter}
	platformLinkGenerator := execution.NewPlatformLinkGenerator()
	installStatus := execution.NewInstallStatus(statusSubscribers, platformLinkGenerator)

	r := NewScriptEvaluationRecipeFilterer(installStatus)

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.Error(t, err)

	require.Equal(t, 1, mockInstallEventsClient.CreateRecipeEventCallCount)
	require.Equal(t, 1, mockReporter.RecipeUnsupportedCallCount)
}
