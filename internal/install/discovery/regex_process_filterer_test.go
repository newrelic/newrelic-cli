// +build unit

package discovery

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestFilter(t *testing.T) {
	r := []types.Recipe{
		{
			ID:           "test",
			Name:         "java",
			ProcessMatch: []string{"java"},
		},
	}

	mockRecipeFetcher := recipes.NewMockRecipeFetcher()
	mockRecipeFetcher.FetchRecipesVal = r

	processes := []types.GenericProcess{
		mockProcess{
			name: "java",
		},
		mockProcess{
			name: "somethingElse",
		},
	}

	f := NewRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes)

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 1, len(filtered))
}
