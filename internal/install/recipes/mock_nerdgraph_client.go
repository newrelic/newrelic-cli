package recipes

import (
	"context"
	"reflect"
)

type mockNerdGraphClient struct {
	RespBody interface{}
}

func NewMockNerdGraphClient() *mockNerdGraphClient {
	return &mockNerdGraphClient{
		RespBody: struct{}{},
	}
}

func (c *mockNerdGraphClient) QueryWithResponseAndContext(ctx context.Context, query string, variables map[string]interface{}, respBody interface{}) error {
	respBodyPtrValue := reflect.ValueOf(respBody)
	respBodyValue := reflect.Indirect(respBodyPtrValue)
	respBodyValue.Set(reflect.ValueOf(c.RespBody))

	return nil
}
