package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrInterrupt represents a context cancellation.
	ErrInterrupt       = errors.New("operation canceled")
	maxTaskPathNesting = 5
)

type GoTaskError interface {
	error
	TaskPath() []string
	SetError(msg string)
}

// GoTaskError represents a task failure reported by go-task.
type GoTaskGeneralError struct {
	error
	taskPath []string
}

func (e GoTaskGeneralError) Error() string {
	return e.error.Error()
}

func (e GoTaskGeneralError) SetError(msg string) {
	//nolint:staticcheck
	e.error = errors.New(msg)
}

func (e GoTaskGeneralError) TaskPath() []string {
	return e.taskPath
}

func NewGoTaskGeneralError(err error) GoTaskError {
	pattern := `task: Failed to run task \"(.+?)\": `
	str := strings.Repeat("(?:%[1]s)?", maxTaskPathNesting)
	re := regexp.MustCompile(fmt.Sprintf(str, pattern))
	parsed := re.FindStringSubmatch(err.Error())

	taskPath := []string{}
	for i, p := range parsed {
		if i == 0 {
			continue
		}

		if p != "" {
			taskPath = append(taskPath, p)
		}
	}

	stripped := re.ReplaceAllString(err.Error(), "")

	return GoTaskGeneralError{
		taskPath: taskPath,
		error:    errors.New(stripped),
	}
}

// ErrNonZeroExitCode represents a non-zero exit code error reported by go-task.
type ErrNonZeroExitCode struct {
	GoTaskError
	additionalContext string
}

func NewNonZeroExitCode(originalError GoTaskError, additionalContext string) ErrNonZeroExitCode {
	return ErrNonZeroExitCode{
		GoTaskError:       originalError,
		additionalContext: additionalContext,
	}
}

func (e ErrNonZeroExitCode) Error() string {
	if e.additionalContext != "" {
		return fmt.Sprintf("%s: %s", e.GoTaskError.Error(), e.additionalContext)
	}

	return e.GoTaskError.Error()
}

type UnsupportedOperatingSytemError struct {
	Err error
}

func (e *UnsupportedOperatingSytemError) Error() string {
	return e.Err.Error()
}

// UpdateRequiredError represents when a user is using an older version
// of the CLI and is required to update when running the `newrelic install` command.
type UpdateRequiredError struct {
	Err     error
	Details string
}

func (e *UpdateRequiredError) Error() string {
	return e.Err.Error()
}

type UncaughtError struct {
	Err error
}

func (e *UncaughtError) Error() string {
	return e.Err.Error()
}

// IncomingMessage represents a standardized recipe message object
// passed back to the CLI via stderr. The standard error is represented as
// a JSON string within a recipe.
//
// Example of:
//   echo ""{\"message\":\"something happened\",\"metadata\":{\"key\":\"relevant data\"}}"" >&2
//
//
type IncomingMessage struct {
	// The primary message that was redirected to stderr
	Message string

	// The exit code used at point of failure in the recipe
	ExitCode int

	// JSON string that contains additional information if the
	// recipe provides it via stderr. Use IncomingMessage.ParseMetadata()
	// access the data in Go.
	Metadata string
}

func (e IncomingMessage) Error() string {
	return e.Message
}

// ParseMetadata converts the incoming JSON string to a map[string]interface{}.
// If the incoming metadata is a simple string, we still return `metadata` as
// a map to keep data structure consistent.
func (e IncomingMessage) ParseMetadata() map[string]interface{} {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(e.Metadata), &data); err != nil {
		log.Debugf("\n Could not unmarshal e.Metadata:  %+v \n", err)

		fmt.Printf("\b ParseMetadata didnt unmarshal:   %+v \n", e.Metadata)

		return map[string]interface{}{
			"metadata": map[string]string{
				"message": e.Metadata,
			},
		}
	}

	if m, ok := data["metadata"].(string); ok {
		return map[string]interface{}{
			"metadata": map[string]string{
				"message": m,
			},
		}
	}

	if m, ok := data["metadata"].(map[string]interface{}); ok {
		return m
	}

	return data
}

// nolint: golint
var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic")
var ErrorFetchingInsightsInsertKey = errors.New("error retrieving Insights insert key")
