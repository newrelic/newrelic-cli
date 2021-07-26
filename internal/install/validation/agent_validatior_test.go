package validation

import (
	"context"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/stretchr/testify/require"
)

const infraAgentValidationURL = "http://localhost:18003/v1/status/entity"

func TestAgentValidator(t *testing.T) {
	c := utils.NewMockHTTPClient()
	av := NewAgentValidator(c)

	ctx, cancel := context.WithCancel(utils.SignalCtx)
	defer cancel()

	guid, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, guid, "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw")
	require.Equal(t, err, nil)
}

// func TestHttpClientDataNotAvailable(t *testing.T) {
// 	t.Parallel()

// 	response := `{}`
// 	statusCode := 204

// 	ts := testServer(response, statusCode)
// 	defer ts.Close()

// 	ctx := context.Background()

// 	c := NewValidationClient()
// 	data, err := c.Get(ctx, ts.URL)

// 	require.Equal(t, data, nil)
// 	require.NotEqual(t, err, nil)
// }

// func TestHttpClientNotFound(t *testing.T) {
// 	t.Parallel()

// 	response := ``
// 	statusCode := 404

// 	ts := testServer(response, statusCode)
// 	defer ts.Close()

// 	ctx := context.Background()

// 	c := NewValidationClient()
// 	data, err := c.Get(ctx, ts.URL)

// 	require.Equal(t, data, nil)
// 	require.NotEqual(t, err, nil)
// }

// func TestHttpClientConnectionRefused(t *testing.T) {
// 	t.Parallel()

// 	response := ``
// 	statusCode := 500

// 	ts := testServer(response, statusCode)
// 	defer ts.Close()

// 	ctx := context.Background()

// 	c := NewValidationClient()
// 	data, err := c.Get(ctx, ts.URL)

// 	require.Equal(t, data, nil)
// 	require.NotEqual(t, err, nil)
// }

// func TestHttpClientInternalServerError(t *testing.T) {
// 	t.Parallel()

// 	response := ``
// 	statusCode := 500

// 	ts := testServer(response, statusCode)
// 	defer ts.Close()

// 	ctx := context.Background()

// 	c := NewValidationClient()
// 	data, err := c.Get(ctx, ts.URL)

// 	require.Equal(t, data, nil)
// 	require.NotEqual(t, err, nil)
// }

// func TestHttpClientTimeout(t *testing.T) {
// 	t.Parallel()

// 	response := ``
// 	statusCode := 500

// 	ts := testServer(response, statusCode)
// 	defer ts.Close()

// 	ctx := context.Background()

// 	c := NewValidationClient()
// 	data, err := c.Get(ctx, ts.URL)

// 	require.Equal(t, data, nil)
// 	require.NotEqual(t, err, nil)
// }

// func testServer(response string, statusCode int) *httptest.Server {
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(statusCode)
// 		_, _ = w.Write([]byte(response))
// 	}))

// 	return ts
// }
