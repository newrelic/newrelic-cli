// +build unit

package install

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

var (
	emptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 0.0,
		},
	}
	nonEmptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 1.0,
		},
	}
)

func TestValidate(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := newMockNrdbClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 1)

	v := newPollingRecipeValidator(c)

	r := recipe{}
	m := discoveryManifest{}

	ok, err := v.validate(context.Background(), m, r)

	require.NoError(t, err)
	require.True(t, ok)
}

func TestValidate_PassAfterNAttempts(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := newMockNrdbClient()
	v := newPollingRecipeValidator(c)
	v.maxAttempts = 5
	v.interval = 10 * time.Millisecond

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 5)

	r := recipe{}
	m := discoveryManifest{}

	ok, err := v.validate(context.Background(), m, r)

	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, 5, c.Attempts())
}

func TestValidate_FailAfterNAttempts(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := newMockNrdbClient()
	v := newPollingRecipeValidator(c)
	v.maxAttempts = 3
	v.interval = 10 * time.Millisecond

	r := recipe{}
	m := discoveryManifest{}

	ok, err := v.validate(context.Background(), m, r)

	require.NoError(t, err)
	require.False(t, ok)
	require.Equal(t, 3, c.Attempts())
}

func TestValidate_FailAfterMaxAttempts(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := newMockNrdbClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)

	v := newPollingRecipeValidator(c)
	v.maxAttempts = 1
	v.interval = 10 * time.Millisecond

	r := recipe{}
	m := discoveryManifest{}

	ok, err := v.validate(context.Background(), m, r)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestValidate_FailIfContextDone(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := newMockNrdbClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)

	v := newPollingRecipeValidator(c)
	v.interval = 1 * time.Second

	r := recipe{}
	m := discoveryManifest{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ok, err := v.validate(ctx, m, r)

	require.NoError(t, err)
	require.False(t, ok)
}

func TestValidate_QueryError(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := newMockNrdbClient()

	c.ThrowError("test error")

	v := newPollingRecipeValidator(c)

	r := recipe{}
	m := discoveryManifest{}

	ok, err := v.validate(context.Background(), m, r)

	require.False(t, ok)
	require.EqualError(t, err, "test error")
}
