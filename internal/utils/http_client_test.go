// +build unit

package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()

	response := `{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`
	statusCode := 200

	ts := testServer(response, statusCode)
	defer ts.Close()

	ctx := context.Background()

	c := NewValidationClient()
	data, err := c.Get(ctx, ts.URL)

	require.Equal(t, string(data), response)
	require.Equal(t, err, nil)
}

func TestHttpClientDataNotAvailable(t *testing.T) {
	t.Parallel()

	response := `{}`
	statusCode := 204

	ts := testServer(response, statusCode)
	defer ts.Close()

	ctx := context.Background()

	c := NewValidationClient()
	data, err := c.Get(ctx, ts.URL)

	require.Equal(t, data, nil)
	require.NotEqual(t, err, nil)
}

func TestHttpClientNotFound(t *testing.T) {
	t.Parallel()

	response := ``
	statusCode := 404

	ts := testServer(response, statusCode)
	defer ts.Close()

	ctx := context.Background()

	c := NewValidationClient()
	data, err := c.Get(ctx, ts.URL)

	require.Equal(t, data, nil)
	require.NotEqual(t, err, nil)
}

func TestHttpClientConnectionRefused(t *testing.T) {
	t.Parallel()

	response := ``
	statusCode := 500

	ts := testServer(response, statusCode)
	defer ts.Close()

	ctx := context.Background()

	c := NewValidationClient()
	data, err := c.Get(ctx, ts.URL)

	require.Equal(t, data, nil)
	require.NotEqual(t, err, nil)
}

func TestHttpClientInternalServerError(t *testing.T) {
	t.Parallel()

	response := ``
	statusCode := 500

	ts := testServer(response, statusCode)
	defer ts.Close()

	ctx := context.Background()

	c := NewValidationClient()
	data, err := c.Get(ctx, ts.URL)

	require.Equal(t, data, nil)
	require.NotEqual(t, err, nil)
}

func TestHttpClientTimeout(t *testing.T) {
	t.Parallel()

	response := ``
	statusCode := 500

	ts := testServer(response, statusCode)
	defer ts.Close()

	ctx := context.Background()

	c := NewValidationClient()
	data, err := c.Get(ctx, ts.URL)

	require.Equal(t, data, nil)
	require.NotEqual(t, err, nil)
}

func testServer(response string, statusCode int) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write([]byte(response))
	}))

	return ts
}
