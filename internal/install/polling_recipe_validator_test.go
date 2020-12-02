// +build unit

package install

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

func TestValidate(t *testing.T) {
	c := newMockNrdbClient()

	results := []nrdb.NrdbResult{
		map[string]interface{}{},
	}

	c.ReturnResultsAfterNAttempts(results, 1)

	v := newPollingRecipeValidator(c)

	r := recipe{}

	ok, err := v.validate(context.Background(), r)

	require.NoError(t, err)
	require.True(t, ok)
}

func TestValidate_PassAfterNAttempts(t *testing.T) {
	c := newMockNrdbClient()
	v := newPollingRecipeValidator(c)
	v.maxAttempts = 5
	v.interval = 10 * time.Millisecond

	results := []nrdb.NrdbResult{
		map[string]interface{}{},
	}

	c.ReturnResultsAfterNAttempts(results, 5)

	r := recipe{}

	ok, err := v.validate(context.Background(), r)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, 5, c.Attempts())
}

func TestValidate_FailAfterNAttempts(t *testing.T) {
	c := newMockNrdbClient()
	v := newPollingRecipeValidator(c)
	v.maxAttempts = 5
	v.interval = 10 * time.Millisecond

	results := []nrdb.NrdbResult{}

	c.ReturnResultsAfterNAttempts(results, 5)

	r := recipe{}

	ok, err := v.validate(context.Background(), r)

	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, 5, c.Attempts())
}

func TestValidate_FailAfterMaxAttempts(t *testing.T) {
	c := newMockNrdbClient()

	results := []nrdb.NrdbResult{
		map[string]interface{}{},
	}

	c.ReturnResultsAfterNAttempts(results, 2)

	v := newPollingRecipeValidator(c)
	v.maxAttempts = 1
	v.interval = 10 * time.Millisecond

	r := recipe{}

	ok, err := v.validate(context.Background(), r)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestValidate_FailIfContextDone(t *testing.T) {
	c := newMockNrdbClient()

	results := []nrdb.NrdbResult{
		map[string]interface{}{},
	}

	c.ReturnResultsAfterNAttempts(results, 2)

	v := newPollingRecipeValidator(c)
	v.interval = 1 * time.Second

	r := recipe{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ok, err := v.validate(ctx, r)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestValidate_QueryError(t *testing.T) {
	c := newMockNrdbClient()

	c.ThrowError("test error")

	v := newPollingRecipeValidator(c)

	r := recipe{}

	ok, err := v.validate(context.Background(), r)

	require.False(t, ok)
	require.EqualError(t, err, "test error")
}
