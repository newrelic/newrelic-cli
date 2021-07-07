// +build unit

package utils

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

func TestHttpClient(t *testing.T) {
	t.Parallel()

	client := NewValidationClient()
	resp, err := client.Get("https://af062943-dc76-45d1-8067-7849cbfe0d98.mock.pstmn.io/v1/status")

	require.NoError(t, err)

	// TODO: add more test cases
}
