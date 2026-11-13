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

// versionEntity builds one entity as returned by actor.entitySearch, with its fields carried as tags.
func versionEntity(datasetID string, version int) map[string]interface{} {
	tag := func(key string, values ...string) map[string]interface{} {
		return map[string]interface{}{"key": key, "values": values}
	}
	return map[string]interface{}{
		"guid": fmt.Sprintf("%s-v%d", datasetID, version),
		"name": fmt.Sprintf("%s v%d", datasetID, version),
		"tags": []map[string]interface{}{
			tag("aimDataset", datasetID),
			tag("version", fmt.Sprint(version)),
			tag("recordCount", "10"),
			tag("variables", "input", "context"),
			tag("metadata.createdAt", "1793815928536"),
			tag("metadata.updatedAt", "1793815928713"),
		},
	}
}

func versionSearchBody(entities []map[string]interface{}, nextCursor string) string {
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"entitySearch": map[string]interface{}{
					"results": map[string]interface{}{"entities": entities, "nextCursor": nextCursor},
				},
			},
		},
	})
	return string(body)
}

// capturedRequest is what the test server saw for one entitySearch call.
type capturedRequest struct {
	header http.Header
	body   []byte
}

// stubVersionSearch points newRelicHTTPClient at a test server that serves the given bodies in order,
// one per call, using the real aimTransport. Returns the requests received, in order.
func stubVersionSearch(t *testing.T, status int, bodies ...string) *[]capturedRequest {
	t.Helper()

	var requests []capturedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requests = append(requests, capturedRequest{header: r.Header, body: buf})

		if len(requests) > len(bodies) {
			t.Fatalf("unexpected extra entitySearch call (call #%d)", len(requests))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, writeErr := w.Write([]byte(bodies[len(requests)-1]))
		require.NoError(t, writeErr)
	}))
	t.Cleanup(server.Close)

	tc := testhelpers.NewTestConfig(t, server)
	tc.HTTPTransport = aimTransport{}

	original := newRelicHTTPClient
	newRelicHTTPClient = &newrelic.NewRelic{NerdGraph: nerdgraph.New(tc)}
	t.Cleanup(func() { newRelicHTTPClient = original })

	return &requests
}

func requestVariables(t *testing.T, r capturedRequest) map[string]interface{} {
	t.Helper()
	var req struct {
		Variables map[string]interface{} `json:"variables"`
	}
	require.NoError(t, json.Unmarshal(r.body, &req))
	return req.Variables
}

func TestListDatasetVersions_PagesSortsAndSendsCapabilityHeader(t *testing.T) {
	requests := stubVersionSearch(t, http.StatusOK,
		versionSearchBody([]map[string]interface{}{versionEntity("ds-1", 2)}, "cursor-page-2"),
		versionSearchBody([]map[string]interface{}{versionEntity("ds-1", 1)}, ""),
	)

	got, err := listDatasetVersions("ds-1")
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, datasetVersion{
		ID:          "ds-1-v1",
		Name:        "ds-1 v1",
		AimDataset:  "ds-1",
		Version:     new(int64(1)),
		RecordCount: new(int64(10)),
		Variables:   []string{"input", "context"},
		CreatedAt:   new(int64(1793815928536)),
		UpdatedAt:   new(int64(1793815928713)),
	}, got[0])
	assert.Equal(t, int64(2), *got[1].Version)

	require.Len(t, *requests, 2)
	for _, r := range *requests {
		assert.Equal(t, "AI_MONITORING", r.header.Get("X-Query-Source-Capability-Id"))
	}
	first := requestVariables(t, (*requests)[0])
	assert.Equal(t, "type = 'AIM_DATASET_VERSION' AND tags.aimDataset = 'ds-1'", first["query"])
	assert.Nil(t, first["cursor"])
	assert.Equal(t, "cursor-page-2", requestVariables(t, (*requests)[1])["cursor"])
}

func TestListDatasetVersions_MissingTagsLeftEmptyAndSortLast(t *testing.T) {
	stubVersionSearch(t, http.StatusOK, versionSearchBody([]map[string]interface{}{
		{"guid": "no-tags", "name": "no tags", "tags": []map[string]interface{}{}},
		{"guid": "ds-1-v7", "name": "v7", "tags": []map[string]interface{}{{"key": "version", "values": []string{"7"}}}},
	}, ""))

	got, err := listDatasetVersions("ds-1")
	require.NoError(t, err)
	assert.Equal(t, []datasetVersion{
		{ID: "ds-1-v7", Name: "v7", Version: new(int64(7))},
		{ID: "no-tags", Name: "no tags"},
	}, got)

	out, err := json.Marshal(got[1])
	require.NoError(t, err)
	assert.JSONEq(t, `{"id": "no-tags", "name": "no tags"}`, string(out))
}

func TestListDatasetVersions_NoMatch(t *testing.T) {
	logs := captureLogOutput(t)
	stubVersionSearch(t, http.StatusOK, versionSearchBody([]map[string]interface{}{}, ""))

	got, err := listDatasetVersions("ds-1")
	require.NoError(t, err)
	assert.Equal(t, []datasetVersion{}, got)
	assert.Contains(t, logs.String(), "no versions were found for AI Monitoring dataset ds-1")
}

func TestListDatasetVersions_RejectsQuoteInID(t *testing.T) {
	_, err := listDatasetVersions("ds-1' OR type = 'AIM_DATASET")
	assert.Error(t, err)
}

func TestListDatasetVersions_ErrorPropagates(t *testing.T) {
	stubVersionSearch(t, http.StatusBadRequest, `{"errors":[{"message":"internal error"}]}`)

	got, err := listDatasetVersions("ds-1")
	assert.Error(t, err)
	assert.Nil(t, got)
}
