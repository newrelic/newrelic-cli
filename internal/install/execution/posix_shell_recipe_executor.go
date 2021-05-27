package execution

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

const (
	BashPath = "/bin/bash"
	ShPath   = "/bin/sh"
)

type PosixShellRecipeExecutor struct {
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	ShellPath string
	Dir       string
}

func NewPosixShellRecipeExecutor() *PosixShellRecipeExecutor {
	return &PosixShellRecipeExecutor{
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		ShellPath: BashPath,
	}
}

func (e *PosixShellRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.Install, v)
}

func (e *PosixShellRecipeExecutor) ExecuteDiscovery(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.PreInstall.RequireAtDiscovery, v)
}

func (e *PosixShellRecipeExecutor) execute(ctx context.Context, script string, v types.RecipeVars) error {
	c := exec.Command(e.ShellPath, "-c", script)

	stdoutCapture := NewLineCaptureBuffer(e.Stdout)
	stderrCapture := NewLineCaptureBuffer(e.Stderr)
	c.Stdout = stdoutCapture
	c.Stderr = stderrCapture
	c.Stdin = e.Stdin
	c.Env = v.ToSlice()
	c.Dir = e.Dir

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			fmt.Println(stderrCapture.LastFullLine)
			re := regexp.MustCompile(".+?: (.+?: )(.*)")
			m := re.FindStringSubmatch(stderrCapture.LastFullLine)

			return fmt.Errorf("%w: %s%s", err, m[1], m[2])
		}

		return err
	}

	return nil
}
