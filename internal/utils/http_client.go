package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type HTTPClientInterface interface {
	Get(ctx context.Context, url string) ([]byte, error)
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewHTTPClient(apiKey string) *HTTPClient {
	return &HTTPClient{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}
}

func (c *HTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Debug(err.Error())
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
		log.Debug(err.Error())
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Api-Key", c.apiKey)
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
