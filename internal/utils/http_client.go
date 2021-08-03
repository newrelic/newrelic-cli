package utils

import (
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
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		httpClient: &http.Client{},
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

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.httpClient.Do(req)
}
