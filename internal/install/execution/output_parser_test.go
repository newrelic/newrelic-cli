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
