// +build unit

package validation

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
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
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 1)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.progressIndicator = pi

	r := types.Recipe{}
	m := types.DiscoveryManifest{}

	_, err := v.Validate(getTestContext(), m, r)

	require.NoError(t, err)
}

func TestValidate_PassAfterNAttempts(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := NewMockNRDBClient()
	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.progressIndicator = pi
	v.maxAttempts = 5
	v.interval = 10 * time.Millisecond

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 5)

	r := types.Recipe{}
	m := types.DiscoveryManifest{}

	_, err := v.Validate(getTestContext(), m, r)

	require.NoError(t, err)
	require.Equal(t, 5, c.Attempts())
}

func TestValidate_FailAfterNAttempts(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := NewMockNRDBClient()
	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.progressIndicator = pi
	v.maxAttempts = 3
	v.interval = 10 * time.Millisecond

	r := types.Recipe{}
	m := types.DiscoveryManifest{}

	_, err := v.Validate(getTestContext(), m, r)

	require.Error(t, err)
	require.Equal(t, 3, c.Attempts())
}

func TestValidate_FailAfterMaxAttempts(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.progressIndicator = pi
	v.maxAttempts = 1
	v.interval = 10 * time.Millisecond

	r := types.Recipe{}
	m := types.DiscoveryManifest{}

	_, err := v.Validate(getTestContext(), m, r)

	require.Error(t, err)
}

func TestValidate_FailIfContextDone(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.progressIndicator = pi
	v.interval = 1 * time.Second

	r := types.Recipe{}
	m := types.DiscoveryManifest{}

	ctx, cancel := context.WithCancel(getTestContext())
	cancel()

	_, err := v.Validate(ctx, m, r)

	require.Error(t, err)
}

func TestValidate_QueryError(t *testing.T) {
	credentials.SetDefaultProfile(credentials.Profile{AccountID: 12345})
	c := NewMockNRDBClient()

	c.ThrowError("test error")

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.progressIndicator = pi

	r := types.Recipe{}
	m := types.DiscoveryManifest{}

	_, err := v.Validate(getTestContext(), m, r)

	require.EqualError(t, err, "test error")
}

func getTestContext() context.Context {
	return context.WithValue(context.Background(), TestIdentifierKey, true)
}
