// +build unit

package validation

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/install/ux"
	utilsvalidation "github.com/newrelic/newrelic-cli/internal/utils/validation"
	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

const (
	TestIdentifierKey contextKey = iota
)

var (
	emptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count": 0.0,
		},
	}
	nonEmptyResults = []nrdb.NRDBResult{
		map[string]interface{}{
			"count":      1.0,
			"entityGuid": "an entity guid",
		},
	}
)

func TestValidate_shouldSucceed(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 1)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 1

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	_, err := v.ValidateRecipe(getTestContext(), m, r)

	require.NoError(t, err)
}

func TestValidate_shouldFailEmpty(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, emptyResults, 1)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 1

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	_, err := v.ValidateRecipe(getTestContext(), m, r)

	require.Error(t, err)
}

func TestValidate_PassAfterNAttempts(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()
	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 5
	v.IntervalMilliSeconds = 1

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 5)

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	_, err := v.ValidateRecipe(getTestContext(), m, r)

	require.NoError(t, err)
	require.Equal(t, 5, c.Attempts())
}

func TestValidate_FailAfterNAttempts(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()
	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 3
	v.IntervalMilliSeconds = 1

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	_, err := v.ValidateRecipe(getTestContext(), m, r)

	require.Error(t, err)
	require.Equal(t, 3, c.Attempts())
}

func TestValidate_FailAfterMaxAttempts(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 1
	v.IntervalMilliSeconds = 1

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	_, err := v.ValidateRecipe(getTestContext(), m, r)

	require.Error(t, err)
}

func TestValidate_FailIfContextDone(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()

	c.ReturnResultsAfterNAttempts(emptyResults, nonEmptyResults, 2)

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 2
	v.IntervalMilliSeconds = 1

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	ctx, cancel := context.WithCancel(getTestContext())
	cancel()

	_, err := v.ValidateRecipe(ctx, m, r)

	require.Error(t, err)
}

func TestValidate_QueryError(t *testing.T) {
	os.Setenv("NEW_RELIC_ACCOUNT_ID", "12345")
	c := NewMockNRDBClient()

	c.ThrowError("test error")

	pi := ux.NewMockProgressIndicator()
	v := NewPollingRecipeValidator(c)
	v.ProgressIndicator = pi
	v.MaxAttempts = 2
	v.IntervalMilliSeconds = 1

	r := types.OpenInstallationRecipe{}
	m := types.DiscoveryManifest{}

	_, err := v.ValidateRecipe(getTestContext(), m, r)

	require.EqualError(t, err, utilsvalidation.ReachexMaxValidationMsg+": test error")
}

func getTestContext() context.Context {
	return context.WithValue(context.Background(), TestIdentifierKey, true)
}
