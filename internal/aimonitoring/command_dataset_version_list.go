package aimonitoring

import (
	"cmp"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
)

const datasetVersionSearchQuery = `query($query: String!, $cursor: String) {
  actor {
    entitySearch(query: $query) {
      results(cursor: $cursor) {
        entities {
          guid
          name
          tags {
            key
            values
          }
        }
        nextCursor
      }
    }
  }
}`

var datasetID string

var newRelicHTTPClient *newrelic.NewRelic

// Add the X-Query-Source-Capability-Id: AI_MONITORING header to every request. Without it,
// entity search silently returns no AIM_DATASET_VERSION entities.
type aimTransport struct{}

func (aimTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("X-Query-Source-Capability-Id", "AI_MONITORING")
	return http.DefaultTransport.RoundTrip(r)
}

func getNewRelicHTTPClient() (*newrelic.NewRelic, error) {
	if newRelicHTTPClient == nil {
		c, err := client.NewClient(configAPI.GetActiveProfileName(), newrelic.ConfigHTTPTransport(aimTransport{}))
		if err != nil {
			return nil, err
		}
		newRelicHTTPClient = c
	}
	return newRelicHTTPClient, nil
}

var cmdDatasetVersion = &cobra.Command{
	Use:     "dataset-version",
	Short:   "Interact with New Relic AI Monitoring dataset versions",
	Example: "newrelic aimonitoring dataset-version --help",
	Long:    "Interact with New Relic AI Monitoring dataset versions",
}

var cmdDatasetVersionList = &cobra.Command{
	Use:   "list",
	Short: "List the versions of an AI Monitoring dataset",
	Long: `List the versions of an AI Monitoring dataset

The list command returns every version of the dataset with the given ID,
ordered by version number. Use "newrelic aimonitoring dataset search" to
find a dataset ID.
`,
	Example: `newrelic aimonitoring dataset-version list --dataset-id "MTIzfEFJTXxEQVRBU0VUfDE"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		results, err := listDatasetVersions(datasetID)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(results))
	},
}

// Fields other than ID and Name come from tags, none of which are guaranteed
// to be present, so missing ones are omitted from the output.
type datasetVersion struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	AimDataset  string   `json:"aimDataset,omitempty"`
	Version     *int64   `json:"version,omitempty"`
	RecordCount *int64   `json:"recordCount,omitempty"`
	Variables   []string `json:"variables,omitempty"`
	CreatedAt   *int64   `json:"createdAt,omitempty"`
	UpdatedAt   *int64   `json:"updatedAt,omitempty"`
}

type datasetVersionSearchResponse struct {
	Actor struct {
		EntitySearch struct {
			Results struct {
				Entities []struct {
					GUID string       `json:"guid"`
					Name string       `json:"name"`
					Tags []datasetTag `json:"tags"`
				} `json:"entities"`
				NextCursor string `json:"nextCursor"`
			} `json:"results"`
		} `json:"entitySearch"`
	} `json:"actor"`
}

func listDatasetVersions(id string) ([]datasetVersion, error) {
	if id == "" || strings.ContainsAny(id, `'\`) {
		return nil, fmt.Errorf("invalid dataset ID %q", id)
	}

	nrClient, err := getNewRelicHTTPClient()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"query":  fmt.Sprintf("type = 'AIM_DATASET_VERSION' AND tags.aimDataset = '%s'", id),
		"cursor": nil,
	}

	versions := []datasetVersion{}
	for range maxEntitySearchPages {
		var resp datasetVersionSearchResponse
		if err := nrClient.NerdGraph.QueryWithResponseAndContext(utils.SignalCtx, datasetVersionSearchQuery, vars, &resp); err != nil {
			return nil, fmt.Errorf("failed to list dataset versions: %w", err)
		}

		results := resp.Actor.EntitySearch.Results
		for _, e := range results.Entities {
			tags := entityTags{}
			for _, t := range e.Tags {
				tags[t.Key] = t.Values
			}

			versions = append(versions, datasetVersion{
				ID:          e.GUID,
				Name:        e.Name,
				AimDataset:  tags.first("aimDataset"),
				Version:     tags.number("version"),
				RecordCount: tags.number("recordCount"),
				Variables:   tags["variables"],
				CreatedAt:   tags.number("metadata.createdAt"),
				UpdatedAt:   tags.number("metadata.updatedAt"),
			})
		}

		if results.NextCursor == "" {
			break
		}
		vars["cursor"] = results.NextCursor
	}
	// Versions without a version number sort last.
	sortKey := func(v datasetVersion) int64 {
		if v.Version == nil {
			return math.MaxInt64
		}
		return *v.Version
	}
	slices.SortFunc(versions, func(a, b datasetVersion) int { return cmp.Compare(sortKey(a), sortKey(b)) })

	if len(versions) == 0 {
		log.Infof("no versions were found for AI Monitoring dataset %s", id)
	}

	return versions, nil
}

type entityTags map[string][]string

func (t entityTags) first(key string) string {
	if len(t[key]) == 0 {
		return ""
	}
	return t[key][0]
}

// number returns nil for a missing or non-numeric value.
func (t entityTags) number(key string) *int64 {
	n, err := strconv.ParseInt(t.first(key), 10, 64)
	if err != nil {
		return nil
	}
	return &n
}

func init() {
	Command.AddCommand(cmdDatasetVersion)

	cmdDatasetVersion.AddCommand(cmdDatasetVersionList)
	cmdDatasetVersionList.Flags().StringVar(&datasetID, "dataset-id", "", "the ID of the dataset to list versions for")
	utils.LogIfError(cmdDatasetVersionList.MarkFlagRequired("dataset-id"))
}
