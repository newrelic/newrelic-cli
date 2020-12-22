package execution

import (
	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
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
	status *StatusRollup
}

// NewNerdStorageStatusReporter returns a new instance of NerdStorageExecutionStatusReporter.
func NewNerdStorageStatusReporter(client NerdStorageClient) *NerdstorageStatusReporter {
	rollup := NewStatusRollup()
	r := NerdstorageStatusReporter{
		client: client,
		status: &rollup,
	}

	return &r
}

// ReportRecipesAvailable reports that recipes are available for installation on
// the underlying host.
func (r NerdstorageStatusReporter) ReportRecipesAvailable(recipes []types.Recipe) error {
	r.status.withAvailableRecipes(recipes)
	if err := r.writeStatus(""); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeFailed(e RecipeStatusEvent) error {
	r.status.withRecipeEvent(e, StatusTypes.FAILED)
	if err := r.writeStatus(e.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeInstalled(e RecipeStatusEvent) error {
	r.status.withRecipeEvent(e, StatusTypes.INSTALLED)
	if err := r.writeStatus(e.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportRecipeSkipped(e RecipeStatusEvent) error {
	r.status.withRecipeEvent(e, StatusTypes.SKIPPED)
	if err := r.writeStatus(e.EntityGUID); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) ReportComplete() error {
	r.status.Complete = true
	r.status.Timestamp = utils.GetTimestamp()
	if err := r.writeStatus(""); err != nil {
		return err
	}

	return nil
}

func (r NerdstorageStatusReporter) writeStatus(entityGUID string) error {
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

func (r NerdstorageStatusReporter) buildExecutionStatusDocument() nerdstorage.WriteDocumentInput {
	return nerdstorage.WriteDocumentInput{
		PackageID:  packageID,
		Collection: collectionID,
		DocumentID: r.status.DocumentID,
		Document:   r.status,
	}
}
