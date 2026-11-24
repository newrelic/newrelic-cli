package aimonitoring

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
)

const maxRecordPages = 1000

var (
	datasetRecordsID        string
	datasetRecordsVersionID string
	datasetRecordsIDs       []int
	datasetRecordsLimit     int
)

var cmdDatasetRecords = &cobra.Command{
	Use:   "records",
	Short: "Get the records of an AI Monitoring dataset version",
	Long: `Get the records of an AI Monitoring dataset version

The records command returns the records of the given dataset version, optionally
narrowed down to the given record IDs. Use "newrelic aimonitoring dataset-version list"
to find a version ID. By default every record is returned; pass --limit to cap the
number of records.
`,
	Example: `newrelic aimonitoring dataset records --id <datasetId> --version-id <versionId> --limit 10`,
	Args:    cobra.NoArgs,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()

		results, err := getDatasetRecords(accountID, datasetRecordsID, datasetRecordsVersionID, datasetRecordsIDs, datasetRecordsLimit)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(results))
	},
}

func getDatasetRecords(accountID int, datasetID string, versionID string, recordIDs []int, limit int) ([]aimonitoring.AimDatasetRecord, error) {
	records := []aimonitoring.AimDatasetRecord{}

	var cursor *string
	for range maxRecordPages {
		page, err := client.NRClient.AIMonitoring.GetRecordsWithContext(utils.SignalCtx, accountID, datasetID, cursor, recordIDs, versionID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dataset records: %w", err)
		}

		records = append(records, page.Records...)

		if page.NextCursor == "" || (limit > 0 && len(records) >= limit) {
			break
		}
		cursor = &page.NextCursor
	}

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	if len(records) == 0 {
		log.Infof("no records were found for AI Monitoring dataset %s version %s", datasetID, versionID)
	}

	return records, nil
}

func init() {
	cmdDataset.AddCommand(cmdDatasetRecords)
	cmdDatasetRecords.Flags().StringVar(&datasetRecordsID, "id", "", "the ID of the dataset to get records from")
	cmdDatasetRecords.Flags().StringVar(&datasetRecordsVersionID, "version-id", "", "the ID of the dataset version to get records from")
	cmdDatasetRecords.Flags().IntSliceVar(&datasetRecordsIDs, "record-ids", nil, "only return the records with these IDs, comma-separated")
	cmdDatasetRecords.Flags().IntVarP(&datasetRecordsLimit, "limit", "l", 0, "the maximum number of records to return; 0 returns every record")
	utils.LogIfError(cmdDatasetRecords.MarkFlagRequired("id"))
	utils.LogIfError(cmdDatasetRecords.MarkFlagRequired("version-id"))
}
