package types

import (
	"errors"
	"fmt"
	"regexp"
)

// ErrInterrupt represents a context cancellation.
var ErrInterrupt = errors.New("operation canceled")

type GoTaskError interface {
	error
	Tasks() []string
	SetError(msg string)
}

// GoTaskError represents a task failure reported by go-task.
type GoTaskGeneralError struct {
	error
	tasks []string
}

func (e GoTaskGeneralError) Error() string {
	return e.error.Error()
}

func (e GoTaskGeneralError) SetError(msg string) {
	e.error = errors.New(msg)
}

func (e GoTaskGeneralError) Tasks() []string {
	return e.tasks
}

func NewGoTaskGeneralError(err error) GoTaskError {
	re := regexp.MustCompile(`task: Failed to run task \"default\": task: Failed to run task \"(.+)\": `)

	parsed := re.FindAllSubmatch([]byte(err.Error()), 1)

	var task string
	if len(parsed) > 0 && len(parsed[0]) > 0 {
		task = string(parsed[0][1])
	}

	stripped := re.ReplaceAllString(err.Error(), "")

	return GoTaskGeneralError{
		tasks: []string{task},
		error: errors.New(stripped),
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

// nolint: golint
var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic")
var ErrorFetchingInsightsInsertKey = errors.New("error retrieving Insights insert key")
