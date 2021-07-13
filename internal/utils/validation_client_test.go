// +build unit

package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()

	// httpURL := "https://af062943-dc76-45d1-8067-7849cbfe0d98.mock.pstmn.io/v1/status"

	c := NewTestAPIClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOk)
		_, _ = w.Write([]byte(`{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`))
	}))

	client := NewValidationClient()
	ctx := context.Background()
	resp, err := client.Get(ctx, httpURL)

	require.NoError(t, err)

	entityGUID := resp.GUID
	require.Equal(t, entityGUID, "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw")
}

func NewTestAPIClient(t *testing.T, handler http.Handler) ValidationClient {
	ts := httptest.NewServer(handler)
	tc := mock.NewTestConfig(t, ts)
	return NewClient(tc)
}

/*
Pass:
200 Status (ready, entity is connected to NR)
204 Status (ready, but data not available)

Error:
404 Status -> Alert us.
Connection Refused
500 Status
Timeout

*/

// func TestHttpClient_HttpError(t *testing.T) {
// 	t.Parallel()

// 	badHttpURL := ""

// 	client := NewValidationClient()
// 	ctx := context.Background()
// 	_, err := client.Get(ctx, badHttpURL)

// 	require.Error(t, err)
// 	// require.Equal(t, resp)
// }
