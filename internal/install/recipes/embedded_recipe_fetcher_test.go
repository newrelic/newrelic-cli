//go:build unit
// +build unit

package recipes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ShouldTrimVersion(t *testing.T) {
	f := NewEmbeddedRecipeFetcher()

	version := f.getLibraryVersion([]byte("v0.1.2.3"))
	require.Equal(t, version, "0.1.2.3")
}

func Test_ShouldNotTrimMissingVersion(t *testing.T) {
	f := NewEmbeddedRecipeFetcher()

	version := f.getLibraryVersion([]byte("0.4.5.6"))
	require.Equal(t, version, "0.4.5.6")
}
