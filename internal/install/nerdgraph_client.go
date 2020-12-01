package install

import "context"

type nerdGraphClient interface {
	QueryWithResponseAndContext(context.Context, string, map[string]interface{}, interface{}) error
}
