package execution

import (
	"bytes"
	"context"
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

func TestPosixShellExecution_RecipeVars(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	v := types.RecipeVars{
		"TEST_VAR": "testValue",
	}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			RequireAtDiscovery: "echo $TEST_VAR",
		},
	}

	err := e.ExecuteDiscovery(context.Background(), r, v)
	require.NoError(t, err)
	require.Equal(t, "testValue\n", stdoutBuf.String())
}
