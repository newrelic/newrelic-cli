// +build unit

package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	errorAfterAllRetry string = "all retry attempts have been made"
)

func TestShouldRetryAndPass(t *testing.T) {
	m := MockFunc{
		CallsBeforeSuccess: 3,
	}
	r := NewRetry(3, 0, m.testErrorUntilFunc)
	err := r.ExecWithRetries(context.Background())
	require.Equal(t, 3, m.CallCount)
	require.NoError(t, err)
}

func TestShouldRetryAndFail(t *testing.T) {
	m := MockFunc{}
	r := NewRetry(3, 0, m.testErrorFunc)
	err := r.ExecWithRetries(context.Background())
	require.Equal(t, 3, m.CallCount)
	require.Equal(t, err.Error(), errorAfterAllRetry)
}

func TestShouldNotRetry(t *testing.T) {
	m := MockFunc{
		CallsBeforeSuccess: 3,
	}
	r := NewRetry(3, 0, m.testOkFunc)
	err := r.ExecWithRetries(context.Background())
	require.Equal(t, 1, m.CallCount)
	require.NoError(t, err)
}

type MockFunc struct {
	CallCount          int
	CallsBeforeSuccess int
}

func (m *MockFunc) testErrorUntilFunc() error {
	m.CallCount++

	if m.CallCount < m.CallsBeforeSuccess {
		return errors.New(errorAfterAllRetry)
	}

	return nil
}

func (m *MockFunc) testErrorFunc() error {
	m.CallCount++

	return errors.New(errorAfterAllRetry)
}

func (m *MockFunc) testOkFunc() error {
	m.CallCount++

	return nil
}
