//go:build unit
// +build unit

package fleetcontrol

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// entitySearchPage is one canned response for a single GetEntitySearch call.
type entitySearchPage struct {
	fleetNames []string
	nextCursor string
	statusCode int
	graphQLErr string
}

func (p entitySearchPage) body() string {
	if p.graphQLErr != "" {
		errBody, _ := json.Marshal(map[string]interface{}{
			"errors": []map[string]string{{"message": p.graphQLErr}},
		})
		return string(errBody)
	}

	entities := make([]map[string]string, 0, len(p.fleetNames))
	for i, name := range p.fleetNames {
		entities = append(entities, map[string]string{
			"__typename": "EntityManagementFleetEntity",
			"id":         name + "-id",
			"name":       name,
			"type":       "FLEET",
		})
		_ = i
	}

	respBody, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"entityManagement": map[string]interface{}{
					"entitySearch": map[string]interface{}{
						"entities":   entities,
						"nextCursor": p.nextCursor,
					},
				},
			},
		},
	})
	return string(respBody)
}

// requestCursor pulls the "cursor" GraphQL variable out of a captured request body,
// so tests can assert the CLI actually forwards NextCursor back on the next call.
func requestCursor(t *testing.T, body []byte) interface{} {
	t.Helper()
	var req struct {
		Variables map[string]interface{} `json:"variables"`
	}
	require.NoError(t, json.Unmarshal(body, &req))
	return req.Variables["cursor"]
}

// stubEntitySearch points client.NRClient.FleetControl at a test server that serves
// the given pages in order, one per GetEntitySearch call, and restores the original
// client.NRClient when the test finishes. Returns the raw request bodies received, in order.
func stubEntitySearch(t *testing.T, pages []entitySearchPage) *[][]byte {
	t.Helper()

	var requests [][]byte
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		requests = append(requests, buf)

		if callCount >= len(pages) {
			t.Fatalf("unexpected extra GetEntitySearch call (call #%d), only %d pages configured", callCount+1, len(pages))
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
	client.NRClient = &newrelic.NewRelic{FleetControl: fleetcontrol.New(tc)}
	t.Cleanup(func() { client.NRClient = original })

	return &requests
}

func TestSearchFleetEntities_MatchOnFirstPage(t *testing.T) {
	stubEntitySearch(t, []entitySearchPage{
		{fleetNames: []string{"alpha-fleet", "target-fleet"}},
	})

	got, err := searchFleetEntities(SearchFlags{NameEquals: "target-fleet"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "target-fleet", got[0].Name)
}

// TestSearchFleetEntities_MatchOnSecondPage is the regression test for the pagination bug:
// without looping on NextCursor, a fleet that only exists on page 2+ is never seen and
// --name-equals silently returns empty.
func TestSearchFleetEntities_MatchOnSecondPage(t *testing.T) {
	requests := stubEntitySearch(t, []entitySearchPage{
		{fleetNames: []string{"alpha-fleet", "beta-fleet"}, nextCursor: "cursor-page-2"},
		{fleetNames: []string{"target-fleet", "gamma-fleet"}},
	})

	got, err := searchFleetEntities(SearchFlags{NameEquals: "target-fleet"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "target-fleet", got[0].Name)

	require.Len(t, *requests, 2, "expected exactly two GetEntitySearch calls (one per page)")
	assert.Equal(t, "cursor-page-2", requestCursor(t, (*requests)[1]),
		"second call should forward the NextCursor returned by the first")
}

func TestSearchFleetEntities_NoMatchAcrossAllPages(t *testing.T) {
	stubEntitySearch(t, []entitySearchPage{
		{fleetNames: []string{"alpha-fleet"}, nextCursor: "cursor-page-2"},
		{fleetNames: []string{"beta-fleet"}},
	})

	got, err := searchFleetEntities(SearchFlags{NameEquals: "does-not-exist"})
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestSearchFleetEntities_ErrorOnLaterPagePropagates(t *testing.T) {
	stubEntitySearch(t, []entitySearchPage{
		{fleetNames: []string{"alpha-fleet"}, nextCursor: "cursor-page-2"},
		{graphQLErr: "internal error", statusCode: http.StatusBadRequest},
	})

	got, err := searchFleetEntities(SearchFlags{NameEquals: "alpha-fleet"})
	assert.Error(t, err, "an error on a later page should propagate, not be silently swallowed")
	assert.Nil(t, got, "no partial results should be returned when a later page errors")
}

func TestSearchFleetEntities_NameContainsAcrossPages(t *testing.T) {
	stubEntitySearch(t, []entitySearchPage{
		{fleetNames: []string{"prod-fleet-a"}, nextCursor: "cursor-page-2"},
		{fleetNames: []string{"prod-fleet-b", "staging-fleet"}},
	})

	got, err := searchFleetEntities(SearchFlags{NameContains: "prod-fleet"})
	require.NoError(t, err)
	require.Len(t, got, 2)

	var names []string
	for _, f := range got {
		names = append(names, f.Name)
	}
	assert.ElementsMatch(t, []string{"prod-fleet-a", "prod-fleet-b"}, names)
}
