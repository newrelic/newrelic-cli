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

const maxDatasetSearchPages = 1000

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

type datasetSearchResponse struct {
	Actor struct {
		EntityManagement struct {
			EntitySearch struct {
				Entities   []dataset `json:"entities"`
				NextCursor string    `json:"nextCursor"`
			} `json:"entitySearch"`
		} `json:"entityManagement"`
	} `json:"actor"`
}

// searchDatasets fetches every page of type='AIM_DATASET' entities, following
// nextCursor until the API reports no more pages, then applies the name filter
// against the full accumulated set. The API only supports filtering by type, so name
// filtering has to happen client-side - which only works correctly if every page has
// been fetched first. It returns an empty (non-nil) slice, logging why, when nothing
// is found.
func searchDatasets(name string) ([]dataset, error) {
	var allDatasets []dataset

	var cursor *string
	for range maxDatasetSearchPages {
		var resp datasetSearchResponse
		vars := map[string]interface{}{"cursor": cursor}
		if err := client.NRClient.NerdGraph.QueryWithResponseAndContext(utils.SignalCtx, datasetSearchQuery, vars, &resp); err != nil {
			return nil, fmt.Errorf("failed to search datasets: %w", err)
		}

		search := resp.Actor.EntityManagement.EntitySearch
		allDatasets = append(allDatasets, search.Entities...)

		if search.NextCursor == "" {
			break
		}
		nextCursor := search.NextCursor
		cursor = &nextCursor
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
