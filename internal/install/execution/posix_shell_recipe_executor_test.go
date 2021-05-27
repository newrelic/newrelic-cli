package execution

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestPosixShellExecution_Basic(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo 1234",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "1234\n", stdoutBuf.String())
}

func TestPosixShellExecution_Error(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "asdf",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.Error(t, err)
	require.Equal(t, "exit status 127: asdf: command not found", err.Error())
}

func TestPosixShellExecution_AllVars(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	os.Setenv("ENV_VAR", "envVarValue")
	v := types.RecipeVars{
		"RECIPE_VAR": "recipeVarValue",
	}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: `
			echo $ENV_VAR
			echo $RECIPE_VAR
			`,
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Contains(t, stdoutBuf.String(), "recipeVarValue")
	require.Contains(t, stdoutBuf.String(), "envVarValue")
}

func TestPosixShellExecution_RecipeVars(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	v := types.RecipeVars{
		"RECIPE_VAR": "recipeVarValue",
	}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo $RECIPE_VAR",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "recipeVarValue\n", stdoutBuf.String())
}

func TestPosixShellExecution_EnvVars(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	os.Setenv("ENV_VAR", "envVarValue")

	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo $ENV_VAR",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "envVarValue\n", stdoutBuf.String())
}
