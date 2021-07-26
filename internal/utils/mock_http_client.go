package utils

import (
	"context"
)

type MockHTTPClient struct {
	GetErr       error
	GetCallCount int
	GetVal       []byte
}

func NewMockHTTPClient(mockGetResponse string) *MockHTTPClient {
	c := MockHTTPClient{}
	c.GetVal = []byte(mockGetResponse)
	return &c
}

func (c *MockHTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	c.GetCallCount++
	return c.GetVal, c.GetErr
}
