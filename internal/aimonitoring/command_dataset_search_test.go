//go:build unit

package aimonitoring

import (
	"encoding/json"
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

func TestAIMonitoringDataset(t *testing.T) {
	assert.Equal(t, "dataset", cmdDataset.Name())

	testcobra.CheckCobraMetadata(t, cmdDataset)
	testcobra.CheckCobraRequiredFlags(t, cmdDataset, []string{})
}

func TestAIMonitoringDatasetSearch(t *testing.T) {
	assert.Equal(t, "search", cmdDatasetSearch.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetSearch)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetSearch, []string{})
}

// datasetSearchPage is one canned response for a single entitySearch call.
type datasetSearchPage struct {
	datasetNames []string
	nextCursor   string
	statusCode   int
	graphQLErr   string
}

func (p datasetSearchPage) body() string {
	if p.graphQLErr != "" {
		errBody, _ := json.Marshal(map[string]interface{}{
			"errors": []map[string]string{{"message": p.graphQLErr}},
		})
		return string(errBody)
	}

	ents := make([]map[string]interface{}, 0, len(p.datasetNames))
	for _, name := range p.datasetNames {
		ents = append(ents, map[string]interface{}{
			"id":              name + "-id",
			"name":            name,
			"description":     name + " description",
			"variables":       []string{"input"},
			"latestVersion":   2,
			"latestVersionId": name + "-version-id",
			"mode":            "SINGLE_TURN",
			"tags":            []map[string]interface{}{{"key": "env", "values": []string{"prod"}}},
		})
	}

	respBody, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"entityManagement": map[string]interface{}{
					"entitySearch": map[string]interface{}{
						"entities":   ents,
						"nextCursor": p.nextCursor,
					},
				},
			},
		},
	})
	return string(respBody)
}

// datasetRequestCursor pulls the "cursor" GraphQL variable out of a captured request body,
// so tests can assert the CLI actually forwards nextCursor back on the next call.
func datasetRequestCursor(t *testing.T, body []byte) interface{} {
	t.Helper()
	var req struct {
		Variables map[string]interface{} `json:"variables"`
	}
	require.NoError(t, json.Unmarshal(body, &req))
	return req.Variables["cursor"]
}

// stubDatasetSearch points client.NRClient.NerdGraph at a test server that serves
// the given pages in order, one per entitySearch call, and restores the original
// client.NRClient when the test finishes. Returns the raw request bodies received, in order.
func stubDatasetSearch(t *testing.T, pages []datasetSearchPage) *[][]byte {
	t.Helper()

	var requests [][]byte
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requests = append(requests, buf)

		if callCount >= len(pages) {
			t.Fatalf("unexpected extra entitySearch call (call #%d), only %d pages configured", callCount+1, len(pages))
		}
		page := pages[callCount]
		callCount++

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

func TestSearchDatasets_MatchIgnoresCase(t *testing.T) {
	stubDatasetSearch(t, []datasetSearchPage{
		{datasetNames: []string{"Alpha Dataset", "Target Dataset"}},
	})

	got, err := searchDatasets("target dataset")
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, dataset{
		ID:              "Target Dataset-id",
		Name:            "Target Dataset",
		Description:     "Target Dataset description",
		Variables:       []string{"input"},
		LatestVersion:   2,
		LatestVersionID: "Target Dataset-version-id",
		Mode:            "SINGLE_TURN",
		Tags:            []datasetTag{{Key: "env", Values: []string{"prod"}}},
	}, got[0])
}

func TestSearchDatasets_PartialMatchIgnoresCaseAcrossPages(t *testing.T) {
	stubDatasetSearch(t, []datasetSearchPage{
		{datasetNames: []string{"My Blah Dataset", "other"}, nextCursor: "cursor-page-2"},
		{datasetNames: []string{"blah-v2", "unrelated"}},
	})

	got, err := searchDatasets("BLAH")
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "My Blah Dataset", got[0].Name)
	assert.Equal(t, "blah-v2", got[1].Name)
}

// TestSearchDatasets_MatchOnSecondPage guards pagination: without looping on nextCursor,
// a dataset that only exists on page 2+ is never seen and --name silently returns empty.
func TestSearchDatasets_MatchOnSecondPage(t *testing.T) {
	requests := stubDatasetSearch(t, []datasetSearchPage{
		{datasetNames: []string{"alpha", "beta"}, nextCursor: "cursor-page-2"},
		{datasetNames: []string{"target", "gamma"}},
	})

	got, err := searchDatasets("target")
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "target", got[0].Name)

	require.Len(t, *requests, 2, "expected exactly two entitySearch calls (one per page)")
	assert.Nil(t, datasetRequestCursor(t, (*requests)[0]), "first call should not send a cursor")
	assert.Equal(t, "cursor-page-2", datasetRequestCursor(t, (*requests)[1]),
		"second call should forward the nextCursor returned by the first")
}

func TestSearchDatasets_NoNameReturnsAllPages(t *testing.T) {
	stubDatasetSearch(t, []datasetSearchPage{
		{datasetNames: []string{"alpha"}, nextCursor: "cursor-page-2"},
		{datasetNames: []string{"beta"}},
	})

	got, err := searchDatasets("")
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "alpha", got[0].Name)
	assert.Equal(t, "beta", got[1].Name)
}

func TestSearchDatasets_NoMatchAcrossAllPages(t *testing.T) {
	logs := captureLogOutput(t)
	stubDatasetSearch(t, []datasetSearchPage{
		{datasetNames: []string{"alpha"}, nextCursor: "cursor-page-2"},
		{datasetNames: []string{"beta"}},
	})

	got, err := searchDatasets("does-not-exist")
	require.NoError(t, err)
	assert.Equal(t, []dataset{}, got)
	assert.Contains(t, logs.String(), "no AI Monitoring datasets were found")
}

func TestSearchDatasets_ErrorOnLaterPagePropagates(t *testing.T) {
	stubDatasetSearch(t, []datasetSearchPage{
		{datasetNames: []string{"alpha"}, nextCursor: "cursor-page-2"},
		{graphQLErr: "internal error", statusCode: http.StatusBadRequest},
	})

	got, err := searchDatasets("alpha")
	assert.Error(t, err, "an error on a later page should propagate, not be silently swallowed")
	assert.Nil(t, got, "no partial results should be returned when a later page errors")
}
