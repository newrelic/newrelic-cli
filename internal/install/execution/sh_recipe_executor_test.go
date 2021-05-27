package execution

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestExecuteDiscovery_Basic(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo 1234",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "1234\n", b.String())
}

func TestExecuteDiscovery_Error(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "bogus command",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.Error(t, err)
}

func TestExecuteDiscovery_RecipeVars(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	v := types.RecipeVars{
		"TEST_VAR": "testValue",
	}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo $TEST_VAR",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "testValue\n", b.String())
}
