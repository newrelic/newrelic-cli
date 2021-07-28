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
	response := `{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`
	c := utils.CreateMockGetResponse(response, nil)
	av := NewAgentValidator(c)

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw", result)
	require.Equal(t, err, nil)
	require.Equal(t, 1, av.Count)
}

func TestAgentValidator_AllAttemptFailed(t *testing.T) {
	t.Parallel()

	c := utils.CreateMockGetResponse("", errors.New("an error was returned"))
	av := NewAgentValidator(c)
	av.MaxAttempts = 2
	av.IntervalSeconds = 1

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "", result)
	require.Error(t, err)
	require.Equal(t, 2, av.Count)
}

func TestAgentValidator_EmptyGuidShouldFail(t *testing.T) {
	t.Parallel()

	c := utils.CreateMockGetResponse("", nil)
	av := NewAgentValidator(c)
	av.MaxAttempts = 2
	av.IntervalSeconds = 1

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "", result)
	require.Error(t, err)
	require.Equal(t, 2, av.Count)
}
