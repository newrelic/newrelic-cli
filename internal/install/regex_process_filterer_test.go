// +build unit

package install

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter(t *testing.T) {
	recipes := []recipe{
		{
			ID: "test",
			Metadata: recipeMetadata{
				Name:         "java",
				ProcessMatch: []string{"java"},
			},
		},
	}

	mockRecipeFetcher := newMockRecipeFetcher()
	mockRecipeFetcher.fetchRecipesFunc = func() ([]recipe, error) {
		return recipes, nil
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
