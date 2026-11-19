package aimonitoring

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
)

const datasetVersionSearchQuery = "type = 'AIM_DATASET_VERSION'"

var datasetID string

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

// Nullable fields in the schema are omitted from the output when unset.
type datasetVersion struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	AimDataset  string   `json:"aimDataset"`
	Version     int      `json:"version"`
	Message     *string  `json:"message,omitempty"`
	RecordCount *int     `json:"recordCount,omitempty"`
	Variables   []string `json:"variables,omitempty"`
	Metadata    struct {
		CreatedAt *int64 `json:"createdAt,omitempty"`
		UpdatedAt *int64 `json:"updatedAt,omitempty"`
	} `json:"metadata"`
}

func listDatasetVersions(id string) ([]datasetVersion, error) {
	allVersions, err := fetchAllEntities[*aimonitoring.EntityManagementAimDatasetVersionEntity](datasetVersionSearchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to list dataset versions: %w", err)
	}

	matched := make([]datasetVersion, 0, len(allVersions))
	for _, v := range allVersions {
		if v.AimDataset == id {
			matched = append(matched, toDatasetVersion(v))
		}
	}
	slices.SortFunc(matched, func(a, b datasetVersion) int { return cmp.Compare(a.Version, b.Version) })

	if len(matched) == 0 {
		log.Infof("no versions were found for AI Monitoring dataset %s", id)
	}

	return matched, nil
}

// The client decodes null as the zero value, so zero values are treated as unset.
// ponytail: a real recordCount of 0 is omitted too; needs a *int in client-go to tell them apart.
func toDatasetVersion(v *aimonitoring.EntityManagementAimDatasetVersionEntity) datasetVersion {
	out := datasetVersion{
		ID:         v.ID,
		Name:       v.Name,
		AimDataset: v.AimDataset,
		Version:    v.Version,
		Variables:  v.Variables,
	}
	if v.Message != "" {
		out.Message = &v.Message
	}
	if v.RecordCount != 0 {
		out.RecordCount = &v.RecordCount
	}
	out.Metadata.CreatedAt = epochMillis(time.Time(v.Metadata.CreatedAt))
	out.Metadata.UpdatedAt = epochMillis(time.Time(v.Metadata.UpdatedAt))
	return out
}

func epochMillis(t time.Time) *int64 {
	if t.IsZero() {
		return nil
	}
	ms := t.UnixMilli()
	return &ms
}

func init() {
	Command.AddCommand(cmdDatasetVersion)

	cmdDatasetVersion.AddCommand(cmdDatasetVersionList)
	cmdDatasetVersionList.Flags().StringVar(&datasetID, "dataset-id", "", "the ID of the dataset to list versions for")
	utils.LogIfError(cmdDatasetVersionList.MarkFlagRequired("dataset-id"))
}
