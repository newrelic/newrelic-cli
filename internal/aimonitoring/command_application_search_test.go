//go:build unit

package aimonitoring

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

// stubApplicationSearchAPI points client.NRClient at a test server that answers
// both the NRQL (uniques(entityGuid)) call and the Entities(guids) call made by
// searchApplications, routed by request body content rather than call order.
// It restores the original client.NRClient when the test finishes.
func stubApplicationSearchAPI(t *testing.T, nrqlBody, entitiesBody string) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		w.Header().Set("Content-Type", "application/json")

		switch {
		case strings.Contains(string(buf), "$accountId"):
			_, writeErr := w.Write([]byte(nrqlBody))
			require.NoError(t, writeErr)
		case strings.Contains(string(buf), "$guids"):
			_, writeErr := w.Write([]byte(entitiesBody))
			require.NoError(t, writeErr)
		default:
			t.Fatalf("unexpected request body: %s", buf)
		}
	}))
	t.Cleanup(server.Close)

	tc := testhelpers.NewTestConfig(t, server)

	original := client.NRClient
	client.NRClient = &newrelic.NewRelic{
		Nrdb:     nrdb.New(tc),
		Entities: entities.New(tc),
	}
	t.Cleanup(func() { client.NRClient = original })
}

// nrqlGUIDsBody builds the NerdGraph response for the uniques(entityGuid) query.
func nrqlGUIDsBody(guids []string) string {
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"account": map[string]interface{}{
					"nrql": map[string]interface{}{
						"results": []map[string]interface{}{
							{"guids": guids},
						},
					},
				},
			},
		},
	})
	return string(body)
}

// nrqlErrorBody builds a GraphQL error response.
func nrqlErrorBody(message string) string {
	body, _ := json.Marshal(map[string]interface{}{
		"errors": []map[string]string{{"message": message}},
	})
	return string(body)
}

// entityStub describes one ApmApplicationEntity to serve from the Entities(guids) call.
type entityStub struct {
	name string
	guid string
	tags map[string]string
}

func entitiesResultBody(stubs []entityStub) string {
	ents := make([]map[string]interface{}, 0, len(stubs))
	for _, s := range stubs {
		tags := make([]map[string]interface{}, 0, len(s.tags))
		for k, v := range s.tags {
			tags = append(tags, map[string]interface{}{"key": k, "values": []string{v}})
		}
		ents = append(ents, map[string]interface{}{
			"__typename": "ApmApplicationEntity",
			"name":       s.name,
			"guid":       s.guid,
			"tags":       tags,
		})
	}

	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"actor": map[string]interface{}{
				"entities": ents,
			},
		},
	})
	return string(body)
}

// captureLogOutput redirects logrus's package-level output to a buffer for the
// duration of the test, so assertions can check what was logged.
func captureLogOutput(t *testing.T) *bytes.Buffer {
	t.Helper()

	var buf bytes.Buffer
	original := log.StandardLogger().Out
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(original) })

	return &buf
}

func TestSearchApplications_NoResultsFound(t *testing.T) {
	logs := captureLogOutput(t)
	stubApplicationSearchAPI(t, nrqlGUIDsBody([]string{}), "")

	got, err := searchApplications(12345, "7 days ago", "", nil)
	require.NoError(t, err)
	assert.Equal(t, []applicationSearchResult{}, got)
	assert.Contains(t, logs.String(), "no applications reporting AI Monitoring telemetry were found")
}

func TestSearchApplications_MatchFound(t *testing.T) {
	stubApplicationSearchAPI(t,
		nrqlGUIDsBody([]string{"guid-checkout"}),
		entitiesResultBody([]entityStub{
			{name: "checkout-service", guid: "guid-checkout", tags: map[string]string{"aiEnabledApp": "true"}},
		}),
	)

	got, err := searchApplications(12345, "7 days ago", "", []string{"aiEnabledApp:true"})
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, applicationSearchResult{Name: "checkout-service", GUID: "guid-checkout"}, got[0])
}

func TestSearchApplications_NoTagMatch(t *testing.T) {
	logs := captureLogOutput(t)
	stubApplicationSearchAPI(t,
		nrqlGUIDsBody([]string{"guid-checkout"}),
		entitiesResultBody([]entityStub{
			{name: "checkout-service", guid: "guid-checkout", tags: map[string]string{"aiEnabledApp": "false"}},
		}),
	)

	got, err := searchApplications(12345, "7 days ago", "", []string{"aiEnabledApp:true"})
	require.NoError(t, err)
	assert.Equal(t, []applicationSearchResult{}, got)
	assert.Contains(t, logs.String(), "no applications matched the given tags")
}

func TestSearchApplications_ResultsTruncated(t *testing.T) {
	logs := captureLogOutput(t)

	guids := make([]string, maxApplicationSearchResults)
	stubs := make([]entityStub, maxApplicationSearchResults)
	for i := range guids {
		guid := fmt.Sprintf("guid-%d", i)
		guids[i] = guid
		stubs[i] = entityStub{name: fmt.Sprintf("app-%d", i), guid: guid}
	}

	stubApplicationSearchAPI(t, nrqlGUIDsBody(guids), entitiesResultBody(stubs))

	got, err := searchApplications(12345, "7 days ago", "", nil)
	require.NoError(t, err)
	assert.Len(t, got, maxApplicationSearchResults)
	assert.Contains(t, logs.String(), "results were truncated")
}

func TestSearchApplications_NRQLErrorPropagates(t *testing.T) {
	stubApplicationSearchAPI(t, nrqlErrorBody("internal error"), "")

	got, err := searchApplications(12345, "7 days ago", "", nil)
	assert.Error(t, err)
	assert.Nil(t, got)
}

func TestSearchApplications_InvalidTagFormat(t *testing.T) {
	// No HTTP calls should be made: an invalid --tags value is rejected before
	// any query is built.
	got, err := searchApplications(12345, "7 days ago", "", []string{"not-a-key-value-pair"})
	assert.Error(t, err)
	assert.Nil(t, got)
}
