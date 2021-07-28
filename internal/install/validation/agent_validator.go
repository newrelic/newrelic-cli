package validation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/utils"
	utilsvalidation "github.com/newrelic/newrelic-cli/internal/utils/validation"
)

// AgentValidator attempts to validate that the infra agent
// was successfully installed and is sending data to New Relic.
type AgentValidator struct {
	clientFunc           func(context.Context, string) ([]byte, error)
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
func NewAgentValidator(clientFunc func(context.Context, string) ([]byte, error)) *AgentValidator {
	v := AgentValidator{
		MaxAttempts:          utilsvalidation.DefaultMaxAttempts,
		IntervalMilliSeconds: utilsvalidation.DefaultIntervalSeconds * 1000,
		ProgressIndicator:    ux.NewSpinner(),
		clientFunc:           clientFunc,
	}

	return &v
}

// Validate performs the attempt(s) to validate successful
// installation if the infra agent. If successful, it returns
// the installed entity's GUID.
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
			return "", fmt.Errorf("%s: %s", utilsvalidation.ReachexMaxValidationMsg, err)
		}

		if entityGUID != "" {
			v.ProgressIndicator.Success("")
			return entityGUID, nil
		}

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
	clientFunc := func() ([]byte, error) {
		return v.clientFunc(ctx, url)
	}

	validator := NewAgentValidatorFunc(clientFunc)

	retry := utils.NewRetry(v.MaxAttempts, v.IntervalMilliSeconds, validator.Func)
	if err := retry.ExecWithRetries(ctx); err != nil {
		v.Count = validator.Count
		return "", err
	}

	v.Count = validator.Count
	return validator.GUID, nil
}
