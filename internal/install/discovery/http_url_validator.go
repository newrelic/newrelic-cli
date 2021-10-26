package discovery

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// HttpURLValidator is a struct that implements the URLValidator interface.
type HttpUrlValidator struct{}

// NewUrlValidator returns a new instance of HttpUrlValidator.
func NewHttpUrlValidator() *HttpUrlValidator {
	return &HttpUrlValidator{}
}

// This function will return true if an HTTP GET using the endpoint returns a status code of 200.
// It will return false otherwise.
// It will time out after .25 seconds.
func (v *HttpUrlValidator) Validate(ctx context.Context, endpoint string) bool {
	protocol := "http://"
	address := "169.254.169.254"
	endpoint = "/latest/meta-data/instance-id"
	url := protocol + address + endpoint

	client := http.Client{
		Timeout: 250 * time.Millisecond,
	}

	resp, err := client.Get(url)

	if err != nil {
		log.Debugf("Error validating URL: %s", err)
		return false
	}

	if resp.StatusCode != 200 {
		log.Debugf("Error validating URL, response code: %s", resp.Status)
		return false
	}

	return true
}
