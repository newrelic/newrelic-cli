//go:build unit
// +build unit

package types

import (
	"errors"
	"fmt"
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

func TestIncomingMessage_ShouldParseMetadata(t *testing.T) {
	mockIncomingMessage := IncomingMessage{
		Metadata: "{\"message\":\"original message\",\"metadata\":{\"someKey\":\"some value\"}}",
	}

	parsedMetadata := mockIncomingMessage.ParseMetadata()

	require.Equal(t, "some value", parsedMetadata["someKey"].(string))
}

func TestIncomingMessage_ShouldReturnRawMetadataWhenNonJSONString(t *testing.T) {
	mockIncomingMessage := IncomingMessage{
		Metadata: "This is a regular string",
	}

	parsedMetadata := mockIncomingMessage.ParseMetadata()

	fmt.Printf("\n\n parsedMetadata: %+v \n\n", parsedMetadata)

	require.Equal(t, "This is a regular string", parsedMetadata["metadata"])
}
