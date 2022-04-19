package diagnose

import "errors"

//nolint:golint
var (
	ErrConnectionStringFormat = "there was an error connecting to New Relic platform. This is an indication that your firewall or proxy settings do not allow outbound traffic to the New Relic platform. To configure the use of an HTTP proxy, use the HTTP_PROXY and/or HTTPS_PROXY environment variables. For more details visit https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md#using-an-http-proxy. Details: %s"
	ErrValidation             = errors.New("there was a failure locating the data that was posted to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
	ErrDiscovery              = errors.New("failed to detect your system's hostname. Please contact New Relic support")
	ErrPostEvent              = errors.New("there was a failure posting data to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
	ErrLicenseKey             = errors.New("the configured license key is invalid for the configured account. Please set a valid license key with the `newrelic profile` command. For more details visit https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#ingest-license-key")
)

const PaymentRequiredExceptionMessage = `
  Your account has exceeded your plan's data limit.
  Take full advantage of New Relic's platform by managing
  the account's plan and payment options at the URL below.`
