//go:build unit
// +build unit

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

func TestExecute_AllMetadataKeysAreCollectedFromRecipe_appending(t *testing.T) {

	// recipe using tee to append metadata
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"something":"else"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 3)
	assert.Equal(t, "thing", e.GetOutput().Metadata()["some"])
	assert.Equal(t, "else", e.GetOutput().Metadata()["something"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])

	// recipe using echo to append metadata
	recipe = types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' >>  {{.NR_CLI_OUTPUT}}

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"something":"else"}}' >> {{.NR_CLI_OUTPUT}} 

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' >> {{.NR_CLI_OUTPUT}} 
`,
	}

	e = NewGoTaskRecipeExecutor()
	err = e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 3)
	assert.Equal(t, "thing", e.GetOutput().Metadata()["some"])
	assert.Equal(t, "else", e.GetOutput().Metadata()["something"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])
}

func TestExecute_ComplexMetadataKeysCollectedFromRecipeAreIgnored_appending(t *testing.T) {

	// recipe using tee to append metadata
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"parents": { "mom":{"first":"linda", "last":"jones"}, "dad":{"first":"otto", "last":"jones"}} }}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 2)
	assert.Equal(t, "thing", e.GetOutput().Metadata()["some"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])

	// recipe using echo to append metadata
	recipe = types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' >>  {{.NR_CLI_OUTPUT}}

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"parents": { "mom":{"first":"linda", "last":"jones"}, "dad":{"first":"otto", "last":"jones"}} }}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' >> {{.NR_CLI_OUTPUT}} 
`,
	}

	e = NewGoTaskRecipeExecutor()
	err = e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 2)
	assert.Equal(t, "thing", e.GetOutput().Metadata()["some"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])
}

func TestExecute_MalformedJSONMetadataKeysAreSkipped(t *testing.T) {

	// first metadata is bad
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing went horribly wrong}}}}}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"something":"else"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 2)
	assert.Empty(t, e.GetOutput().Metadata()["some"])
	assert.Equal(t, "else", e.GetOutput().Metadata()["something"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])

	// second metadata is bad
	recipe = types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"somethi' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null
`,
	}

	e = NewGoTaskRecipeExecutor()
	err = e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 2)
	assert.Equal(t, "thing", e.GetOutput().Metadata()["some"])
	assert.Empty(t, e.GetOutput().Metadata()["something"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])

}

func TestExecute_LastMetadataChildKeyWins(t *testing.T) {
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: write_meta_again_again

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' | tee -a  {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"some":"more"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again_again:
     cmds:
      - |
        echo '{"Metadata":{"some":"last-write"}}' | tee -a {{.NR_CLI_OUTPUT}} > /dev/null
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 1)
	assert.Equal(t, "last-write", e.GetOutput().Metadata()["some"])
}

func TestExecute_AllMetadataKeysAreCollectedFromRecipe_notAppending(t *testing.T) {
	// recipe using tee to write metadata, not append
	recipe := types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' | tee  {{.NR_CLI_OUTPUT}} > /dev/null

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"something":"else"}}' | tee  {{.NR_CLI_OUTPUT}} > /dev/null

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' | tee {{.NR_CLI_OUTPUT}} > /dev/null
`,
	}

	e := NewGoTaskRecipeExecutor()
	err := e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 1)
	assert.Empty(t, e.GetOutput().Metadata()["some"])
	assert.Empty(t, e.GetOutput().Metadata()["something"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])

	// recipe using echo to write metadata, not append
	recipe = types.OpenInstallationRecipe{
		Name: "test-recipe",
		Install: `
version: '3'
tasks:
  default:
     cmds:
      - task: write_meta
      - task: write_meta_again
      - task: cleanup

  write_meta:
     cmds:
      - |
        echo '{"Metadata":{"some":"thing"}}' > {{.NR_CLI_OUTPUT}} 

  write_meta_again:
     cmds:
      - |
        echo '{"Metadata":{"something":"else"}}' > {{.NR_CLI_OUTPUT}}

  cleanup:
     cmds:
      - |
        echo '{"Metadata":{"clean":"very"}}' > {{.NR_CLI_OUTPUT}} 
`,
	}

	e = NewGoTaskRecipeExecutor()
	err = e.Execute(context.Background(), recipe, types.RecipeVars{})
	require.NoError(t, err)
	assert.Len(t, e.GetOutput().Metadata(), 1)
	assert.Empty(t, e.GetOutput().Metadata()["some"])
	assert.Empty(t, e.GetOutput().Metadata()["something"])
	assert.Equal(t, "very", e.GetOutput().Metadata()["clean"])
}
