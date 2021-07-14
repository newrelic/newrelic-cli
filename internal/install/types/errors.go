package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
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

// nolint: golint
var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic")
var ErrorFetchingInsightsInsertKey = errors.New("error retrieving Insights insert key")
