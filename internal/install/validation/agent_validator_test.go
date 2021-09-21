//go:build unit
// +build unit

package validation

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

const infraAgentValidationURL = "infra validation url"

func TestAgentValidator_EntityGUID(t *testing.T) {
	t.Parallel()

	response := `{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`
	f := utils.CreateMockHTTPDoFunc(response, 200, nil)
	c := utils.NewMockHTTPClient(f)

	av := NewAgentValidator()
	av.fn = getMockClientFunc(c)

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw", result)
	require.Equal(t, err, nil)
	require.Equal(t, 1, c.GetCallCount)
}

func TestAgentValidator_AllAttemptFailed(t *testing.T) {
	t.Parallel()

	f := utils.CreateMockHTTPDoFunc("", 500, errors.New("an error was returned"))
	c := utils.NewMockHTTPClient(f)

	av := NewAgentValidator()
	av.fn = getMockClientFunc(c)
	av.MaxAttempts = 2
	av.IntervalMilliSeconds = 1

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "", result)
	require.Error(t, err)
	require.Equal(t, 2, c.GetCallCount)
}

func TestAgentValidator_EmptyGuidShouldFail(t *testing.T) {
	t.Parallel()

	f := utils.CreateMockHTTPDoFunc("", 200, nil)
	c := utils.NewMockHTTPClient(f)

	av := NewAgentValidator()
	av.fn = getMockClientFunc(c)
	av.MaxAttempts = 2
	av.IntervalMilliSeconds = 1

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "", result)
	require.Error(t, err)
	require.Equal(t, 2, c.GetCallCount)
}

func getMockClientFunc(c *utils.MockHTTPClient) clientFunc {
	return func(ctx context.Context, url string) ([]byte, error) { return c.Get(ctx, url) }
}
