package execution

import (
	log "github.com/sirupsen/logrus"

	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
)

const (
	packageID    = "00000000-0000-0000-0000-000000000000"
	collectionID = "openInstallLibrary"
)

// NerdstorageStatusReporter is an implementation of the ExecutionStatusReporter
// interface that reports execution status into NerdStorage.
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

func (r NerdstorageStatusReporter) RecipeDetected(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	return nil
}

func (r NerdstorageStatusReporter) RecipesSelected(status *InstallStatus, recipes []types.OpenInstallationRecipe) error {
	return nil
}

// RecipeAvailable reports that a recipe is available for installation on
// the underlying host.
func (r NerdstorageStatusReporter) RecipeAvailable(status *InstallStatus, recipe types.OpenInstallationRecipe) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) RecipeFailed(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) RecipeInstalling(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) RecipeInstalled(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) RecipeRecommended(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) RecipeSkipped(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) InstallStarted(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) InstallComplete(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) InstallCanceled(status *InstallStatus) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) RecipeUnsupported(status *InstallStatus, event RecipeStatusEvent) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) DiscoveryComplete(status *InstallStatus, dm types.DiscoveryManifest) error {
	return r.writeStatus(status)
}

func (r NerdstorageStatusReporter) UpdateRequired(status *InstallStatus) error {
	return nil
}

func (r NerdstorageStatusReporter) writeStatus(status *InstallStatus) error {
	i := r.buildExecutionStatusDocument(status)
	_, err := r.client.WriteDocumentWithUserScope(i)
	if err != nil {
		return err
	}

	for _, g := range status.EntityGUIDs {
		_, err = r.client.WriteDocumentWithEntityScope(g, i)
		if err != nil {
			return err
		}
	}

	if len(status.EntityGUIDs) == 0 {
		log.Debug("no entity GUIDs available, skipping entity-scoped status updates")
	}

	accountID := configAPI.GetActiveProfileAccountID()
	_, err = r.client.WriteDocumentWithAccountScope(accountID, i)
	if err != nil {
		log.Debug("failed to write to account scoped nerd storage")
	}

	return nil
}

func (r NerdstorageStatusReporter) buildExecutionStatusDocument(status *InstallStatus) nerdstorage.WriteDocumentInput {
	return nerdstorage.WriteDocumentInput{
		PackageID:  packageID,
		Collection: collectionID,
		DocumentID: status.DocumentID,
		Document:   status,
	}
}
