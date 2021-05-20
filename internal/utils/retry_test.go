// +build unit

package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	m := MockFunc{
		CallsBeforeSuccess: 3,
	}
	r := NewRetry(3, 0, m.testFunc)
	err := r.ExecWithRetries()
	require.Equal(t, 3, m.CallCount)
	require.NoError(t, err)
}

type MockFunc struct {
	CallCount          int
	CallsBeforeSuccess int
}

func (m *MockFunc) testFunc() error {
	m.CallCount++

	if m.CallCount < m.CallsBeforeSuccess {
		return errors.New("")
	}

	return nil
}
