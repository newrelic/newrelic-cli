package recipes

import (
	"context"
	"reflect"
)

type MockNerdGraphClient struct {
	RespBody interface{}
}

func NewMockNerdGraphClient() *MockNerdGraphClient {
	return &MockNerdGraphClient{
		RespBody: struct{}{},
	}
}

func (c *MockNerdGraphClient) QueryWithResponseAndContext(ctx context.Context, query string, variables map[string]interface{}, respBody interface{}) error {
	respBodyPtrValue := reflect.ValueOf(respBody)
	respBodyValue := reflect.Indirect(respBodyPtrValue)
	respBodyValue.Set(reflect.ValueOf(c.RespBody))

	return nil
}
