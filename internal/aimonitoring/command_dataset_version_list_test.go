//go:build unit

package aimonitoring

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdgraph"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestAIMonitoringDatasetVersion(t *testing.T) {
	assert.Equal(t, "dataset-version", cmdDatasetVersion.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetVersion)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetVersion, []string{})
}

func TestAIMonitoringDatasetVersionList(t *testing.T) {
	assert.Equal(t, "list", cmdDatasetVersionList.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetVersionList)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetVersionList, []string{"dataset-id"})
}

// versionEntity builds one AIM_DATASET_VERSION entity as returned by entityManagement.entitySearch.
func versionEntity(datasetID string, version int) map[string]interface{} {
	return map[string]interface{}{
		"id":          fmt.Sprintf("%s-v%d", datasetID, version),
		"name":        fmt.Sprintf("%s v%d", datasetID, version),
		"aimDataset":  datasetID,
		"version":     version,
		"message":     "commit",
		"recordCount": 10,
		"variables":   []string{"input", "context"},
		"metadata":    map[string]interface{}{"createdAt": 1793815928536, "updatedAt": 1793815928713},
	}
}

type versionSearchPage struct {
	entities   []map[string]interface{}
	nextCursor string
	statusCode int
}

func (p versionSearchPage) body() string {
	if p.statusCode != 0 {
		return `{"errors":[{"message":"internal error"}]}`
	}
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"entityManagement": map[string]interface{}{
					"entitySearch": map[string]interface{}{"entities": p.entities, "nextCursor": p.nextCursor},
				},
			},
		},
	})
	return string(body)
}

// stubVersionSearch points client.NRClient at a test server that serves the given pages in order,
// one per entitySearch call. Returns the raw request bodies received, in order.
func stubVersionSearch(t *testing.T, pages ...versionSearchPage) *[][]byte {
	t.Helper()

	var requests [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requests = append(requests, buf)

		if len(requests) > len(pages) {
			t.Fatalf("unexpected extra entitySearch call (call #%d)", len(requests))
		}
		page := pages[len(requests)-1]
		status := page.statusCode
		if status == 0 {
			status = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, writeErr := w.Write([]byte(page.body()))
		require.NoError(t, writeErr)
	}))
	t.Cleanup(server.Close)

	tc := testhelpers.NewTestConfig(t, server)

	original := client.NRClient
	client.NRClient = &newrelic.NewRelic{NerdGraph: nerdgraph.New(tc)}
	t.Cleanup(func() { client.NRClient = original })

	return &requests
}

func TestListDatasetVersions_FiltersByDatasetAcrossPagesAndSorts(t *testing.T) {
	requests := stubVersionSearch(t,
		versionSearchPage{entities: []map[string]interface{}{versionEntity("ds-1", 2), versionEntity("ds-2", 1)}, nextCursor: "cursor-page-2"},
		versionSearchPage{entities: []map[string]interface{}{versionEntity("ds-1", 1)}},
	)

	got, err := listDatasetVersions("ds-1")
	require.NoError(t, err)
	require.Len(t, got, 2)

	want := datasetVersion{
		ID:          "ds-1-v1",
		Name:        "ds-1 v1",
		AimDataset:  "ds-1",
		Version:     1,
		Message:     new("commit"),
		RecordCount: new(10),
		Variables:   []string{"input", "context"},
	}
	want.Metadata.CreatedAt = new(int64(1793815928536))
	want.Metadata.UpdatedAt = new(int64(1793815928713))
	assert.Equal(t, want, got[0])
	assert.Equal(t, 2, got[1].Version)

	require.Len(t, *requests, 2)
	assert.Equal(t, "cursor-page-2", datasetRequestCursor(t, (*requests)[1]))
}

func TestListDatasetVersions_OmitsUnsetFields(t *testing.T) {
	stubVersionSearch(t, versionSearchPage{entities: []map[string]interface{}{{
		"id": "ds-1-v7", "name": "v7", "aimDataset": "ds-1", "version": 7,
		"message": nil, "recordCount": nil, "variables": nil, "metadata": map[string]interface{}{},
	}}})

	got, err := listDatasetVersions("ds-1")
	require.NoError(t, err)
	require.Len(t, got, 1)

	out, err := json.Marshal(got[0])
	require.NoError(t, err)
	assert.JSONEq(t, `{"id": "ds-1-v7", "name": "v7", "aimDataset": "ds-1", "version": 7, "metadata": {}}`, string(out))
}

func TestListDatasetVersions_NoMatch(t *testing.T) {
	logs := captureLogOutput(t)
	stubVersionSearch(t, versionSearchPage{entities: []map[string]interface{}{versionEntity("ds-2", 1)}})

	got, err := listDatasetVersions("ds-1")
	require.NoError(t, err)
	assert.Equal(t, []datasetVersion{}, got)
	assert.Contains(t, logs.String(), "no versions were found for AI Monitoring dataset ds-1")
}

func TestListDatasetVersions_ErrorOnLaterPagePropagates(t *testing.T) {
	stubVersionSearch(t,
		versionSearchPage{entities: []map[string]interface{}{versionEntity("ds-1", 1)}, nextCursor: "cursor-page-2"},
		versionSearchPage{statusCode: http.StatusBadRequest},
	)

	got, err := listDatasetVersions("ds-1")
	assert.Error(t, err)
	assert.Nil(t, got)
}
