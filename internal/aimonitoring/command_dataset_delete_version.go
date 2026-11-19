package aimonitoring

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
)

var (
	versionDeleteDatasetID string
	versionDeleteVersionID string
	versionDeleteYes       bool
)

var cmdDatasetDeleteVersion = &cobra.Command{
	Use:   "delete-version",
	Short: "Delete a single version of an AI Monitoring dataset",
	Long: `Delete a single version of an AI Monitoring dataset

The delete-version command permanently deletes one version of a dataset. If it is
the only version left, the dataset itself is deleted too. This cannot be undone.
You are asked to confirm before anything is deleted; pass --yes to skip the prompt,
which is required when no terminal is available.
`,
	Example: `newrelic aimonitoring dataset delete-version --id <datasetId> --version-id <versionId>`,
	Args:    cobra.NoArgs,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := deleteDatasetVersion(versionDeleteDatasetID, versionDeleteVersionID, versionDeleteYes)
		utils.LogIfFatal(err)

		if result == nil {
			log.Info("delete cancelled")
			return
		}

		utils.LogIfFatal(output.Print(result))
	},
}

func deleteDatasetVersion(datasetID, versionID string, assumeYes bool) (*aimonitoring.AimDatasetDeleteVersionResult, error) {
	msg := fmt.Sprintf("Permanently delete version %s of dataset %s? If it is the only version left, the dataset is deleted too. This cannot be undone.", versionID, datasetID)
	confirmed, err := confirmAction(assumeYes, msg)
	if err != nil || !confirmed {
		return nil, err
	}

	input := aimonitoring.AimDatasetDeleteVersionInput{DatasetId: datasetID, VersionId: versionID}
	result, err := client.NRClient.AIMonitoring.AimDatasetDeleteVersionWithContext(utils.SignalCtx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to delete dataset version: %w", err)
	}

	return result, nil
}

func init() {
	cmdDataset.AddCommand(cmdDatasetDeleteVersion)
	cmdDatasetDeleteVersion.Flags().StringVar(&versionDeleteDatasetID, "id", "", "the ID of the dataset the version belongs to")
	cmdDatasetDeleteVersion.Flags().StringVar(&versionDeleteVersionID, "version-id", "", "the ID of the dataset version to delete")
	cmdDatasetDeleteVersion.Flags().BoolVarP(&versionDeleteYes, "yes", "y", false, "delete without asking for confirmation")
	utils.LogIfError(cmdDatasetDeleteVersion.MarkFlagRequired("id"))
	utils.LogIfError(cmdDatasetDeleteVersion.MarkFlagRequired("version-id"))
}
