package utils

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strings"

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
		escapedUrl := html.EscapeString(url)
		log.WithFields(log.Fields{
			"url":   escapedUrl,
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
		escapedUrl := html.EscapeString(url)
		log.WithFields(log.Fields{
			"url":   escapedUrl,
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

	if err != nil {
		if strings.Contains(err.Error(), "context canceled") {
			return resp, err
		}
	}

	respStatusCode := getResponseCodeString(resp)
	url := req.URL.String()
	escapedUrl := html.EscapeString(url)
	if err != nil || !isResponseSuccess(resp) {
		log.WithFields(log.Fields{
			"method":     req.Method,
			"statusCode": respStatusCode,
			"url":        escapedUrl,
			"error":      err,
		}).Debug("HTTPClient: error performing request")
	}

	if !isResponseSuccess(resp) {
		err = fmt.Errorf("error performing %s request to %s: %s", req.Method, escapedUrl, respStatusCode)
	}

	return resp, err
}

func getResponseCodeString(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	return fmt.Sprintf("%d", resp.StatusCode)
}

// Ensures the response status code falls within the
// status codes that are commonly considered successful.
func isResponseSuccess(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	statusCode := resp.StatusCode

	return statusCode >= http.StatusOK && statusCode <= 299
}
