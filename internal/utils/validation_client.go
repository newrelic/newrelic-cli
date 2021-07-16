package utils

import (
	"context"
	"io/ioutil"
	"net/http"
)

// Should ValidationClient live in it's own file (e.g. validation_client.go)
type ValidationClient struct{}

// TODO: move to interfaces.go
type HTTPClient interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

func NewValidationClient() *ValidationClient {
	return &ValidationClient{}
}

func (c *ValidationClient) Get(ctx context.Context, url string) ([]byte, error) {
	httpClient := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, _ := ioutil.ReadAll(resp.Body)

	return data, nil
}
