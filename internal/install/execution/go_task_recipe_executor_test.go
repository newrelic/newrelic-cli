package execution

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestExecute_WritesCliOutputVar(t *testing.T) {
	v := types.RecipeVars{
		"TEST_VAR":         "testValue",
		"CAPTURE_CLI_LOGS": "true",
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
	_, containsKey := e.GetOutput().output["CapturedCliOutput"]
	assert.Truef(t, containsKey, "Output file missing key for 'CapturedCliOutput")
	require.Error(t, err)
	require.Contains(t, err.Error(), "testValue")
}

func TestExecute_RecipesGetTheirOwnMetadata(t *testing.T) {
	firstRecipeExecuted := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
       - |
         echo '{"Metadata":{"first-recipe":"firstRecipeVal"}}' | tee {{.NR_CLI_OUTPUT}} > /dev/null 
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), firstRecipeExecuted, types.RecipeVars{})
	require.NoError(t, err)
	val := e.GetOutput().Metadata()["first-recipe"]
	assert.Equal(t, "firstRecipeVal", val)
	assert.Len(t, e.GetOutput().Metadata(), 1)

	secondRecipeExecuted := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
       - |
         echo no metadata here 
`,
	}

	err = e.Execute(context.Background(), secondRecipeExecuted, types.RecipeVars{})
	require.NoError(t, err)
	_, firstKeyPresent := e.GetOutput().Metadata()["first-recipe"]
	assert.False(t, firstKeyPresent)
	assert.Len(t, e.GetOutput().Metadata(), 0)

	thirdRecipeExecuted := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
       - |
         echo '{"Metadata":{"third-recipe":"thirdRecipeVal"}}' | tee {{.NR_CLI_OUTPUT}} > /dev/null 
`,
	}

	err = e.Execute(context.Background(), thirdRecipeExecuted, types.RecipeVars{})
	require.NoError(t, err)
	_, firstKeyPresent = e.GetOutput().Metadata()["first-recipe"]
	assert.False(t, firstKeyPresent)
	val = e.GetOutput().Metadata()["third-recipe"]
	assert.Equal(t, "thirdRecipeVal", val)
	assert.Len(t, e.GetOutput().Metadata(), 1)

}
