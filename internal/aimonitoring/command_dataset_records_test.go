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
	"github.com/newrelic/newrelic-client-go/v2/pkg/aimonitoring"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdgraph"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestAIMonitoringDatasetRecords(t *testing.T) {
	assert.Equal(t, "records", cmdDatasetRecords.Name())

	testcobra.CheckCobraMetadata(t, cmdDatasetRecords)
	testcobra.CheckCobraRequiredFlags(t, cmdDatasetRecords, []string{"id", "version-id"})
}

// recordJSON builds one record as returned by aimDataset.records.
func recordJSON(id int) map[string]interface{} {
	return map[string]interface{}{
		"id":               id,
		"input":            map[string]interface{}{"question": "q"},
		"expected":         "answer",
		"context":          "ctx",
		"additionalFields": map[string]interface{}{"source": "csv"},
	}
}

type recordsPage struct {
	records    []map[string]interface{}
	nextCursor string
	statusCode int
}

func (p recordsPage) body() string {
	if p.statusCode != 0 {
		return `{"errors":[{"message":"internal error"}]}`
	}
	records := p.records
	if records == nil {
		records = []map[string]interface{}{}
	}
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"account": map[string]interface{}{
					"aimDataset": map[string]interface{}{
						"records": map[string]interface{}{"records": records, "nextCursor": p.nextCursor, "totalCount": 99},
					},
				},
			},
		},
	})
	return string(body)
}

// stubRecords points client.NRClient at a test server that serves the given pages in order,
// one per records call. Returns the raw request bodies received, in order.
func stubRecords(t *testing.T, pages ...recordsPage) *[][]byte {
	t.Helper()

	var requests [][]byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requests = append(requests, buf)

		if len(requests) > len(pages) {
			t.Fatalf("unexpected extra records call (call #%d)", len(requests))
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
	client.NRClient = &newrelic.NewRelic{NerdGraph: nerdgraph.New(tc), AIMonitoring: aimonitoring.New(tc)}
	t.Cleanup(func() { client.NRClient = original })

	return &requests
}

func TestGetDatasetRecords_FollowsCursorAcrossPages(t *testing.T) {
	requests := stubRecords(t,
		recordsPage{records: []map[string]interface{}{recordJSON(1), recordJSON(2)}, nextCursor: "cursor-page-2"},
		recordsPage{records: []map[string]interface{}{recordJSON(3)}},
	)

	got, err := getDatasetRecords(12345, "ds-1", "ver-1", nil, 0)
	require.NoError(t, err)
	require.Len(t, got, 3)

	assert.Equal(t, aimonitoring.AimDatasetRecord{
		ID:               1,
		Input:            map[string]interface{}{"question": "q"},
		Expected:         "answer",
		Context:          "ctx",
		AdditionalFields: map[string]interface{}{"source": "csv"},
	}, got[0])
	assert.Equal(t, 3, got[2].ID)

	require.Len(t, *requests, 2)
	first := (*requests)[0]
	assert.Equal(t, float64(12345), requestVariable(t, first, "accountID"))
	assert.Equal(t, "ds-1", requestVariable(t, first, "aimDatasetId"))
	assert.Equal(t, "ver-1", requestVariable(t, first, "versionId"))
	assert.Nil(t, requestVariable(t, first, "cursor"))
	assert.Nil(t, requestVariable(t, first, "recordIds"))
	assert.Equal(t, "cursor-page-2", requestVariable(t, (*requests)[1], "cursor"))
}

func TestGetDatasetRecords_SendsRecordIDs(t *testing.T) {
	requests := stubRecords(t, recordsPage{records: []map[string]interface{}{recordJSON(4), recordJSON(7)}})

	got, err := getDatasetRecords(12345, "ds-1", "ver-1", []int{4, 7}, 0)
	require.NoError(t, err)
	require.Len(t, got, 2)

	require.Len(t, *requests, 1)
	assert.Equal(t, []interface{}{float64(4), float64(7)}, requestVariable(t, (*requests)[0], "recordIds"))
}

func TestGetDatasetRecords_LimitStopsPagingAndTrims(t *testing.T) {
	requests := stubRecords(t,
		recordsPage{records: []map[string]interface{}{recordJSON(1), recordJSON(2)}, nextCursor: "cursor-page-2"},
		recordsPage{records: []map[string]interface{}{recordJSON(3), recordJSON(4)}, nextCursor: "cursor-page-3"},
	)

	got, err := getDatasetRecords(12345, "ds-1", "ver-1", nil, 3)
	require.NoError(t, err)
	require.Len(t, got, 3)
	assert.Equal(t, 3, got[2].ID)
	assert.Len(t, *requests, 2)
}

func TestGetDatasetRecords_LimitWithinFirstPage(t *testing.T) {
	requests := stubRecords(t,
		recordsPage{records: []map[string]interface{}{recordJSON(1), recordJSON(2)}, nextCursor: "cursor-page-2"},
	)

	got, err := getDatasetRecords(12345, "ds-1", "ver-1", nil, 1)
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, 1, got[0].ID)
	assert.Len(t, *requests, 1)
}

func TestGetDatasetRecords_NoRecords(t *testing.T) {
	logs := captureLogOutput(t)
	stubRecords(t, recordsPage{})

	got, err := getDatasetRecords(12345, "ds-1", "ver-1", nil, 0)
	require.NoError(t, err)
	assert.Equal(t, []aimonitoring.AimDatasetRecord{}, got)
	assert.Contains(t, logs.String(), "no records were found for AI Monitoring dataset ds-1 version ver-1")
}

func TestGetDatasetRecords_ErrorOnLaterPagePropagates(t *testing.T) {
	stubRecords(t,
		recordsPage{records: []map[string]interface{}{recordJSON(1)}, nextCursor: "cursor-page-2"},
		recordsPage{statusCode: http.StatusBadRequest},
	)

	got, err := getDatasetRecords(12345, "ds-1", "ver-1", nil, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get dataset records")
	assert.Nil(t, got)
}
