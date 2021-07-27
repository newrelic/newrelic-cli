package validation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/newrelic/newrelic-cli/internal/install/ux"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

const validationInProgressMsg string = "Checking for data in New Relic (this may take a few minutes)..."

// AgentValidator attempts to validate that the infra agent
// was successfully installed and is sending data to New Relic.
type AgentValidator struct {
	httpClient        utils.HTTPClientInterface
	MaxAttempts       int
	IntervalSeconds   int
	ProgressIndicator ux.ProgressIndicator
}

// AgentSuccessResponse represents the response object
// returned from infra agent `/v1/status/entity` endpoint.
//
// Docs: https://github.com/newrelic/infrastructure-agent/blob/master/docs/status_api.md#report-entity
type AgentSuccessResponse struct {
	GUID string `json:"guid"`
}

// NewAgentValidator returns a new instance of AgentValidator.
func NewAgentValidator(httpClient utils.HTTPClientInterface) *AgentValidator {
	v := AgentValidator{
		MaxAttempts:       60,
		IntervalSeconds:   5,
		ProgressIndicator: ux.NewSpinner(),
		httpClient:        httpClient,
	}

	return &v
}

// Validate performs the attempt(s) to validate successful
// installation if the infra agent. If successful, it returns
// the installed entity's GUID.
func (v *AgentValidator) Validate(ctx context.Context, url string) (string, error) {
	return v.waitForData(ctx, url)
}

func (v *AgentValidator) waitForData(ctx context.Context, url string) (string, error) {
	ticker := time.NewTicker(time.Duration(v.IntervalSeconds) * time.Second)
	defer ticker.Stop()

	v.ProgressIndicator.Start(validationInProgressMsg)
	defer v.ProgressIndicator.Stop()

	for {
		// log.WithFields(log.Fields{
		// 	"retryAttempts": count,
		// 	"url":           url,
		// }).Trace("validating installation")

		entityGUID, err := v.doValidate(ctx, url)
		if err != nil {
			v.ProgressIndicator.Fail("")
			return "", fmt.Errorf("reached max validation attempts: %s", err)
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

func (v *AgentValidator) doValidate(ctx context.Context, url string) (string, error) {
	var guid string

	retryFunc := func() error {
		var err error
		guid, err = v.executeAgentValidationRequest(ctx, url)
		if err != nil {
			return err
		}

		return nil
	}

	r := utils.NewRetry(v.MaxAttempts, v.IntervalSeconds, retryFunc)
	if err := r.ExecWithRetries(ctx); err != nil {
		return "", err
	}

	return guid, nil
}

func (v *AgentValidator) executeAgentValidationRequest(ctx context.Context, url string) (string, error) {
	data, err := v.httpClient.Get(ctx, url)
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
