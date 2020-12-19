package recipes

import "context"

type NerdGraphClient interface {
	QueryWithResponseAndContext(context.Context, string, map[string]interface{}, interface{}) error
}
