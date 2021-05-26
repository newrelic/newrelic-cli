package execution

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestExecution_Basic(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	m := types.DiscoveryManifest{}
	v := types.RecipeVars{}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			ExecDiscovery: "echo 1234",
		},
	}

	err := e.Execute(context.Background(), m, r, v)
	require.NoError(t, err)
	require.Equal(t, "1234\n", b.String())
}

func TestExecution_RecipeVars(t *testing.T) {
	e := NewShRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b

	m := types.DiscoveryManifest{}
	v := types.RecipeVars{
		"TEST_VAR": "testValue",
	}
	r := types.OpenInstallationRecipe{
		PreInstall: types.OpenInstallationPreInstallConfiguration{
			ExecDiscovery: "echo $TEST_VAR",
		},
	}

	err := e.Execute(context.Background(), m, r, v)
	require.NoError(t, err)
	require.Equal(t, "testValue\n", b.String())
}
