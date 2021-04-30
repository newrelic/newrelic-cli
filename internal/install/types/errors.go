package types

import (
	"errors"
)

// ErrInterrupt represents a context cancellation.
var ErrInterrupt = errors.New("operation canceled")

// ErrNonZeroExitCode represents a non-zero exit code bieing returned from a child process.
type ErrNonZeroExitCode struct {
	err string
}

func NewNonZeroExitCode(err string) ErrNonZeroExitCode {
	return ErrNonZeroExitCode{
		err: err,
	}
}

func (e ErrNonZeroExitCode) Error() string {
	return e.err
}

// nolint: golint
var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic")
var ErrorFetchingInsightsInsertKey = errors.New("error retrieving Insights insert key")
