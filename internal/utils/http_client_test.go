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

	c.Get()

	expected := `{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`

	require.Equal(t, c.GetCallCount, 1)
	require.Equal(t, c.GetVar, expected)
	require.Equal(t, c.GetErr, nil)
}
