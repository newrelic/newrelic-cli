package aimonitoring

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

const datasetSearchQuery = `query($cursor: String) {
  actor {
    entityManagement {
      entitySearch(query: "type = 'AIM_DATASET'", cursor: $cursor) {
        entities {
          id
          name
          ... on EntityManagementAimDatasetEntity {
            description
            variables
            latestVersion
            latestVersionId
            mode
            tags {
              key
              values
            }
          }
        }
        nextCursor
      }
    }
  }
}`

var datasetName string

var cmdDataset = &cobra.Command{
	Use:     "dataset",
	Short:   "Interact with New Relic AI Monitoring datasets",
	Example: "newrelic aimonitoring dataset --help",
	Long:    "Interact with New Relic AI Monitoring datasets",
}

var cmdDatasetSearch = &cobra.Command{
	Use:   "search",
	Short: "Search for AI Monitoring datasets",
	Long: `Search for AI Monitoring datasets

The search command lists every AI Monitoring dataset the user has access to,
optionally narrowed down to the datasets whose name contains the given name,
ignoring case.
`,
	Example: `newrelic aimonitoring dataset search --name "My Dataset"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		results, err := searchDatasets(datasetName)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(results))
	},
}

type datasetTag struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

type dataset struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	Variables       []string     `json:"variables"`
	LatestVersion   int          `json:"latestVersion"`
	LatestVersionID string       `json:"latestVersionId"`
	Mode            string       `json:"mode"`
	Tags            []datasetTag `json:"tags"`
}

func searchDatasets(name string) ([]dataset, error) {
	allDatasets, err := fetchAllEntities[dataset](datasetSearchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search datasets: %w", err)
	}

	lowerName := strings.ToLower(name)
	matched := make([]dataset, 0, len(allDatasets))
	for _, d := range allDatasets {
		if strings.Contains(strings.ToLower(d.Name), lowerName) {
			matched = append(matched, d)
		}
	}

	if len(matched) == 0 {
		log.Info("no AI Monitoring datasets were found")
	}

	return matched, nil
}

func init() {
	Command.AddCommand(cmdDataset)

	cmdDataset.AddCommand(cmdDatasetSearch)
	cmdDatasetSearch.Flags().StringVarP(&datasetName, "name", "n", "", "narrow results to datasets whose name contains the given value, ignoring case")
}
