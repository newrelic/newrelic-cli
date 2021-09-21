//go:build integration
// +build integration

package recipes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFetchRecipes(t *testing.T) {
	f := NewEmbeddedRecipeFetcher()

	r, err := f.FetchRecipes(context.Background())
	require.NoError(t, err)
	require.Greater(t, len(r), 0)
}
