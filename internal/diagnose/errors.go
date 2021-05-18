package diagnose

import "errors"

var (
	ErrConnection        = errors.New("There was an error connecting to New Relic platform. This is an indication that your firewall or proxy settings do not allow outbound traffic to the New Relic platform. To configure the use of an HTTP proxy, use the HTTP_PROXY and/or HTTPS_PROXY environment variables. For more details visit https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md#using-an-http-proxy")
	ErrValidation        = errors.New("There was a failure locating the data that was posted to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
	ErrDiscovery         = errors.New("Failed to detect your system's hostname. Please contact New Relic support.")
	ErrPostEvent         = errors.New("There was a failure posting data to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
	ErrLicenseKey        = errors.New("The configured license key is invalid for the configured account. Please set a valid license key with the `newrelic profile` command. For more details visit https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#ingest-license-key")
	ErrInsightsInsertKey = errors.New("The configured Insights insert key is invalid for the configured account. Please set a valid Insights insert key with the `newrelic profile` command. For more details visit https://docs.newrelic.com/docs/telemetry-data-platform/ingest-apis/introduction-event-api/#register")
)
