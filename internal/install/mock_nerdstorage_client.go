package install

import "github.com/newrelic/newrelic-client-go/pkg/nerdstorage"

// nolint:unused,deadcode
type mockNerdstorageClient struct {
	respBody                              interface{}
	userScopeError                        error
	entityScopeError                      error
	writeDocumentWithUserScopeCallCount   int
	writeDocumentWithEntityScopeCallCount int
}

// nolint:unused,deadcode
func newMockNerdstorageClient() *mockNerdstorageClient {
	return &mockNerdstorageClient{
		respBody:         struct{}{},
		userScopeError:   nil,
		entityScopeError: nil,
	}
}

func (c *mockNerdstorageClient) WriteDocumentWithUserScope(nerdstorage.WriteDocumentInput) (interface{}, error) {
	c.writeDocumentWithUserScopeCallCount++
	return c.respBody, c.userScopeError
}

func (c *mockNerdstorageClient) WriteDocumentWithEntityScope(string, nerdstorage.WriteDocumentInput) (interface{}, error) {
	c.writeDocumentWithEntityScopeCallCount++
	return c.respBody, c.entityScopeError
}
