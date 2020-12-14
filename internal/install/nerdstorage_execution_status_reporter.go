package install

import "github.com/newrelic/newrelic-client-go/pkg/nerdstorage"

const (
	packageID    = "badfa35a-827d-428d-8f5b-33b836b0e2dd"
	collectionID = "openInstallationStatus"
	documentID   = "status"
)

type nerdstorageExecutionStatusReporter struct {
	client nerdstorageClient
}

func newNerdStorageExecutionStatusReporter(client nerdstorageClient) *nerdstorageExecutionStatusReporter {
	r := nerdstorageExecutionStatusReporter{
		client: client,
	}

	return &r
}

func (r nerdstorageExecutionStatusReporter) reportUserStatus(status executionStatus) error {
	i := createExecutionStatusDocument(status)

	_, err := r.client.WriteDocumentWithUserScope(i)
	if err != nil {
		return err
	}

	return nil
}

func (r nerdstorageExecutionStatusReporter) reportEntityStatus(entityGUID string, status executionStatus) error {
	i := createExecutionStatusDocument(status)

	_, err := r.client.WriteDocumentWithEntityScope(entityGUID, i)
	if err != nil {
		return err
	}

	return nil
}

func createExecutionStatusDocument(status executionStatus) nerdstorage.WriteDocumentInput {
	return nerdstorage.WriteDocumentInput{
		PackageID:  packageID,
		Collection: collectionID,
		DocumentID: documentID,
		Document:   status,
	}
}
