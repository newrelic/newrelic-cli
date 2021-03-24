package types

import (
	"errors"
)

// ErrInterrupt represents a context cancellation.
var ErrInterrupt = errors.New("operation canceled")

var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://developer.newrelic.com/automate-workflows/get-started-new-relic-cli")
