package execution

import (
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
)

type NerdStorageClient interface {
	WriteDocumentWithUserScope(nerdstorage.WriteDocumentInput) (interface{}, error)
	WriteDocumentWithEntityScope(string, nerdstorage.WriteDocumentInput) (interface{}, error)
}
