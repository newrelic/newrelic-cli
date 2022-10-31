//go:build unit
// +build unit

package entities

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEntitiesGetTags(t *testing.T) {
	assert.Equal(t, "get", cmdTagsGet.Name())
}

func TestEntitiesDeleteTags(t *testing.T) {
	assert.Equal(t, "delete", cmdTagsDelete.Name())

	requiredFlags := []string{"guid", "tag"}

	for _, r := range requiredFlags {
		x := cmdTagsDelete.Flag(r)
		if x == nil {
			t.Errorf("missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestEntitiesDeleteTagValues(t *testing.T) {
	assert.Equal(t, "delete-values", cmdTagsDeleteValues.Name())

	requiredFlags := []string{"guid", "value"}

	for _, r := range requiredFlags {
		x := cmdTagsDeleteValues.Flag(r)
		if x == nil {
			t.Errorf("missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestEntitiesCreateTags(t *testing.T) {
	cur := *cmdTagsCreate
	assert.Equal(t, "create", cur.Name())

	requiredFlags := []string{"guid", "tag"}

	for _, r := range requiredFlags {
		x := cur.Flag(r)
		if x == nil {
			t.Errorf("missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestEntitiesReplaceTags(t *testing.T) {
	cur := *cmdTagsReplace
	assert.Equal(t, "replace", cur.Name())

	requiredFlags := []string{"guid", "tag"}

	for _, r := range requiredFlags {
		x := cur.Flag(r)
		if x == nil {
			t.Errorf("missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}
