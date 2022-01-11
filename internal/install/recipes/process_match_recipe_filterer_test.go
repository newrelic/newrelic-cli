//go:build unit
// +build unit

package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestProcessMatchRecipeFilterer_ShouldMatchProcess(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name:         "test-recipe",
		ProcessMatch: []string{"php-fpm"},
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "exit 0",
		},
	}
	matchedProcess := mockProcess{
		cmdline: "php-fpm",
		name:    `php-fpm`,
		pid:     int32(1234),
	}

	m := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
	}

	r := NewProcessMatchRecipeFilterer()

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.NoError(t, err)
}

func TestProcessMatchRecipeFilterer_ShouldNotMatchProcess(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name:         "test-recipe",
		ProcessMatch: []string{"node"},
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "exit 0",
		},
	}
	matchedProcess := mockProcess{
		cmdline: "php-fpm",
		name:    `php-fpm`,
		pid:     int32(1234),
	}

	m := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{matchedProcess},
	}

	r := NewProcessMatchRecipeFilterer()

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.Error(t, err)
}

func TestProcessMatchRecipeFilterer_ShouldIgnoreIfNoProcessMatchConfigured(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name:         "test-recipe",
		ProcessMatch: []string{}, // empty process match array in recipe should be ignored
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "exit 0",
		},
	}

	m := types.DiscoveryManifest{
		DiscoveredProcesses: []types.GenericProcess{},
	}

	r := NewProcessMatchRecipeFilterer()

	err := r.CheckCompatibility(context.Background(), &recipe, &m)
	require.NoError(t, err)
}
