//go:build unit
// +build unit

package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoTaskGeneralError(t *testing.T) {
	err := errors.New(`some error`)
	e := NewGoTaskGeneralError(err)
	require.Equal(t, []string{}, e.TaskPath())
	require.Equal(t, e.Error(), "some error")

	err = errors.New(`task: Failed to run task "default": some error`)
	e = NewGoTaskGeneralError(err)
	require.Equal(t, []string{"default"}, e.TaskPath())
	require.Equal(t, e.Error(), "some error")

	err = errors.New(`task: Failed to run task "default": task: Failed to run task "subTask": some error`)
	e = NewGoTaskGeneralError(err)
	require.Equal(t, []string{"default", "subTask"}, e.TaskPath())
	require.Equal(t, e.Error(), "some error")

	err = errors.New(`task: Failed to run task "default": task: Failed to run task "subTask": task: Failed to run task "nestedSubTask": some error`)
	e = NewGoTaskGeneralError(err)
	require.Equal(t, []string{"default", "subTask", "nestedSubTask"}, e.TaskPath())
	require.Equal(t, e.Error(), "some error")
}

func TestNewCustomStdError(t *testing.T) {
	expected := map[string]interface{}{
		"someKey": "value",
	}

	stderr := errors.New("exit status 1")

	input := "{\"metadata\":{\"someKey\":\"value\"}}"
	incomingMessage := NewCustomStdError(stderr, input)
	require.NotNil(t, incomingMessage)
	require.Equal(t, expected, incomingMessage.Metadata)
}

func TestNewCustomStdError_ShouldReturnNilIfInvalidJSON(t *testing.T) {
	stderr := errors.New("exit status 1")

	input := "not a JSON string"
	incomingMessage := NewCustomStdError(stderr, input)
	require.Nil(t, incomingMessage)

	input = ""
	incomingMessage = NewCustomStdError(stderr, input)
	require.Nil(t, incomingMessage)

	input = "{\"metadata\":\"not a map\"}"
	incomingMessage = NewCustomStdError(stderr, input)
	require.Nil(t, incomingMessage)
}

func TestNewCustomStdError_ShouldSupportNestedJSONObjects(t *testing.T) {
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": "c",
		},
	}

	stderr := errors.New("exit status 1")

	input := "{\"metadata\":{\"a\":{\"b\":\"c\"}}}"
	incomingMessage := NewCustomStdError(stderr, input)
	require.NotNil(t, incomingMessage)
	require.Equal(t, expected, incomingMessage.Metadata)
}
