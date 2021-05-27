package execution

import (
	"bytes"
	"context"
	"os"
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
			RequireAtDiscovery: "bogus",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.Error(t, err)
	require.Equal(t, "exit status 127: \"bogus\": executable file not found in $PATH", err.Error())
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

func TestExecuteDiscovery_EnvVars(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	os.Setenv("ENV_VAR", "envVarValue")
	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo $ENV_VAR",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "envVarValue\n", b.String())
}

func TestExecuteDiscovery_AllVars(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	os.Setenv("ENV_VAR", "envVarValue")
	v := types.RecipeVars{
		"RECIPE_VAR": "recipeVarValue",
	}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: `
			echo $ENV_VAR
			echo $RECIPE_VAR
			`,
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Contains(t, b.String(), "envVarValue")
	require.Contains(t, b.String(), "recipeVarValue")
}
