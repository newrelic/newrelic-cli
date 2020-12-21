package execution

import "github.com/newrelic/newrelic-client-go/pkg/nerdstorage"

type MockNerdStorageClient struct {
	WriteDocumentWithUserScopeVal         interface{}
	WriteDocumentWithEntityScopeVal       interface{}
	WriteDocumentWithUserScopeErr         error
	WriteDocumentWithEntityScopeErr       error
	writeDocumentWithUserScopeCallCount   int
	writeDocumentWithEntityScopeCallCount int
}

func NewMockNerdStorageClient() *MockNerdStorageClient {
	return &MockNerdStorageClient{
		WriteDocumentWithUserScopeVal:   struct{}{},
		WriteDocumentWithEntityScopeVal: struct{}{},
		WriteDocumentWithUserScopeErr:   nil,
		WriteDocumentWithEntityScopeErr: nil,
	}
}

func (c *MockNerdStorageClient) WriteDocumentWithUserScope(nerdstorage.WriteDocumentInput) (interface{}, error) {
	c.writeDocumentWithUserScopeCallCount++
	return c.WriteDocumentWithUserScopeVal, c.WriteDocumentWithUserScopeErr
}

func (c *MockNerdStorageClient) WriteDocumentWithEntityScope(string, nerdstorage.WriteDocumentInput) (interface{}, error) {
	c.writeDocumentWithEntityScopeCallCount++
	return c.WriteDocumentWithEntityScopeVal, c.WriteDocumentWithEntityScopeErr
}
