package utils

import (
	"errors"
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func AssembleTagsInput(tags []string) ([]entities.TaggingTagInput, error) {
	var t []entities.TaggingTagInput

	tagBuilder := make(map[string][]string)

	for _, x := range tags {
		if !strings.Contains(x, ":") {
			return []entities.TaggingTagInput{}, errors.New("tags must be specified as colon separated key:value pairs")
		}

		v := strings.SplitN(x, ":", 2)

		tagBuilder[v[0]] = append(tagBuilder[v[0]], v[1])
	}

	for k, v := range tagBuilder {
		tag := entities.TaggingTagInput{
			Key:    k,
			Values: v,
		}

		t = append(t, tag)
	}

	return t, nil
}

func AssembleTagValuesInput(values []string) ([]entities.TaggingTagValueInput, error) {
	var tagValues []entities.TaggingTagValueInput

	for _, x := range values {
		key, value, err := AssembleTagValue(x)

		if err != nil {
			return []entities.TaggingTagValueInput{}, err
		}

		tagValues = append(tagValues, entities.TaggingTagValueInput{Key: key, Value: value})
	}

	return tagValues, nil
}

func AssembleTagValue(tagValueString string) (string, string, error) {
	tagFormatError := errors.New("tag values must be specified as colon separated key:value pairs")

	if !strings.Contains(tagValueString, ":") {
		return "", "", tagFormatError
	}

	v := strings.SplitN(tagValueString, ":", 2)

	// Handle incomplete tag where the value portion is empty
	if v[1] == "" {
		return "", "", tagFormatError
	}

	return v[0], v[1], nil
}
