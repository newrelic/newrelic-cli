package execution

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestExecute_SystemVariableInterpolation(t *testing.T) {
	//r := NewRecipeBuilder().Vars("TEST_VAR", "testValue").InstallShell("echo {{.TEST_VAR}}").Build()
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
	require.Equal(t, "testValue\n", b.String())
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

func TestExecute_HandleInfrastructureRecipeFileCreation(t *testing.T) {
	v := types.RecipeVars{
		"TEST_VAR":         "testValue",
		"CAPTURE_CLI_LOGS": "true",
	}
	r := types.OpenInstallationRecipe{
		Name: "infrastructure-agent-installer",
		Install: `
version: '3'
tasks:
  default:
     cmds:
       - |
         echo "this is random stdout noise"
         echo some error text >&2
         echo {{.TEST_VAR}} >&2
         exit 1
`,
	}
	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), r, v)

	require.Error(t, err)
	require.Contains(t, err.Error(), "testValue")

	require.NotEqual(t, "", e.GetOutput().FailedRecipeOutput())
	_, fileErr := os.Stat(e.GetOutput().FailedRecipeOutput())
	require.NoError(t, fileErr)

	data, fileErr := os.ReadFile(e.GetOutput().FailedRecipeOutput())
	require.NoError(t, fileErr)
	fileContents := string(data)
	require.True(t, strings.Contains(fileContents, "this is random stdout noise"))
	require.True(t, strings.Contains(fileContents, "some error text"))
	require.True(t, strings.Contains(fileContents, "testValue"))
}
