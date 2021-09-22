package utils

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HTTPClientInterface interface {
	Get(ctx context.Context, url string) ([]byte, error)
	Post(ctx context.Context, url string, requestBody []byte) ([]byte, error)
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewHTTPClient(apiKey string) HTTPClientInterface {
	return &HTTPClient{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}

func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"url":   url,
			"error": err.Error(),
		}).Debug("HTTPClient: error creating new GET request")
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	return data, nil
}

func (c *HTTPClient) Post(ctx context.Context, url string, requestBody []byte) ([]byte, error) {
	rBody := bytes.NewBuffer(requestBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, rBody)
	if err != nil {
		log.WithFields(log.Fields{
			"url":   url,
			"error": err.Error(),
		}).Debug("HTTPClient: error creating new POST request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Api-Key", c.apiKey)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)

	if err != nil || !isResponseSuccess(resp) {
		log.WithFields(log.Fields{
			"method":     req.Method,
			"statusCode": resp.StatusCode,
			"url":        req.URL.String(),
			"error":      err,
		}).Debug("HTTPClient: error performing request")
	}

	if !isResponseSuccess(resp) {
		err = fmt.Errorf("error performing %s request to %s: %d", req.Method, req.URL.String(), resp.StatusCode)
	}

	return resp, err
}

// Ensures the response status code falls within the
// status codes that are commonly considered successful.
func isResponseSuccess(resp *http.Response) bool {
	statusCode := resp.StatusCode

	return statusCode >= http.StatusOK && statusCode <= 299
}
