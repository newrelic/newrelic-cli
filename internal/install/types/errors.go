package types

import (
	"errors"
)

// ErrInterrupt represents a context cancellation.
var ErrInterrupt = errors.New("operation canceled")

var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic")
