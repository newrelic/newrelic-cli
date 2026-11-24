package aimonitoring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
)

var datasetModes = []aimonitoring.AimDatasetMode{
	aimonitoring.AimDatasetModeTypes.TEXT,
	aimonitoring.AimDatasetModeTypes.SINGLE_TURN,
	aimonitoring.AimDatasetModeTypes.MULTI_TURN,
}

var (
	datasetCreateName        string
	datasetCreateMode        string
	datasetCreateFile        string
	datasetCreateDescription string
)

var cmdDatasetCreateJSON = &cobra.Command{
	Use:   "create-json",
	Short: "Create an AI Monitoring dataset from a JSON file of records",
	Long: `Create an AI Monitoring dataset from a JSON file of records

The create-json command creates a new dataset in the current account from a JSON
file containing an array of records. Each record must have an "input", and may have
"expected", "context" and "additionalFields" (an object):

  [
    {"input": "What is New Relic?", "expected": "An observability platform"}
  ]

For --mode TEXT, "input" is a string. For SINGLE_TURN and MULTI_TURN it is an array
of turns, each an object of template variables:

  [
    {"input": [{"question": "What is APM?"}], "additionalFields": {"difficulty": "easy"}}
  ]
`,
	Example: `newrelic aimonitoring dataset create-json --name "My Dataset" --mode TEXT --file records.json`,
	Args:    cobra.NoArgs,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()

		result, err := createDatasetFromJSON(accountID, datasetCreateName, datasetCreateMode, datasetCreateFile, datasetCreateDescription)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(result))
	},
}

func createDatasetFromJSON(accountID int, name, mode, file, description string) (*aimonitoring.AimDatasetCreateResult, error) {
	if !slices.Contains(datasetModes, aimonitoring.AimDatasetMode(mode)) {
		return nil, fmt.Errorf("--mode must be one of: %v", datasetModes)
	}

	records, err := readDatasetRecords(file)
	if err != nil {
		return nil, err
	}

	input := aimonitoring.AimDatasetCreateInput{
		Name:        name,
		Mode:        aimonitoring.AimDatasetMode(mode),
		Description: description,
		Records:     records,
		Scope: aimonitoring.AimDatasetScopeInput{
			ID:   strconv.Itoa(accountID),
			Type: aimonitoring.AimDatasetScopeTypeTypes.ACCOUNT,
		},
	}

	result, err := client.NRClient.AIMonitoring.AimDatasetCreateWithContext(utils.SignalCtx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create dataset: %w", err)
	}

	return result, nil
}

func readDatasetRecords(file string) ([]aimonitoring.AimDatasetRecordInput, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read records file: %w", err)
	}

	// Unknown fields are rejected so a typo like "expect" fails here instead of being silently dropped.
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var records []aimonitoring.AimDatasetRecordInput
	if err := dec.Decode(&records); err != nil {
		return nil, fmt.Errorf("failed to parse records file %s: expected a JSON array of records: %w", file, err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("records file %s contains no records", file)
	}

	for i := range records {
		if err := encodeRawFields(&records[i]); err != nil {
			return nil, fmt.Errorf("record %d in %s: %w", i, file, err)
		}
	}

	return records, nil
}

// AimDatasetRawRecord is a JSON scalar the NerdGraph gateway only accepts as a string, so
// turns and additionalFields are sent JSON-encoded.
func encodeRawFields(r *aimonitoring.AimDatasetRecordInput) error {
	switch input := r.Input.(type) {
	case nil:
		return fmt.Errorf("missing \"input\"")
	case string:
	case []interface{}:
		if len(input) == 0 {
			return fmt.Errorf("\"input\" has no turns")
		}
		for _, turn := range input {
			if _, ok := turn.(map[string]interface{}); !ok {
				return fmt.Errorf("each turn in \"input\" must be an object")
			}
		}
		r.Input = jsonString(input)
	default:
		return fmt.Errorf("\"input\" must be a string or an array of turns")
	}

	if r.AdditionalFields != nil {
		if _, ok := r.AdditionalFields.(map[string]interface{}); !ok {
			return fmt.Errorf("\"additionalFields\" must be an object")
		}
		r.AdditionalFields = jsonString(r.AdditionalFields)
	}

	return nil
}

func jsonString(v interface{}) string {
	b, _ := json.Marshal(v) // v was just decoded from JSON, so encoding it again cannot fail
	return string(b)
}

func init() {
	cmdDataset.AddCommand(cmdDatasetCreateJSON)
	cmdDatasetCreateJSON.Flags().StringVar(&datasetCreateName, "name", "", "the display name of the dataset")
	cmdDatasetCreateJSON.Flags().StringVar(&datasetCreateMode, "mode", "", "the input mode of the records: TEXT, SINGLE_TURN or MULTI_TURN")
	cmdDatasetCreateJSON.Flags().StringVar(&datasetCreateFile, "file", "", "path to a JSON file containing an array of records")
	cmdDatasetCreateJSON.Flags().StringVar(&datasetCreateDescription, "description", "", "a description of the dataset's purpose")
	utils.LogIfError(cmdDatasetCreateJSON.MarkFlagRequired("name"))
	utils.LogIfError(cmdDatasetCreateJSON.MarkFlagRequired("mode"))
	utils.LogIfError(cmdDatasetCreateJSON.MarkFlagRequired("file"))
}
