package execution

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	Stderr       io.Writer
	Stdin        io.Reader
	Stdout       io.Writer
	Output       *OutputParser
	RecipeOutput []string
}

// NewGoTaskRecipeExecutor returns a new instance of GoTaskRecipeExecutor.
func NewGoTaskRecipeExecutor() *GoTaskRecipeExecutor {
	return &GoTaskRecipeExecutor{
		Stderr:       os.Stderr,
		Stdin:        os.Stdin,
		Stdout:       os.Stdout,
		Output:       NewOutputParser(map[string]interface{}{}),
		RecipeOutput: []string{},
	}
}

func (re *GoTaskRecipeExecutor) ExecutePreInstall(ctx context.Context, r types.OpenInstallationRecipe, recipeVars types.RecipeVars) error {
	return errors.New("not implemented")
}

func (re *GoTaskRecipeExecutor) Execute(ctx context.Context, r types.OpenInstallationRecipe, recipeVars types.RecipeVars) (retErr error) {

	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				retErr = errors.New(x)
			case error:
				retErr = x
			default:
				retErr = errors.New("unknown panic")
			}
		}
	}()

	log.Debugf("executing recipe %s", r.Name)

	// unmarshall task file & create/write to temp file
	taskFile, err := createRecipeTempFile(r)
	if err != nil {
		return err
	}
	defer os.Remove(taskFile.Name())

	// Create temp 'flags' file.
	outputJSONFile, err := createOutputJSONFile(r, recipeVars)
	if err != nil {
		return err
	}
	defer outputJSONFile.Close()
	defer os.Remove(outputJSONFile.Name())

	e := task.Executor{
		Dir:        os.TempDir(),
		Entrypoint: filepath.Base(taskFile.Name()),
		Stdin:      re.Stdin,
	}
	if err = e.Setup(); err != nil {
		return fmt.Errorf("could not set up task executor: %s", err)
	}

	calls, globals := taskargs.ParseV3()
	e.Taskfile.Vars.Merge(globals)
	for k, val := range recipeVars {
		e.Taskfile.Vars.Set(k, taskfile.Var{Static: val})
	}

	// configure cli output capture, potentially creating temp file to post logs to New Relic
	// var cliOutputFile *os.File
	// var stdoutCapture *LineCaptureBuffer
	// var stderrCapture *LineCaptureBuffer
	// stdoutCapture, stderrCapture, cliOutputFile, err = configureRecipeOutputCapture(re, r.Name, e.Taskfile.Vars.ToCacheMap(), cliOutputFile)
	// if err != nil {
	// 	return err
	// }
	// if cliOutputFile != nil {
	// 	defer cliOutputFile.Close()
	// }

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
	e.Stdout = stdoutCapture
	e.Stderr = stderrCapture

	if err = e.Run(ctx, calls...); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Debug("Task execution returned error")

		// set flag to point to cli output file
		// if nil != cliOutputFile {
		// 	_, writeErr := outputJSONFile.WriteString(fmt.Sprintf(`{"FailedRecipeOutput":"%s"}`, cliOutputFile.Name()))
		// 	if nil != writeErr {
		// 		log.Debugf("Could not update FailedRecipeOutput flag to reference cli output file: %e", writeErr)
		// 	}
		// }
		re.RecipeOutput = stdoutCapture.GetFullRecipeOutput()
		re.setOutput(outputJSONFile.Name())

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
			lastStderr := stderrCapture.LastFullLine

			return types.NewNonZeroExitCode(goTaskError, lastStderr)
		}

		return goTaskError
	}

	// remove output file immediately since install did not error
	// will be nil for recipes other than infrastructure-agent-installer
	// if nil != cliOutputFile {
	// 	os.Remove(cliOutputFile.Name())
	// }
	re.setOutput(outputJSONFile.Name())

	return nil
}

// func configureRecipeOutputCapture(re *GoTaskRecipeExecutor, recipeName string, recipeVars map[string]interface{}, cliOutputFile *os.File) (*LineCaptureBuffer, *LineCaptureBuffer, *os.File, error) {
// 	silentInstall := false
// 	if assumeYes, ok := recipeVars["assumeYes"]; ok {
// 		silentInstall, _ = strconv.ParseBool(assumeYes.(string))
// 		if silentInstall {
// 			return NewLineCaptureBuffer(&bytes.Buffer{}), NewLineCaptureBuffer(&bytes.Buffer{}), nil, nil
// 		}
// 	}

// 	// Create a temporary cli output file only for non-silent infra installs
// 	sendLogs := false
// 	if captureLogs, ok := recipeVars["CAPTURE_CLI_LOGS"]; ok {
// 		sendLogs, _ = strconv.ParseBool(captureLogs.(string))
// 		if sendLogs {
// 			var err error
// 			cliOutputFile, err = ioutil.TempFile("", fmt.Sprintf("%s_cli_stderr_", recipeName))
// 			if err != nil {
// 				log.Debugf("Could not create temp file for recipe cli output: %e", err)
// 				cliOutputFile = nil
// 			}
// 		}
// 	}
// 	return NewLineCaptureToFileBuffer(re.Stdout, cliOutputFile), NewLineCaptureToFileBuffer(re.Stderr, cliOutputFile), cliOutputFile, nil
// }

func createRecipeTempFile(r types.OpenInstallationRecipe) (*os.File, error) {
	out := []byte(r.Install)
	err := yaml.Unmarshal(out, &taskfile.Taskfile{})
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal taskfile: %s", err)
	}
	taskFile, err := ioutil.TempFile("", r.Name)
	if err != nil {
		return nil, err
	}
	_, err = taskFile.Write(out)
	if err != nil {
		return nil, err
	}
	return taskFile, nil
}

func createOutputJSONFile(r types.OpenInstallationRecipe, recipeVars types.RecipeVars) (*os.File, error) {
	outputJSONFile, err := ioutil.TempFile("", fmt.Sprintf("%s_out", r.Name))
	if err != nil {
		return nil, err
	}
	recipeVars["NR_CLI_OUTPUT"] = outputJSONFile.Name()
	return outputJSONFile, nil
}

func (re *GoTaskRecipeExecutor) setOutput(outputFileName string) {
	outputFile, err := os.Open(outputFileName)
	if err != nil {
		log.Debugf("error openning json output file %s", outputFileName)
		return
	}

	defer outputFile.Close()

	outputBytes, err := ioutil.ReadAll(outputFile)
	if err == nil && len(outputBytes) > 0 {
		var result map[string]interface{}
		if err := json.Unmarshal(outputBytes, &result); err == nil {
			re.Output = NewOutputParser(result)
		} else {
			log.Debugf("error while unmarshaling json output from file %s details:%s", outputFileName, err.Error())
		}
	}
}

func (re *GoTaskRecipeExecutor) GetOutput() *OutputParser {
	return re.Output
}

func isExitStatusCode(exitCode int, err error) bool {
	exitCodeString := fmt.Sprintf("exit status %d", exitCode)
	return strings.Contains(err.Error(), exitCodeString)
}
