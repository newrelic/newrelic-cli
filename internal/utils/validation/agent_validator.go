package validation

import (
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// AgentValidator polls NRDB to assert data is being reported for the given query.
type AgentValidator struct {
	httpClient utils.HTTPClient
}

// NewAgentValidator returns a new instance of AgentValidator.
func NewAgentValidator(c utils.HTTPClient) *AgentValidator {
	v := AgentValidator{
		httpClient: c,
	}

	return &v
}

// Validate
func (v *AgentValidator) Validate() (string, error) {

	_, err := v.httpClient.Get("https://af062943-dc76-45d1-8067-7849cbfe0d98.mock.pstmn.io/v1/status")
	if err != nil {
		return "", err
	}

	return ""
}
