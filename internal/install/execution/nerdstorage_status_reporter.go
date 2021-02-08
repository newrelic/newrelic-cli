package execution

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
)

const (
	packageID    = "00000000-0000-0000-0000-000000000000"
	collectionID = "openInstallLibrary"
)

// NerdstorageStatusReporter is an implementation of the ExecutionStatusReporter
// interface that reports esecution status into NerdStorage.
type NerdstorageStatusReporter struct {
	client NerdStorageClient
}

// NewNerdStorageStatusReporter returns a new instance of NerdStorageExecutionStatusReporter.
func NewNerdStorageStatusReporter(client NerdStorageClient) *NerdstorageStatusReporter {
	r := NerdstorageStatusReporter{
		client: client,
	}

	return &r
}

// ReportRecipesAvailable reports that recipes are available for installation on
// the underlying host.
func (r NerdstorageStatusReporter) ReportRecipesAvailable(status *StatusRollup, recipes []types.Recipe) error {
	if err := r.writeStatus(status, ""); err != nil {
		return err
	}

	return nil
}

// ReportRecipeAvailable reports that a recipe is available for installation on
// the underlying host.
func (r NerdstorageStatusReporter) ReportRecipeAvailable(status *StatusRollup, recipe types.Recipe) error {
	if err := r.writeStatus(status, ""); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeFailed(status *StatusRollup, event RecipeStatusEvent) error {
	if err := r.writeStatus(status, event.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeInstalling(status *StatusRollup, event RecipeStatusEvent) error {
	if err := r.writeStatus(status, event.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeInstalled(status *StatusRollup, event RecipeStatusEvent) error {
	if err := r.writeStatus(status, event.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeRecommended(status *StatusRollup, event RecipeStatusEvent) error {
	if err := r.writeStatus(status, event.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeSkipped(status *StatusRollup, event RecipeStatusEvent) error {
	if err := r.writeStatus(status, event.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportComplete(status *StatusRollup) error {
	if err := r.writeStatus(status, ""); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) writeStatus(status *StatusRollup, entityGUID string) error {
	i := r.buildExecutionStatusDocument(status)
	_, err := r.client.WriteDocumentWithUserScope(i)
	if err != nil {
		return err
	}

	if entityGUID != "" {
		_, err := r.client.WriteDocumentWithEntityScope(entityGUID, i)
		if err != nil {
			return err
		}
	} else {
		log.Debug("No entity GUID available, skipping entity-scoped status update.")
	}

	return nil
}

func (r NerdstorageStatusReporter) buildExecutionStatusDocument(status *StatusRollup) nerdstorage.WriteDocumentInput {
	return nerdstorage.WriteDocumentInput{
		PackageID:  packageID,
		Collection: collectionID,
		DocumentID: status.DocumentID,
		Document:   status,
	}
}
