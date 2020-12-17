package install

import (
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
)

type nerdstorageClient interface {
	WriteDocumentWithUserScope(nerdstorage.WriteDocumentInput) (interface{}, error)
	WriteDocumentWithEntityScope(string, nerdstorage.WriteDocumentInput) (interface{}, error)
}
