// +build unit

package entities

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEntitiesGetTags(t *testing.T) {
	assert.Equal(t, "get", cmdTagsGet.Name())

	requiredFlags := []string{"guid"}

	for _, r := range requiredFlags {
		x := cmdTagsGet.Flag(r)
		if x == nil {
			t.Errorf("missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
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

func TestEntitiesAssembleTags(t *testing.T) {
	var scenarios = []struct {
		tags     []string
		expected []entities.Tag
		err      error
	}{
		{
			[]string{"one"},
			[]entities.Tag{},
			fmt.Errorf("tags must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"tag1:value1", "tag1:value2", "tag2:value1"},
			[]entities.Tag{
				{Key: "tag1", Values: []string{"value1", "value2"}}, {Key: "tag2", Values: []string{"value1"}},
			},
			nil,
		},
	}

	for _, s := range scenarios {

		r, e := assembleTags(s.tags)

		assert.ElementsMatch(t, s.expected, r)
		assert.Equal(t, s.err, e)
	}
}

func TestEntitiesAssembleTagValues(t *testing.T) {
	var scenarios = []struct {
		tags     []string
		expected []entities.TagValue
		err      error
	}{
		{
			[]string{"one"},
			[]entities.TagValue{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"incomplete:"},
			[]entities.TagValue{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"tag1:value1", "tag1:value2", "tag2:value1"},
			[]entities.TagValue{
				{Key: "tag1", Value: "value1"},
				{Key: "tag1", Value: "value2"},
				{Key: "tag2", Value: "value1"},
			},
			nil,
		},
	}

	for _, s := range scenarios {
		r, e := assembleTagValues(s.tags)

		assert.ElementsMatch(t, s.expected, r)
		assert.Equal(t, s.err, e)
	}
}

func TestEntitiesAssembleTagValue(t *testing.T) {
	var scenarios = []struct {
		tag      string
		expected entities.TagValue
		err      error
	}{
		{
			"invalidTag",
			entities.TagValue{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			"incompleteTag:",
			entities.TagValue{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			"validKey:validValue",
			entities.TagValue{
				Key:   "validKey",
				Value: "validValue",
			},
			nil,
		},
	}

	for _, s := range scenarios {
		r, e := assembleTagValue(s.tag)

		assert.Equal(t, s.expected, r)
		assert.Equal(t, s.err, e)
	}
}
