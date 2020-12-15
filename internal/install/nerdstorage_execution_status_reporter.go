package install

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
)

const (
	packageID    = "badfa35a-827d-428d-8f5b-33b836b0e2dd"
	collectionID = "openInstallLibrary"
)

type nerdstorageExecutionStatusReporter struct {
	client          nerdstorageClient
	executionStatus executionStatus
}

func newNerdStorageExecutionStatusReporter(client nerdstorageClient) *nerdstorageExecutionStatusReporter {
	r := nerdstorageExecutionStatusReporter{
		client:          client,
		executionStatus: newExecutionStatus(),
	}

	return &r
}

func (r nerdstorageExecutionStatusReporter) reportRecipesAvailable(recipes []recipe) error {
	r.executionStatus.withAvailableRecipes(recipes)
	if err := r.writeStatus(""); err != nil {
		return err
	}

	return nil
}

func (r nerdstorageExecutionStatusReporter) writeStatus(entityGUID string) error {
	i := r.buildExecutionStatusDocument()
	_, err := r.client.WriteDocumentWithUserScope(i)
	if err != nil {
		return err
	}

	if entityGUID != "" {
		log.Debug("No entity GUID available, skipping entity-scoped status update.")
		_, err := r.client.WriteDocumentWithEntityScope(entityGUID, i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r nerdstorageExecutionStatusReporter) reportRecipeFailed(e recipeStatusEvent) error {
	r.executionStatus.withRecipeEvent(e, executionStatusRecipeStatusTypes.FAILED)
	if err := r.writeStatus(e.entityGUID); err != nil {
		return err
	}

	return nil
}

func (r nerdstorageExecutionStatusReporter) reportRecipeInstalled(e recipeStatusEvent) error {
	r.executionStatus.withRecipeEvent(e, executionStatusRecipeStatusTypes.INSTALLED)
	if err := r.writeStatus(e.entityGUID); err != nil {
		return err
	}

	return nil
}

func (r nerdstorageExecutionStatusReporter) buildExecutionStatusDocument() nerdstorage.WriteDocumentInput {
	return nerdstorage.WriteDocumentInput{
		PackageID:  packageID,
		Collection: collectionID,
		DocumentID: r.executionStatus.DocumentID,
		Document:   r.executionStatus,
	}
}
