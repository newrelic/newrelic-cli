package execution

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutputParserShouldBuild(t *testing.T) {
	result := NewOutputParser(nil)
	assert.NotNil(t, result)
}

func TestOutputParserShouldEntityGuid(t *testing.T) {
	output := givenJSON("{\"EntityGuid\":\"abcd\"}")

	result := NewOutputParser(output)
	assert.Equal(t, "abcd", result.EntityGUID())
}

func TestOutputParserShouldEntityGuidCaseSensitive(t *testing.T) {
	output := givenJSON("{\"entityguid\":\"abcd\"}")

	result := NewOutputParser(output)
	assert.Equal(t, "", result.EntityGUID())
}

func TestOutputParserShouldNotEntityGuid(t *testing.T) {
	output := givenJSON("{\"SomethingElse\":\"abcd\"}")

	result := NewOutputParser(output)
	assert.NotEqual(t, "abcd", result.EntityGUID())
}

func TestOutputParserShouldGetMetadata(t *testing.T) {
	output := givenJSON("{\"Metadata\":{\"key1\":\"abcd\",\"key2\":\"efgh\"}}")

	result := NewOutputParser(output)
	assert.NotNil(t, result.Metadata())
	assert.Equal(t, result.Metadata()["key1"], "abcd")
	assert.Equal(t, result.Metadata()["key2"], "efgh")
}

func TestOutputParserShouldGetNoMetadata(t *testing.T) {
	output := givenJSON("{\"EntityGuid\":\"abcd\"}")
	result := NewOutputParser(output)
	assert.Equal(t, len(result.Metadata()), 0)
	assert.Nil(t, result.Metadata())
}

func TestOutputParserShouldGetMetadataMissing(t *testing.T) {
	output := givenJSON("{\"Metadata\":{}}")
	result := NewOutputParser(output)
	assert.NotNil(t, result.Metadata())
	assert.Equal(t, len(result.Metadata()), 0)
	assert.Equal(t, result.Metadata()["key1"], "")
}

func TestOutputParserShouldGetCapturedCliOutputFlag(t *testing.T) {
	// value is true
	output := givenJSON("{\"Metadata\":{\"CapturedCliOutput\":\"true\"}}")
	result := NewOutputParser(output)
	assert.NotNil(t, result.IsCapturedCliOutput())
	assert.True(t, result.IsCapturedCliOutput())

	// value is false
	output = givenJSON("{\"Metadata\":{\"CapturedCliOutput\":\"false\"}}")
	result = NewOutputParser(output)
	assert.NotNil(t, result.IsCapturedCliOutput())
	assert.False(t, result.IsCapturedCliOutput())

	// value is empty
	output = givenJSON("{\"Metadata\":{\"CapturedCliOutput\":\"\"}}")
	result = NewOutputParser(output)
	assert.False(t, result.IsCapturedCliOutput())

	// key does not exist
	output = givenJSON("{\"EntityGuid\":\"abcd\"}")
	result = NewOutputParser(output)
	assert.False(t, result.IsCapturedCliOutput())
}

func TestAddingMetadataAddsNewAndPreservesExistingValues(t *testing.T) {
	// start with one piece of metadata
	data := NewOutputParser(givenJSON("{\"Metadata\":{\"originalvalue\":\"abc\"}}"))
	assert.Equal(t, 1, len(data.Metadata()))

	// add a second
	data.AddMetadata("addedvalue", "123")
	assert.Equal(t, 2, len(data.Metadata()))

	// original should exist
	originalValue, originalExists := data.Metadata()["originalvalue"]
	assert.True(t, originalExists)
	assert.Equal(t, "abc", originalValue)

	// added value should exist
	addedValue, addedExists := data.Metadata()["addedvalue"]
	assert.True(t, addedExists)
	assert.Equal(t, "123", addedValue)

	// adding value to empty output doesn't break
	empty := NewOutputParser(nil)
	assert.Equal(t, 0, len(empty.Metadata()))

	empty.AddMetadata("newEntry", "456")
	assert.Equal(t, 1, len(empty.Metadata()))

	newEntry, newEntryExists := empty.Metadata()["newEntry"]
	assert.True(t, newEntryExists)
	assert.Equal(t, "456", newEntry)
}

func TestOutputParserShouldBeEmpty(t *testing.T) {
	output := givenJSON("")
	result := NewOutputParser(output)
	assert.Equal(t, "", result.EntityGUID())
}

func givenJSON(value string) map[string]interface{} {
	var result map[string]interface{}
	if value != "" {
		if err := json.Unmarshal([]byte(value), &result); err == nil {
			return result
		}
		log.Fatalf("couldn't unmarshal json for test with %s", value)
	}
	return map[string]interface{}{}
}
