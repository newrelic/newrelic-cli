package execution

import (
	"context"
	"encoding/json"
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

	fmt.Print("\n\n **************************** \n")
	fmt.Printf("\n ShRecipeExecutor - stderrCapture:  %+v \n", stderrCapture.LastFullLine)
	fmt.Print("\n **************************** \n\n")
	// time.Sleep(5 * time.Second)

	if err != nil {
		if exitCode, ok := interp.IsExitStatus(err); ok {
			// If the stderr message is a regular string, we return
			// the original error and last full line of text.
			// This is the original way recipes send messages via stderr,
			// hence we need this check for backwards compatibility.
			if !isJSONString(stderrCapture.LastFullLine) {
				return fmt.Errorf("%w: %s", err, stderrCapture.LastFullLine)
			}

			// When the stderr message is a JSON string, we have a standard object
			// we return a custom error to facilitate capturing additional metadata
			// for debugging purposes.
			return &types.IncomingMessage{
				// Should we use fmt.Errorf here ever? Should we have another field to cover error?
				Message:  fmt.Sprintf("%s: %s", err, stderrCapture.LastFullLine),
				ExitCode: int(exitCode),
				Metadata: stderrCapture.LastFullLine,
			}
		}

		return err
	}

	// Handle scenario when no error occurs, but we still pass
	// a message back
	if stderrCapture.LastFullLine != "" {
		return &types.IncomingMessage{
			Message:  fmt.Sprintf("%s: %s", err, stderrCapture.LastFullLine),
			Metadata: stderrCapture.LastFullLine,
		}
	}

	return nil
}

func isJSONString(input string) bool {
	var x struct{}
	if err := json.Unmarshal([]byte(input), &x); err != nil {
		return false
	}

	return true
}
