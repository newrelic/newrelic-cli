// +build unit

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()

	httpUrl = "https://af062943-dc76-45d1-8067-7849cbfe0d98.mock.pstmn.io/v1/status"

	client := NewValidationClient()
	resp, err := client.Get(httpUrl)

	require.NoError(t, err)

	entityGUID := resp.EntityGUID
	require.Equal(t, entityGUID, "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw")
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

func TestHttpClient_HttpError(t *testing.T) {
	t.Parallel()

	badHttpUrl = ""

	client := NewValidationClient()
	resp, err := client.Get(BadHttpUrl)

	require.Error(t, err)
	require.Equal(t, resp)
}
