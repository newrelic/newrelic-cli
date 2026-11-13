//go:build unit

package aimonitoring

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

const testVersionID = "MTE4NTkwOTh8TkdFUHxBSU1fREFUQVNFVF9WRVJTSU9OfDE"

func TestAIMonitoringDatasetDeleteVersion(t *testing.T) {
	assert.Equal(t, "delete-version", cmdDatasetDeleteVersion.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetDeleteVersion)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetDeleteVersion, []string{"id", "version-id"})
}

func datasetDeleteVersionBody(versionID string) string {
	return graphQLDataBody(map[string]interface{}{
		"aimDatasetDeleteVersion": map[string]string{"versionId": versionID},
	})
}

func TestDeleteDatasetVersion_Confirmed(t *testing.T) {
	prompts := stubPromptConfirm(t, true, nil)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteVersionBody(testVersionID))

	got, err := deleteDatasetVersion(testDatasetID, testVersionID, false)
	require.NoError(t, err)
	assert.Equal(t, &datasetDeleteVersionResult{VersionID: testVersionID}, got)

	require.Len(t, *prompts, 1)
	assert.Contains(t, (*prompts)[0], testVersionID)
	assert.Contains(t, (*prompts)[0], testDatasetID)
	assert.Contains(t, (*prompts)[0], "only version left")

	require.Len(t, *requests, 1)
	assert.Equal(t, map[string]interface{}{"datasetId": testDatasetID, "versionId": testVersionID}, requestInput(t, (*requests)[0]))
}

func TestDeleteDatasetVersion_DeclinedMakesNoRequest(t *testing.T) {
	stubPromptConfirm(t, false, nil)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteVersionBody(testVersionID))

	got, err := deleteDatasetVersion(testDatasetID, testVersionID, false)
	require.NoError(t, err)
	assert.Nil(t, got)
	assert.Empty(t, *requests, "declining must not call the API")
}

func TestDeleteDatasetVersion_PromptErrorMakesNoRequest(t *testing.T) {
	stubPromptConfirm(t, false, errors.New("not a terminal"))
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteVersionBody(testVersionID))

	got, err := deleteDatasetVersion(testDatasetID, testVersionID, false)
	assert.Error(t, err)
	assert.Nil(t, got)
	assert.Empty(t, *requests, "a failed prompt must not call the API")
}

func TestDeleteDatasetVersion_AssumeYes(t *testing.T) {
	prompts := stubPromptConfirm(t, false, nil)
	requests := stubNerdGraph(t, http.StatusOK, datasetDeleteVersionBody(testVersionID))

	got, err := deleteDatasetVersion(testDatasetID, testVersionID, true)
	require.NoError(t, err)
	assert.Equal(t, &datasetDeleteVersionResult{VersionID: testVersionID}, got)
	assert.Empty(t, *prompts)
	assert.Len(t, *requests, 1)
}

func TestDeleteDatasetVersion_APIErrorPropagates(t *testing.T) {
	stubPromptConfirm(t, true, nil)
	stubNerdGraph(t, http.StatusBadRequest, graphQLErrorBody("dataset version not found"))

	got, err := deleteDatasetVersion(testDatasetID, testVersionID, false)
	assert.Error(t, err)
	assert.Nil(t, got)
}
