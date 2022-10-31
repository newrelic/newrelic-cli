//go:build unit
// +build unit

package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func TestEntitiesAssembleTagValues(t *testing.T) {
	var scenarios = []struct {
		tags     []string
		expected []entities.EntitySearchQueryBuilderTag
		err      error
	}{
		{
			[]string{"one"},
			[]entities.EntitySearchQueryBuilderTag{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"incomplete:"},
			[]entities.EntitySearchQueryBuilderTag{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"tag1:value1", "tag1:value2", "tag2:value1"},
			[]entities.EntitySearchQueryBuilderTag{
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

func TestEntitiesAssembleTagsInput(t *testing.T) {
	var scenarios = []struct {
		tags     []string
		expected []entities.TaggingTagInput
		err      error
	}{
		{
			[]string{"one"},
			[]entities.TaggingTagInput{},
			fmt.Errorf("tags must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"tag1:value1", "tag1:value2", "tag2:value1"},
			[]entities.TaggingTagInput{
				{Key: "tag1", Values: []string{"value1", "value2"}}, {Key: "tag2", Values: []string{"value1"}},
			},
			nil,
		},
	}

	for _, s := range scenarios {

		r, e := AssembleTagsInput(s.tags)

		assert.ElementsMatch(t, s.expected, r)
		assert.Equal(t, s.err, e)
	}
}

func assembleTagValues(values []string) ([]entities.EntitySearchQueryBuilderTag, error) {
	var tagValues []entities.EntitySearchQueryBuilderTag

	for _, x := range values {
		key, value, err := AssembleTagValue(x)

		if err != nil {
			return []entities.EntitySearchQueryBuilderTag{}, err
		}

		tagValues = append(tagValues, entities.EntitySearchQueryBuilderTag{Key: key, Value: value})
	}

	return tagValues, nil
}

func TestEntitiesAssembleTagValuesInput(t *testing.T) {
	var scenarios = []struct {
		tags     []string
		expected []entities.TaggingTagValueInput
		err      error
	}{
		{
			[]string{"one"},
			[]entities.TaggingTagValueInput{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"incomplete:"},
			[]entities.TaggingTagValueInput{},
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			[]string{"tag1:value1", "tag1:value2", "tag2:value1"},
			[]entities.TaggingTagValueInput{
				{Key: "tag1", Value: "value1"},
				{Key: "tag1", Value: "value2"},
				{Key: "tag2", Value: "value1"},
			},
			nil,
		},
	}

	for _, s := range scenarios {
		r, e := AssembleTagValuesInput(s.tags)

		assert.ElementsMatch(t, s.expected, r)
		assert.Equal(t, s.err, e)
	}
}
func TestEntitiesAssembleTagValue(t *testing.T) {
	var scenarios = []struct {
		tag           string
		expectedKey   string
		expectedValue string
		err           error
	}{
		{
			"invalidTag",
			"",
			"",
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			"incompleteTag:",
			"",
			"",
			fmt.Errorf("tag values must be specified as colon separated key:value pairs"),
		},
		{
			"validKey:validValue",
			"validKey",
			"validValue",
			nil,
		},
	}

	for _, s := range scenarios {
		k, v, e := AssembleTagValue(s.tag)
		assert.Equal(t, s.expectedKey, k)
		assert.Equal(t, s.expectedValue, v)
		assert.Equal(t, s.err, e)
	}
}
