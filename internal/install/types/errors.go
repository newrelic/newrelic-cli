package types

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrInterrupt represents a context cancellation.
	ErrInterrupt              = errors.New("operation canceled")
	maxTaskPathNesting        = 5
	ErrConnectionStringFormat = "there was an error connecting to New Relic platform. This is an indication that your firewall or proxy settings do not allow outbound traffic to the New Relic platform. To configure the use of an HTTP proxy, use the HTTP_PROXY and/or HTTPS_PROXY environment variables. For more details visit https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md#using-an-http-proxy. Details: %s"
	ErrValidation             = errors.New("there was a failure locating the data that was posted to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
	ErrDiscovery              = errors.New("failed to detect your system's hostname. Please contact New Relic support")
	ErrPostEvent              = errors.New("there was a failure posting data to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
	ErrLicenseKey             = errors.New("the configured license key is invalid for the configured account. Please set a valid license key with the `newrelic profile` command. For more details visit https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#ingest-license-key")
)

type EventType string

var EventTypes = struct {
	InstallStarted             EventType
	AccountIDMissing           EventType
	APIKeyMissing              EventType
	RegionMissing              EventType
	UnableToConnect            EventType
	UnableToFetchLicenseKey    EventType
	LicenseKeyFetchedOk        EventType
	UnableToOverrideClient     EventType
	UnableToPostData           EventType
	InstallCompleted           EventType
	InstallCancelled           EventType
	InvalidIngestKey           EventType
	UnableToDiscover           EventType
	NrIntegrationPollingErrror EventType
	OtherError                 EventType
	UnableToLocatePostedData   EventType
}{
	InstallStarted:             "InstallStarted",
	AccountIDMissing:           "AccountIDMissing",
	APIKeyMissing:              "APIKeyMissing",
	RegionMissing:              "RegionMissing",
	UnableToConnect:            "UnableToConnect",
	UnableToFetchLicenseKey:    "UnableToFetchLicenseKey",
	LicenseKeyFetchedOk:        "LicenseKeyFetchedOk",
	UnableToPostData:           "UnableToPostData",
	UnableToLocatePostedData:   "UnableToLocatePostedData",
	InstallCompleted:           "InstallCompleted",
	InstallCancelled:           "InstallCancelled",
	UnableToOverrideClient:     "UnableToOverrideClient",
	InvalidIngestKey:           "InvalidIngestKey",
	UnableToDiscover:           "UnableToDiscover",
	NrIntegrationPollingErrror: "NrIntegrationPollingErrror",
	OtherError:                 "OtherError",
}

func TryParseEventType(e string) (EventType, bool) {
	switch e {
	case "InstallStarted":
		return EventTypes.InstallStarted, true
	case "AccountIDMissing":
		return EventTypes.AccountIDMissing, true
	case "APIKeyMissing":
		return EventTypes.APIKeyMissing, true
	case "RegionMissing":
		return EventTypes.RegionMissing, true
	case "UnableToConnect":
		return EventTypes.UnableToConnect, true
	case "UnableToFetchLicenseKey":
		return EventTypes.UnableToFetchLicenseKey, true
	case "LicenseKeyFetchedOk":
		return EventTypes.LicenseKeyFetchedOk, true
	case "UnableToPostData":
		return EventTypes.UnableToPostData, true
	case "InstallCompleted":
		return EventTypes.InstallCompleted, true
	case "UnableToOverrideClient":
		return EventTypes.UnableToOverrideClient, true
	case "InvalidIngestKey":
		return EventTypes.InvalidIngestKey, true
	case "UnableToDiscover":
		return EventTypes.UnableToDiscover, true
	case "NrIntegrationPollingErrror":
		return EventTypes.NrIntegrationPollingErrror, true
	}
	return "", false
}

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

type UnsupportedOperatingSystemError struct {
	Err error
}

func (e *UnsupportedOperatingSystemError) Error() string {
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

// nolint: golint
var ErrorFetchingLicenseKey = errors.New("Oops, we're having some difficulties fetching your license key. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic")

type ErrUnalbeToFetchLicenseKey struct {
	Err     error
	Details string
}

func (e *ErrUnalbeToFetchLicenseKey) Error() string {
	return "could not fetch license key"
}

const PaymentRequiredExceptionMessage = `
  Your account has exceeded its plan data limit.
  Take full advantage of New Relic's platform by managing
  your account's plan and payment options at the URL below.`

type ConnectionError struct {
	Err error
}

func (p ConnectionError) Error() string {
	return fmt.Sprintf("Connection Error: %s", p.Err)
}

type DetailError struct {
	EventName EventType
	Details   string
}

func NewDetailError(eventName EventType, details string) *DetailError {
	return &DetailError{
		eventName,
		details,
	}
}

func (e *DetailError) Error() string {
	return e.Details
}
