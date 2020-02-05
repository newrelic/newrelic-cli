package entities

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEntitiesDescribeTags(t *testing.T) {
	assert.Equal(t, "describe-tags", entitiesDescribeTags.Name())

	requiredFlags := []string{"guid"}

	for _, r := range requiredFlags {
		x := entitiesDescribeTags.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestEntitiesDeleteTags(t *testing.T) {
	assert.Equal(t, "delete-tags", entitiesDeleteTags.Name())

	requiredFlags := []string{"guid", "tags"}

	for _, r := range requiredFlags {
		x := entitiesDeleteTags.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestEntitiesDeleteTagValues(t *testing.T) {
	assert.Equal(t, "delete-tag-values", entitiesDeleteTagValues.Name())

	requiredFlags := []string{"guid", "value"}

	for _, r := range requiredFlags {
		x := entitiesDeleteTagValues.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}
