package execution

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
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

func (e *ShRecipeExecutor) Execute(ctx context.Context, m types.DiscoveryManifest, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	p, err := syntax.NewParser().Parse(strings.NewReader(r.PreInstall.ExecDiscovery), "")
	if err != nil {
		return err
	}

	var environ []string
	if len(v) == 0 {
		environ = os.Environ()
	} else {
		environ = v.ToSlice()
	}

	i, err := interp.New(
		interp.Params("-e"),
		interp.Dir(e.Dir),
		interp.Env(expand.ListEnviron(environ...)),
		interp.StdIO(e.Stdin, e.Stdout, e.Stderr),
	)
	if err != nil {
		return err
	}

	return i.Run(ctx, p)
}

func (e *ShRecipeExecutor) Prepare(ctx context.Context, manifest types.DiscoveryManifest, recipe types.OpenInstallationRecipe, assumeYes bool, licenseKey string) (types.RecipeVars, error) {
	return nil, nil
}
