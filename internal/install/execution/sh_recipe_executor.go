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

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

type ShRecipeExecutor struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Dir    string
}

func NewShRecipeExecutor() *ShRecipeExecutor {
	return &ShRecipeExecutor{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (e *ShRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.Install, v)
}

func (e *ShRecipeExecutor) ExecuteDiscovery(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.PreInstall.RequireAtDiscovery, v)
}

func (e *ShRecipeExecutor) execute(ctx context.Context, script string, v types.RecipeVars) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return err
	}

	var environ []string
	if len(v) == 0 {
		environ = os.Environ()
	} else {
		environ = v.ToSlice()
	}

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

	if err := i.Run(ctx, p); err != nil {
		if _, ok := interp.IsExitStatus(err); ok {
			return fmt.Errorf("%w: %s", err, stderrCapture.LastFullLine)
		}

		return err
	}

	return nil
}

func (e *ShRecipeExecutor) Prepare(ctx context.Context, manifest types.DiscoveryManifest, recipe types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error) {
	return nil, nil
}
