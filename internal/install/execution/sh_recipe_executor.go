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

	log.Print("\n****************************\n")

	log.Printf("LastFullLine: %+v \n", stderrCapture.LastFullLine)
	log.Printf("Error:        %+v \n", err)

	if err != nil {
		if exitCode, ok := interp.IsExitStatus(err); ok {
			return &types.ShError{
				Err:      fmt.Errorf("%w: %s", err, stderrCapture.LastFullLine),
				ExitCode: int(exitCode),
				Metadata: stderrCapture.LastFullLine,
			}
		}

		return err
	}

	return nil
}
