package aimonitoring

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

const datasetDeleteMutation = `mutation($input: AimDatasetDeleteInput!) {
  aimDatasetDelete(dataset: $input) {
    datasetId
  }
}`

var (
	datasetDeleteID  string
	datasetDeleteYes bool
)

var cmdDatasetDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete an AI Monitoring dataset and all of its versions",
	Long: `Delete an AI Monitoring dataset and all of its versions

The delete command permanently deletes the dataset and every one of its versions.
This cannot be undone. You are asked to confirm before anything is deleted; pass
--yes to skip the prompt, which is required when no terminal is available.
`,
	Example: `newrelic aimonitoring dataset delete --id <datasetId>`,
	Args:    cobra.NoArgs,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := deleteDataset(datasetDeleteID, datasetDeleteYes)
		utils.LogIfFatal(err)

		if result == nil {
			log.Info("delete cancelled")
			return
		}

		utils.LogIfFatal(output.Print(result))
	},
}

type datasetDeleteResult struct {
	DatasetID string `json:"datasetId"`
}

type datasetDeleteResponse struct {
	AimDatasetDelete datasetDeleteResult `json:"aimDatasetDelete"`
}

func deleteDataset(datasetID string, assumeYes bool) (*datasetDeleteResult, error) {
	msg := fmt.Sprintf("Permanently delete dataset %s and all of its versions? This cannot be undone.", datasetID)
	confirmed, err := confirmAction(assumeYes, msg)
	if err != nil || !confirmed {
		return nil, err
	}

	var resp datasetDeleteResponse
	vars := map[string]interface{}{
		"input": map[string]interface{}{"datasetId": datasetID},
	}
	if err := client.NRClient.NerdGraph.QueryWithResponseAndContext(utils.SignalCtx, datasetDeleteMutation, vars, &resp); err != nil {
		return nil, fmt.Errorf("failed to delete dataset: %w", err)
	}

	return &resp.AimDatasetDelete, nil
}

func init() {
	cmdDataset.AddCommand(cmdDatasetDelete)
	cmdDatasetDelete.Flags().StringVar(&datasetDeleteID, "id", "", "the ID of the dataset to delete")
	cmdDatasetDelete.Flags().BoolVarP(&datasetDeleteYes, "yes", "y", false, "delete without asking for confirmation")
	utils.LogIfError(cmdDatasetDelete.MarkFlagRequired("id"))
}
