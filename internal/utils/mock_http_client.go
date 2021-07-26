package utils

import "context"

type MockHTTPClient struct {
	GetErr       error
	GetCallCount int
	GetVal       []byte
}

func NewMockHTTPClient() *MockHTTPClient {
	mockReponse := `{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`

	c := MockHTTPClient{}
	c.GetVal = []byte(mockReponse)
	return &c
}

func (c *MockHTTPClient) Get(ctx context.Context, url string) ([]byte, error) {
	c.GetCallCount++
	return c.GetVal, c.GetErr
}
