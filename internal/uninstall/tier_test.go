package uninstall

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestSelectTier_Taskfile(t *testing.T) {
	r := types.OpenInstallationRecipe{
		Uninstall: "default:\n  cmds:\n    - true\n",
	}

	require.Equal(t, TierTaskfile, SelectTier(r))
}

func TestSelectTier_Generic(t *testing.T) {
	r := types.OpenInstallationRecipe{
		UninstallMeta: types.OpenInstallationUninstallMeta{
			Packages: []string{"newrelic-infra"},
		},
	}

	require.Equal(t, TierGeneric, SelectTier(r))
}

func TestSelectTier_NoneWhenEmpty(t *testing.T) {
	r := types.OpenInstallationRecipe{}

	require.Equal(t, TierNone, SelectTier(r))
}

func TestSelectTier_TaskfileWinsOverGeneric(t *testing.T) {
	r := types.OpenInstallationRecipe{
		Uninstall: "default:\n  cmds:\n    - true\n",
		UninstallMeta: types.OpenInstallationUninstallMeta{
			Packages: []string{"newrelic-infra"},
		},
	}

	require.Equal(t, TierTaskfile, SelectTier(r))
}

func TestSelectTier_GenericFromServicesOnly(t *testing.T) {
	r := types.OpenInstallationRecipe{
		UninstallMeta: types.OpenInstallationUninstallMeta{
			Services: []string{"newrelic-infra"},
		},
	}

	require.Equal(t, TierGeneric, SelectTier(r))
}

func TestSelectTier_GenericFromPathsOnly(t *testing.T) {
	r := types.OpenInstallationRecipe{
		UninstallMeta: types.OpenInstallationUninstallMeta{
			Paths: []string{"/etc/newrelic-infra"},
		},
	}

	require.Equal(t, TierGeneric, SelectTier(r))
}
