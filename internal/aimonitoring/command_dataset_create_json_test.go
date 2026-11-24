//go:build unit

package aimonitoring

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
)

func TestAIMonitoringDatasetCreateJSON(t *testing.T) {
	assert.Equal(t, "create-json", cmdDatasetCreateJSON.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetCreateJSON)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetCreateJSON, []string{"name", "mode", "file"})
}

func writeRecordsFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "records.json")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestCreateDatasetFromJSON(t *testing.T) {
	file := writeRecordsFile(t, `[{"input": "hi", "expected": "hello", "additionalFields": {"k": 1}}]`)
	requests := stubNerdGraph(t, http.StatusOK, graphQLDataBody(map[string]interface{}{
		"aimDatasetCreate": map[string]string{"datasetId": testDatasetID, "versionId": "v1"},
	}))

	got, err := createDatasetFromJSON(123, "My Dataset", "TEXT", file, "desc")
	require.NoError(t, err)
	assert.Equal(t, &aimonitoring.AimDatasetCreateResult{DatasetId: testDatasetID, VersionId: "v1"}, got)

	require.Len(t, *requests, 1)
	assert.Equal(t, map[string]interface{}{
		"name":        "My Dataset",
		"mode":        "TEXT",
		"description": "desc",
		"scope":       map[string]interface{}{"id": "123", "type": "ACCOUNT"},
		"records": []interface{}{map[string]interface{}{
			"input":            "hi",
			"expected":         "hello",
			"additionalFields": `{"k":1}`,
		}},
	}, requestInput(t, (*requests)[0], "dataset"))
}

func TestCreateDatasetFromJSON_TurnsAreSentAsJSONStrings(t *testing.T) {
	file := writeRecordsFile(t, `[{"input": [{"question": "What is APM?"}]}]`)
	requests := stubNerdGraph(t, http.StatusOK, graphQLDataBody(map[string]interface{}{
		"aimDatasetCreate": map[string]string{"datasetId": testDatasetID, "versionId": "v1"},
	}))

	_, err := createDatasetFromJSON(123, "n", "SINGLE_TURN", file, "")
	require.NoError(t, err)

	records := requestInput(t, (*requests)[0], "dataset")["records"].([]interface{})
	assert.Equal(t, `[{"question":"What is APM?"}]`, records[0].(map[string]interface{})["input"])
}

func TestCreateDatasetFromJSON_InvalidInputMakesNoRequest(t *testing.T) {
	cases := map[string]struct {
		mode    string
		content string
		errText string
	}{
		"bad mode":            {"text", `[{"input": "hi"}]`, "--mode"},
		"not an array":        {"TEXT", `{"input": "hi"}`, "JSON array"},
		"unknown field":       {"TEXT", `[{"input": "hi", "expect": "x"}]`, "expect"},
		"empty array":         {"TEXT", `[]`, "no records"},
		"non-string expected": {"TEXT", `[{"input": "hi", "expected": 5}]`, "expected"},
		"non-string context":  {"TEXT", `[{"input": "hi", "context": {}}]`, "context"},
		"missing input":       {"TEXT", `[{"expected": "x"}]`, "record 0"},
		"object input":        {"SINGLE_TURN", `[{"input": {"q": "x"}}]`, "array of turns"},
		"no turns":            {"SINGLE_TURN", `[{"input": []}]`, "no turns"},
		"non-object turn":     {"SINGLE_TURN", `[{"input": ["x"]}]`, "each turn"},
		"string additional":   {"TEXT", `[{"input": "hi", "additionalFields": "x"}]`, "additionalFields"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			requests := stubNerdGraph(t, http.StatusOK, "{}")

			_, err := createDatasetFromJSON(123, "n", tc.mode, writeRecordsFile(t, tc.content), "")
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errText)
			assert.Empty(t, *requests)
		})
	}
}

func TestCreateDatasetFromJSON_APIErrorPropagates(t *testing.T) {
	stubNerdGraph(t, http.StatusBadRequest, graphQLErrorBody("invalid input"))

	got, err := createDatasetFromJSON(123, "n", "TEXT", writeRecordsFile(t, `[{"input": "hi"}]`), "")
	assert.Error(t, err)
	assert.Nil(t, got)
}
