//go:build unit
// +build unit

package utils

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	errorAfterAllRetry = "all retry attempts have been made"
)

func TestShouldRetryAndPass(t *testing.T) {
	m := MockFunc{
		CallsBeforeSuccess: 3,
	}
	r := NewRetry(3, 0, m.testErrorUntilFunc)
	ctx := r.ExecWithRetries(context.Background())
	require.Equal(t, 3, m.CallCount)
	require.Error(t, ctx.MostRecentError())
	require.True(t, ctx.Success)
	require.Equal(t, 3, ctx.RetryCount)
	require.Equal(t, 2, len(ctx.Errors))
}

func TestShouldRetryAndFail(t *testing.T) {
	m := MockFunc{}
	r := NewRetry(3, 0, m.testErrorFunc)
	ctx := r.ExecWithRetries(context.Background())
	require.Equal(t, 3, m.CallCount)
	require.False(t, ctx.Success)
	require.Equal(t, ctx.MostRecentError().Error(), errorAfterAllRetry)
	require.Equal(t, 3, ctx.RetryCount)
	require.Equal(t, 3, len(ctx.Errors))
}

func TestShouldNotRetry(t *testing.T) {
	m := MockFunc{
		CallsBeforeSuccess: 3,
	}
	r := NewRetry(3, 0, m.testOkFunc)
	ctx := r.ExecWithRetries(context.Background())
	require.Equal(t, 1, m.CallCount)
	require.NoError(t, ctx.MostRecentError())
	require.True(t, ctx.Success)
	require.Equal(t, 1, ctx.RetryCount)
	require.Equal(t, 0, len(ctx.Errors))
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
