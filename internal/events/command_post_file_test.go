package events

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/testcobra"
)

func TestPostFile(t *testing.T) {
	assert.Equal(t, "postFile", cmdPostFile.Name())

	testcobra.CheckCobraMetadata(t, cmdPostFile)
	testcobra.CheckCobraRequiredFlags(t, cmdPostFile, []string{})
}

func TestShouldGetJSONArray(t *testing.T) {
	myJsonString := `[{"key1":"value1"},{"key2":"value2"}]`
	array := getArray([]byte(myJsonString))
	assert.Equal(t, 2, len(*array))
}

func TestShouldNotGetArrayWhenJSON(t *testing.T) {
	myJsonString := `{"some":"json"}`
	array := getArray([]byte(myJsonString))
	assert.Equal(t, 0, len(*array))
}

func TestShouldSliceBy(t *testing.T) {
	myJsonString := `[{"key1":"value1"},{"key2":"value2"},{"key3":"value3"},{"key4":"value5"},{"key5":"value5"}]`
	array := getArray([]byte(myJsonString))
	sliced := sliceBy(array, 3)
	assert.Equal(t, 3, len(sliced[0]))
	assert.Equal(t, 2, len(sliced[1]))
}

func TestShouldSliceByMax(t *testing.T) {
	myJsonString := `[{"key1":"value1"},{"key2":"value2"},{"key3":"value3"},{"key4":"value5"},{"key5":"value5"}]`
	array := getArray([]byte(myJsonString))
	sliced := sliceBy(array, 5)
	assert.Equal(t, 5, len(sliced[0]))
}
