package install

import "github.com/newrelic/newrelic-client-go/pkg/nerdstorage"

type mockNerdstorageClient struct {
	respBody interface{}
	err      error
}

func newMockNerdstorageClient() *mockNerdstorageClient {
	return &mockNerdstorageClient{
		respBody: struct{}{},
		err:      nil,
	}
}

func (c *mockNerdstorageClient) WriteDocumentWithUserScope(nerdstorage.WriteDocumentInput) (interface{}, error) {
	return c.respBody, c.err
}

func (c *mockNerdstorageClient) WriteDocumentWithEntityScope(string, nerdstorage.WriteDocumentInput) (interface{}, error) {
	return c.respBody, c.err
}
