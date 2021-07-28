package utils

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type MockHTTPClient struct {
	GetCallCount int
	MockDoFunc   MockHTTPDoFunc
}

type MockHTTPDoFunc func(req *http.Request) (*http.Response, error)

func NewMockHTTPClient(mockDoFunc MockHTTPDoFunc) *MockHTTPClient {
	c := MockHTTPClient{
		MockDoFunc: mockDoFunc,
	}

	return &c
}

func (c *MockHTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	c.GetCallCount++

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.Do(req)
	if err != nil {
		log.Debug(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.MockDoFunc(req)
}

// CreateMockHTTPDoFunc is a helper function to create mock responses for
// the MockHTTPClient. In short, it simulates the http.Client.Do() method.
func CreateMockHTTPDoFunc(mockResponse string, statusCode int, err error) MockHTTPDoFunc {
	return func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: statusCode,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte(mockResponse))),
		}, err
	}
}

func CreateMockGetResponse(response string, err error) func(ctx context.Context, url string) ([]byte, error) {
	getFunc := func(ctx context.Context, url string) ([]byte, error) {
		return ([]byte)(response), err
	}
	return getFunc
}

func CreateMockEmptyGetResponse() func(ctx context.Context, url string) ([]byte, error) {
	return func(ctx context.Context, url string) ([]byte, error) {
		return ([]byte)(""), nil
	}
}
