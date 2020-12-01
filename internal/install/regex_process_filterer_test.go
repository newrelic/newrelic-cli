// +build unit

package install

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	recipeFilters := []recipeFilter{
		{
			ID: "test",
			Metadata: recipeFilterMetadata{
				Name:         "java",
				ProcessMatch: []string{"java"},
			},
		},
	}

	mockRecipeFetcher := newMockRecipeFetcher()
	mockRecipeFetcher.fetchFiltersFunc = func() ([]recipeFilter, error) {
		return recipeFilters, nil
	}

	processes := []genericProcess{
		mockProcess{
			name: "java",
		},
		mockProcess{
			name: "somethingElse",
		},
	}

	f := newRegexProcessFilterer(mockRecipeFetcher)
	filtered, err := f.filter(context.Background(), processes)

	require.NoError(t, err)
	require.NotNil(t, filtered)
	require.NotEmpty(t, filtered)
	require.Equal(t, 1, len(filtered))
}
