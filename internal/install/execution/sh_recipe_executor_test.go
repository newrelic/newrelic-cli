package execution

import (
	"bytes"
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/stretchr/testify/require"
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
	require.Equal(t, "1234\n", string(b.Bytes()))
}
