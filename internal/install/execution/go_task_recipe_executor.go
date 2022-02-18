package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/go-task/task/v3"
	taskargs "github.com/go-task/task/v3/args"
	"github.com/go-task/task/v3/taskfile"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// GoTaskRecipeExecutor is an implementation of the recipeExecutor interface that
// uses the go-task module to execute the steps defined in each recipe.
type GoTaskRecipeExecutor struct {
	Stderr        io.Writer
	Stdin         io.Reader
	Stdout        io.Writer
	OutputCapture *LineCaptureBuffer
	ErrorCapture  *LineCaptureBuffer
}

// NewGoTaskRecipeExecutor returns a new instance of GoTaskRecipeExecutor.
func NewGoTaskRecipeExecutor() *GoTaskRecipeExecutor {
	return &GoTaskRecipeExecutor{
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
	}
}

func (re *GoTaskRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, recipeVars types.RecipeVars) error {
	return errors.New("not implemented")
}

func (re *GoTaskRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, recipeVars types.RecipeVars) error {
	log.Debugf("executing recipe %s", r.Name)

	out := []byte(r.Install)

	// Create a temporary task file.
	file, err := ioutil.TempFile("", r.Name)
	defer os.Remove(file.Name())
	if err != nil {
		return err
	}

	_, err = file.Write(out)
	if err != nil {
		return err
	}

	//FIXME: do something with output in buffer
	//TODO: allow silent flag to be pass in here

	silentInstall, _ := strconv.ParseBool(recipeVars["assumeYes"])

	var stdoutCapture *LineCaptureBuffer
	var stderrCapture *LineCaptureBuffer

	if silentInstall {
		stdoutCapture = NewLineCaptureBuffer(&bytes.Buffer{})
		stderrCapture = NewLineCaptureBuffer(&bytes.Buffer{})
	} else {
		stdoutCapture = NewLineCaptureBuffer(re.Stdout)
		stderrCapture = NewLineCaptureBuffer(re.Stderr)
	}

	e := task.Executor{
		Entrypoint: file.Name(),
		Stderr:     stderrCapture,
		Stdin:      re.Stdin,
		Stdout:     stdoutCapture,
	}

	if err = e.Setup(); err != nil {
		return fmt.Errorf("could not set up task executor: %s", err)
	}

	var tf taskfile.Taskfile
	err = yaml.Unmarshal(out, &tf)
	if err != nil {
		return fmt.Errorf("could not unmarshal taskfile: %s", err)
	}

	calls, globals := taskargs.ParseV3()
	e.Taskfile.Vars.Merge(globals)

	for k, val := range recipeVars {
		e.Taskfile.Vars.Set(k, taskfile.Var{Static: val})
	}

	if err := e.Run(ctx, calls...); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("Task execution returned error")

		goTaskError := types.NewGoTaskGeneralError(err)

		// go-task does not provide an error type to denote context cancelation
		// Therefore we need to match inside the error message
		if strings.Contains(err.Error(), "context canceled") {
			return types.ErrInterrupt
		}

		// Recipe trigger a canceled event with specific exit code 130 used
		if isExitStatusCode(130, err) {
			return types.ErrInterrupt
		}

		// We return exit code 131 when a user attempts to
		// install a recipe on an unsupported operating system.
		if isExitStatusCode(131, err) {

			return &types.UnsupportedOperatingSystemError{
				Err: errors.New(stderrCapture.LastFullLine),
			}
		}

		// Catchall error formatting for child process errors
		if strings.Contains(err.Error(), "exit status") {
			lastStderr := re.ErrorCapture.LastFullLine

			return types.NewNonZeroExitCode(goTaskError, lastStderr)
		}

		return goTaskError
	}

	return nil
}

func isExitStatusCode(exitCode int, err error) bool {
	exitCodeString := fmt.Sprintf("exit status %d", exitCode)
	return strings.Contains(err.Error(), exitCodeString)
}
