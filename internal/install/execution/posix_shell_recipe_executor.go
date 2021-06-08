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
	Stderr    io.Writer
	Stdin     io.Reader
	Stdout    io.Writer
	ShellPath string
	Dir       string
}

func NewPosixShellRecipeExecutor() *PosixShellRecipeExecutor {
	return &PosixShellRecipeExecutor{
		Stderr:    os.Stderr,
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		ShellPath: BashPath,
	}
}

func (e *PosixShellRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.Install, v)
}

func (e *PosixShellRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, v types.RecipeVars) error {
	return e.execute(ctx, r.PreInstall.RequireAtDiscovery, v)
}

func (e *PosixShellRecipeExecutor) execute(ctx context.Context, script string, v types.RecipeVars) error {
	c := exec.Command(e.ShellPath, "-c", script)

	stdoutCapture := NewLineCaptureBuffer(e.Stdout)
	stderrCapture := NewLineCaptureBuffer(e.Stderr)

	c.Dir = e.Dir
	c.Env = append(os.Environ(), v.ToSlice()...)
	c.Stderr = stderrCapture
	c.Stdin = e.Stdin
	c.Stdout = stdoutCapture

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError
		if ok := errors.As(err, &exitError); ok {
			fmt.Println(stderrCapture.LastFullLine)
			re := regexp.MustCompile(".+?: (.*)")
			m := re.FindStringSubmatch(stderrCapture.LastFullLine)

			return fmt.Errorf("%w: %s", err, m[1])
		}

		return err
	}

	return nil
}
