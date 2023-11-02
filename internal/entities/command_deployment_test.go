package entities

import (
	"errors"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/testcobra"

	"github.com/stretchr/testify/assert"
)

func TestEntityDeployment(t *testing.T) {
	command := cmdEntityDeployment

	assert.Equal(t, "deployment", command.Name())
}

func TestEntityDeploymentCreate(t *testing.T) {
	assert.Equal(t, "create", cmdEntityDeploymentCreate.Name())

	testcobra.CheckCobraMetadata(t, cmdEntityDeploymentCreate)
	testcobra.CheckCobraRequiredFlags(t, cmdEntityDeploymentCreate,
		[]string{"guid", "version"})
}

func TestParseAttributesSingleKeyValue(t *testing.T) {
	a := []string{
		"key:value",
	}

	var want = map[string]interface{}{
		"key": "value",
	}
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, *got)
}

func TestParseAttributesSingleIntegerKeyValue(t *testing.T) {
	a := []string{
		"a:1",
	}

	var want = map[string]interface{}{
		"a": "1",
	}
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, *got)
}

func TestParseAttributesSingleFloatingKeyValue(t *testing.T) {
	a := []string{
		"a:1.5",
	}

	var want = map[string]interface{}{
		"a": "1.5",
	}
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, *got)
}

func TestParseAttributesSingleBooleanKeyValue(t *testing.T) {
	a := []string{
		"a:true",
	}

	var want = map[string]interface{}{
		"a": "true",
	}
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, *got)
}

func TestParseAttributesMultipleTypesKeyValues(t *testing.T) {
	a := []string{
		"a:true",
		"b:1",
		"c:1.5",
		`d:"value"`,
	}

	var want = map[string]interface{}{
		"a": "true",
		"b": "1",
		"c": "1.5",
		"d": `"value"`,
	}
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, *got)
}

func TestParseAttributesTwoKeyValues(t *testing.T) {
	a := []string{
		"key:value",
		"key2:value2",
	}

	var want = map[string]interface{}{
		"key":  "value",
		"key2": "value2",
	}
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, *got)
}

func TestParseAttributesKeyNoValue(t *testing.T) {
	a := []string{
		"key",
	}

	want := nilPointerMapStringString()
	errWant := errors.New("invalid format, please use comma separated key-value pairs (--customAttribute key1:value1,key2:value2)")

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, got)
}

func TestParseAttributesTooManyColons(t *testing.T) {
	a := []string{
		"key:value:extra",
	}

	want := nilPointerMapStringString()
	errWant := errors.New("invalid format, please use comma separated key-value pairs (--customAttribute key1:value1,key2:value2)")

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, got)
}

func TestParseAttributesEmptyStringSlice(t *testing.T) {
	a := []string{}

	want := nilPointerMapStringString()
	var errWant error

	got, errGot := parseCustomAttributes(&a)

	assert.Equal(t, errWant, errGot)
	assert.Equal(t, want, got)
}

func nilPointerMapStringString() *map[string]interface{} {
	return nil
}
