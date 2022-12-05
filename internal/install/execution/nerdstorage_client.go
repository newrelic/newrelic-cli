package execution

import (
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdstorage"
)

type NerdStorageClient interface {
	WriteDocumentWithUserScope(nerdstorage.WriteDocumentInput) (interface{}, error)
	WriteDocumentWithEntityScope(string, nerdstorage.WriteDocumentInput) (interface{}, error)
	WriteDocumentWithAccountScope(int, nerdstorage.WriteDocumentInput) (interface{}, error)
}
