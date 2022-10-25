package execution

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestExecute_SystemVariableInterpolation(t *testing.T) {
	v := types.RecipeVars{
		"TEST_VAR": "testValue",
	}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
    cmds:
      - |
        echo {{.TEST_VAR}}
`,
	}

	e := NewGoTaskRecipeExecutor()
	b := bytes.NewBufferString("")
	e.Stdout = b
	err := e.Execute(context.Background(), r, v)
	require.NoError(t, err)
	require.True(t, strings.Contains(b.String(), "echo testValue"))
}

func TestExecute_HandleRecipeLastError(t *testing.T) {
	v := types.RecipeVars{
		"TEST_VAR": "testValue",
	}
	r := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
       - |
         echo {{.TEST_VAR}} >&2
         exit 1
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), r, v)
	require.Error(t, err)
	require.Contains(t, err.Error(), "testValue")
}
