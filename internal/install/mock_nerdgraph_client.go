package install

import (
	"context"
	"reflect"
)

type mockNerdGraphClient struct {
	respBody interface{}
}

func newMockNerdGraphClient() *mockNerdGraphClient {
	return &mockNerdGraphClient{
		respBody: struct{}{},
	}
}

func (c *mockNerdGraphClient) QueryWithResponseAndContext(ctx context.Context, query string, variables map[string]interface{}, respBody interface{}) error {
	respBodyPtrValue := reflect.ValueOf(respBody)
	respBodyValue := reflect.Indirect(respBodyPtrValue)
	respBodyValue.Set(reflect.ValueOf(c.respBody))

	return nil
}
