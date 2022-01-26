//go:build unit
// +build unit

package recipes

import (
	"context"
	"fmt"
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
			RequireAtDiscovery: `echo "{\"aUsefulKey\":\"a useful value\"}" >&2; exit 132`, // simulate failed preinstall check
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

	// mockReporter := execution.NewMockStatusReporter()
	statusSubscribers := []execution.StatusSubscriber{installEventsReporter}
	platformLinkGenerator := execution.NewPlatformLinkGenerator()
	installStatus := execution.NewInstallStatus(statusSubscribers, platformLinkGenerator)

	r := NewScriptEvaluationRecipeFilterer(installStatus)

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.Error(t, err)

	// c := statusSubscribers[0].(*execution.InstallEventsReporter)
	fmt.Print("\n\n **************************** \n")
	fmt.Printf("\n TEST - err:  %+v \n", err)
	fmt.Print("\n **************************** \n\n")

	// TODO: Figure out how to test the value of `metadata` in the event itself
	require.Equal(t, 1, mockInstallEventsClient.CreateRecipeEventCallCount)
}
