package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/utils"
	utilsvalidation "github.com/newrelic/newrelic-cli/internal/utils/validation"
)

type clientFunc func(ctx context.Context, url string) ([]byte, error)

// AgentValidator attempts to validate that the infra agent
// was successfully installed and is sending data to New Relic.
type AgentValidator struct {
	fn                   clientFunc
	MaxAttempts          int
	IntervalMilliSeconds int
	Count                int
	ProgressIndicator    ux.ProgressIndicator
}

// AgentSuccessResponse represents the response object
// returned from infra agent `/v1/status/entity` endpoint.
//
// Docs: https://github.com/newrelic/infrastructure-agent/blob/master/docs/status_api.md#report-entity
type AgentSuccessResponse struct {
	GUID string `json:"guid"`
}

// NewAgentValidator returns a new instance of AgentValidator.
func NewAgentValidator() *AgentValidator {
	v := AgentValidator{
		MaxAttempts:          utilsvalidation.DefaultMaxAttempts,
		IntervalMilliSeconds: utilsvalidation.DefaultIntervalSeconds * 1000,
		ProgressIndicator:    ux.NewSpinner(),
		fn:                   getDefaultClientFunc(),
	}

	return &v
}

// Validate attempts to validate if the infra agent installation is successful.
// If it is successful, Validate returns the installed entity's GUID.
func (v *AgentValidator) Validate(ctx context.Context, url string) (string, error) {
	ticker := time.NewTicker(time.Duration(v.IntervalMilliSeconds) * time.Millisecond)
	defer ticker.Stop()

	v.ProgressIndicator.Start(utilsvalidation.ValidationInProgressMsg)
	defer v.ProgressIndicator.Stop()

	for {
		entityGUID, err := v.tryValidate(ctx, url)
		if err != nil {
			v.ProgressIndicator.Fail("")
			if strings.Contains(err.Error(), "context canceled") {
				return "", err
			}
			return "", fmt.Errorf("%s: %s", utilsvalidation.ReachedMaxValidationMsg, err)
		}

		if entityGUID != "" {
			v.ProgressIndicator.Success("")
			return entityGUID, nil
		}

		// This is no longer needed, it is implemented in the ExecWithRetries. Remove with the for loop
		select {
		case <-ticker.C:
			continue

		case <-ctx.Done():
			v.ProgressIndicator.Fail("")
			return "", fmt.Errorf("validation cancelled")
		}
	}
}

func (v *AgentValidator) tryValidate(ctx context.Context, url string) (string, error) {
	var guid string
	var err error

	fn := func() error {
		guid, err = v.executeAgentValidationRequest(ctx, url)
		return err
	}

	retry := utils.NewRetry(v.MaxAttempts, v.IntervalMilliSeconds, fn)
	if err := retry.ExecWithRetries(ctx); err != nil {
		return "", err
	}

	return guid, nil
}

func (v *AgentValidator) executeAgentValidationRequest(ctx context.Context, url string) (string, error) {
	data, err := v.fn(ctx, url)
	if err != nil {
		return "", err
	}

	response := AgentSuccessResponse{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return "", err
	}

	if response.GUID == "" {
		return "", errors.New("no entity GUID returned in response")
	}

	return response.GUID, nil
}

func getDefaultClientFunc() clientFunc {
	return func(ctx context.Context, url string) ([]byte, error) {
		return utils.NewHTTPClient().Get(ctx, url)
	}
}
