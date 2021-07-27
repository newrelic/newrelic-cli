// +build unit

package validation

import (
	"context"
	"errors"
	"testing"

	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/stretchr/testify/require"
)

const infraAgentValidationURL = "http://localhost:18003/v1/status/entity"

func TestAgentValidator_EntityGUID(t *testing.T) {
	response := `{"GUID":"MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw"}`
	c := utils.NewMockHTTPClient(utils.CreateMockHTTPDoFunc(response, 200, nil))
	av := NewAgentValidator(c)

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw", result)
	require.Equal(t, err, nil)
	require.Equal(t, 1, c.GetCallCount)
}

func TestAgentValidator_NoContentResponse(t *testing.T) {
	t.Parallel()

	c := utils.NewMockHTTPClient(utils.CreateMockHTTPDoFunc("", 204, nil))
	av := NewAgentValidator(c)

	ctx := context.Background()
	result, err := av.Validate(ctx, infraAgentValidationURL)

	require.Equal(t, "", result)
	require.Equal(t, err, nil)
	require.Equal(t, 1, c.GetCallCount)
}

func TestAgentValidator_InternalServerError(t *testing.T) {
	t.Parallel()

	c := utils.NewMockHTTPClient(utils.CreateMockHTTPDoFunc("", 500, errors.New("Internal Server Error")))
	av := NewAgentValidator(c)

	ctx := context.Background()
	_, err := av.Validate(ctx, infraAgentValidationURL)

	require.Error(t, err)
	require.Equal(t, 1, c.GetCallCount)
}
