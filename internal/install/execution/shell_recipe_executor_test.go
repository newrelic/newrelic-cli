package execution

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestShellExecution_Basic(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	m := types.DiscoveryManifest{}
	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			ExecDiscovery: "echo 1234",
		},
	}

	err := e.Execute(context.Background(), m, r, v)
	require.NoError(t, err)
	require.Equal(t, "1234\n", stdoutBuf.String())
}

func TestShellExecution_Error(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	m := types.DiscoveryManifest{}
	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			ExecDiscovery: "asdf",
		},
	}

	err := e.Execute(context.Background(), m, r, v)
	require.Error(t, err)
	require.Equal(t, "exit status 127: asdf: command not found", err.Error())
}

func TestShellExecution_RecipeVars(t *testing.T) {
	e := NewPosixShellRecipeExecutor()
	stderrBuf := bytes.NewBufferString("")
	stdoutBuf := bytes.NewBufferString("")
	e.Stdout = stdoutBuf
	e.Stderr = stderrBuf

	m := types.DiscoveryManifest{}
	v := types.RecipeVars{
		"TEST_VAR": "testValue",
	}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			ExecDiscovery: "echo $TEST_VAR",
		},
	}

	err := e.Execute(context.Background(), m, r, v)
	require.NoError(t, err)
	require.Equal(t, "testValue\n", stdoutBuf.String())
}
