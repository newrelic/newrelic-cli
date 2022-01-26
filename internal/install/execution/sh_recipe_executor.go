package execution

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type ShRecipeExecutor struct {
	Dir    string
	Stderr io.Writer
	Stdin  io.Reader
	Stdout io.Writer
}

func NewShRecipeExecutor() *ShRecipeExecutor {
	writer := config.Logger.WriterLevel(log.DebugLevel)
	return &ShRecipeExecutor{
		Stdin:  os.Stdin,
		Stdout: writer,
		Stderr: writer,
	}
}

func (e *ShRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.Install, v)
}

func (e *ShRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	log.Tracef("ExecutePreInstall script for recipe %s", r.Name)
	return e.execute(ctx, r.PreInstall.RequireAtDiscovery, v)
}

func (e *ShRecipeExecutor) execute(ctx context.Context, script string, v types.RecipeVars) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}

	environ := append(os.Environ(), v.ToSlice()...)
	stdoutCapture := NewLineCaptureBuffer(e.Stdout)
	stderrCapture := NewLineCaptureBuffer(e.Stderr)

	i, err := interp.New(
		interp.Params("-e"),
		interp.Dir(e.Dir),
		interp.Env(expand.ListEnviron(environ...)),
		interp.StdIO(e.Stdin, stdoutCapture, stderrCapture),
	)
	if err != nil {
		return err
	}

	err = i.Run(ctx, p)

	if err != nil {
		if _, ok := interp.IsExitStatus(err); ok {
			// If the stderr message is a regular string, we return the original error
			// and last full line of text. This is the original way recipes send messages
			// via stderr, hence we need this check for backwards compatibility.
			if !utils.IsJSONString(stderrCapture.LastFullLine) {
				return fmt.Errorf("%w: %s", err, stderrCapture.LastFullLine)
			}

			return types.NewCustomStdError(err, stderrCapture.LastFullLine)

			// When a recipe returns the stderr message is a JSON string we
			// capture the additional metadata for informational purposes.
			// return &types.CustomStdError{
			// 	Message:  fmt.Sprintf("%s: %s", err, stderrCapture.LastFullLine),
			// 	ExitCode: int(exitCode),
			// 	Metadata: stderrCapture.LastFullLine,
			// }
		}

		return err
	}

	// Handle when a recipe sends a JSON string via stderr even if no error occurred.
	// This can occur when a recipe executes a step successfully but still wants to capture
	// metadata in the recipe event.
	// if stderrCapture.LastFullLine != "" && utils.IsJSONString(stderrCapture.LastFullLine) {
	// 	return &types.CustomStdError{
	// 		Message:  fmt.Sprintf("%s: %s", err, stderrCapture.LastFullLine),
	// 		Metadata: stderrCapture.LastFullLine,
	// 	}
	// }

	return nil
}
