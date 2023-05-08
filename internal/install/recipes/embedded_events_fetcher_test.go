//go:build unit
// +build unit

package recipes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldGetWriteKey(t *testing.T) {
	_, err := NewEmbeddedEventsFetcher().GetWriteKey()

	require.NotNil(t, err)
}
