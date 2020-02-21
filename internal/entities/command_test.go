package entities

import (
	"fmt"
	"testing"

	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestEntitiesCommand(t *testing.T) {
	assert.NotEmptyf(t, Command.Use, "Need to set Command.%s on Command %s", "Use", Command.CalledAs())
	assert.NotEmptyf(t, Command.Short, "Need to set Command.%s on Command %s", "Short", Command.CalledAs())

	for _, c := range Command.Commands() {
		assert.NotEmptyf(t, c.Use, "Need to set Command.%s on Command %s", "Use", c.CommandPath())
		assert.NotEmptyf(t, c.Short, "Need to set Command.%s on Command %s", "Short", c.CommandPath())
		assert.NotEmptyf(t, c.Long, "Need to set Command.%s on Command %s", "Long", c.CommandPath())
		assert.NotEmptyf(t, c.Example, "Need to set Command.%s on Command %s", "Example", c.CommandPath())
	}
}

func TestEntitiesSearch(t *testing.T) {
	command := entitiesSearch

	assert.Equal(t, "search", command.Name())
	assert.True(t, command.HasFlags())
}

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

	requiredFlags := []string{"guid", "tag"}

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

func TestEntitiesCreateTags(t *testing.T) {
	cur := *entitiesCreateTags
	assert.Equal(t, "create-tags", cur.Name())

	requiredFlags := []string{"guid", "tag"}

	for _, r := range requiredFlags {
		x := cur.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
			continue
		}

		assert.Equal(t, []string{"true"}, x.Annotations[cobra.BashCompOneRequiredFlag])
	}
}

func TestEntitiesReplaceTags(t *testing.T) {
	cur := *entitiesReplaceTags
	assert.Equal(t, "replace-tags", cur.Name())

	requiredFlags := []string{"guid", "tag"}

	for _, r := range requiredFlags {
		x := cur.Flag(r)
		if x == nil {
			t.Errorf("Missing required flag: %s\n", r)
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
			fmt.Errorf("Tags must be specified as colon separated key:value pairs"),
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
			fmt.Errorf("Tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"incomplete:"},
			[]entities.TagValue{},
			fmt.Errorf("Tag values must be specified as colon separated key:value pairs"),
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
			fmt.Errorf("Tag values must be specified as colon separated key:value pairs"),
		},
		{
			"incompleteTag:",
			entities.TagValue{},
			fmt.Errorf("Tag values must be specified as colon separated key:value pairs"),
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
