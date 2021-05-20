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
