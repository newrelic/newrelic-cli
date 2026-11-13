//go:build unit

package aimonitoring

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdgraph"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

const testDatasetID = "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVHwx"

func TestAIMonitoringDatasetDelete(t *testing.T) {
	assert.Equal(t, "delete", cmdDatasetDelete.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetDelete)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetDelete, []string{"id"})
}

// stubNerdGraph points client.NRClient.NerdGraph at a test server that answers every
// request with the given body, and restores the original client.NRClient when the test
// finishes. Returns the raw request bodies received, in order.
func stubNerdGraph(t *testing.T, statusCode int, body string) *[][]byte {
	t.Helper()

	var requests [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requests = append(requests, buf)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, writeErr := w.Write([]byte(body))
		require.NoError(t, writeErr)
	}))
	t.Cleanup(server.Close)

	tc := testhelpers.NewTestConfig(t, server)

	original := client.NRClient
	client.NRClient = &newrelic.NewRelic{NerdGraph: nerdgraph.New(tc)}
	t.Cleanup(func() { client.NRClient = original })

	return &requests
}

// stubPromptConfirm replaces the interactive prompt with one that returns the given
// answer, and records each message it was asked. Restores the real prompt afterwards.
func stubPromptConfirm(t *testing.T, confirmed bool, err error) *[]string {
	t.Helper()

	var prompts []string
	original := promptConfirm
	promptConfirm = func(msg string) (bool, error) {
		prompts = append(prompts, msg)
		return confirmed, err
	}
	t.Cleanup(func() { promptConfirm = original })

	return &prompts
}

// requestInput pulls the "input" GraphQL variable out of a captured request body.
func requestInput(t *testing.T, body []byte) map[string]interface{} {
	t.Helper()
	var req struct {
		Variables struct {
			Input map[string]interface{} `json:"input"`
		} `json:"variables"`
	}
	require.NoError(t, json.Unmarshal(body, &req))
	return req.Variables.Input
}

func graphQLDataBody(data map[string]interface{}) string {
	body, _ := json.Marshal(map[string]interface{}{"data": data})
	return string(body)
}

func graphQLErrorBody(message string) string {
	body, _ := json.Marshal(map[string]interface{}{
		"errors": []map[string]string{{"message": message}},
	})
	return string(body)
}

func datasetDeleteBody(datasetID string) string {
	return graphQLDataBody(map[string]interface{}{
		"aimDatasetDelete": map[string]string{"datasetId": datasetID},
	})
}

func TestConfirmAction_AssumeYesSkipsPrompt(t *testing.T) {
	prompts := stubPromptConfirm(t, false, nil)

	confirmed, err := confirmAction(true, "Delete?")
	require.NoError(t, err)
	assert.True(t, confirmed)
	assert.Empty(t, *prompts, "--yes should never show the prompt")
}

func TestConfirmAction_InterruptCountsAsDeclined(t *testing.T) {
	stubPromptConfirm(t, false, terminal.InterruptErr)

	confirmed, err := confirmAction(false, "Delete?")
	require.NoError(t, err)
	assert.False(t, confirmed)
}

func TestConfirmAction_PromptErrorPointsAtYesFlag(t *testing.T) {
	stubPromptConfirm(t, false, io.EOF)

	confirmed, err := confirmAction(false, "Delete?")
	require.Error(t, err)
	assert.ErrorIs(t, err, io.EOF)
	assert.Contains(t, err.Error(), "--yes")
	assert.False(t, confirmed)
}

func TestDeleteDataset_Confirmed(t *testing.T) {
	prompts := stubPromptConfirm(t, true, nil)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteBody(testDatasetID))

	got, err := deleteDataset(testDatasetID, false)
	require.NoError(t, err)
	assert.Equal(t, &datasetDeleteResult{DatasetID: testDatasetID}, got)

	require.Len(t, *prompts, 1)
	assert.Contains(t, (*prompts)[0], testDatasetID)

	require.Len(t, *requests, 1)
	assert.Equal(t, map[string]interface{}{"datasetId": testDatasetID}, requestInput(t, (*requests)[0]))
}

func TestDeleteDataset_DeclinedMakesNoRequest(t *testing.T) {
	stubPromptConfirm(t, false, nil)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteBody(testDatasetID))

	got, err := deleteDataset(testDatasetID, false)
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.Empty(t, *requests, "declining must not call the API")
}

func TestDeleteDataset_PromptErrorMakesNoRequest(t *testing.T) {
	stubPromptConfirm(t, false, io.EOF)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteBody(testDatasetID))

	got, err := deleteDataset(testDatasetID, false)
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Empty(t, *requests, "a failed prompt must not call the API")
}

func TestDeleteDataset_AssumeYes(t *testing.T) {
	prompts := stubPromptConfirm(t, false, nil)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteBody(testDatasetID))

	got, err := deleteDataset(testDatasetID, true)
	require.NoError(t, err)
	assert.Equal(t, &datasetDeleteResult{DatasetID: testDatasetID}, got)
	assert.Empty(t, *prompts)
	assert.Len(t, *requests, 1)
}

func TestDeleteDataset_APIErrorPropagates(t *testing.T) {
	stubPromptConfirm(t, true, nil)
	stubNerdGraph(t, http.StatusBadRequest, graphQLErrorBody("dataset not found"))

	got, err := deleteDataset(testDatasetID, false)
	assert.Error(t, err)
	assert.Nil(t, got)
}
