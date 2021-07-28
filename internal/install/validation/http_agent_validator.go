package validation

import (
	"encoding/json"
	"errors"
)

type AgentValidatorFunc struct {
	Count     int
	RetryFunc func() error
	GUID      string
}

// NewAgentValidatorFunc returns a func to validate and get a GUID
func NewAgentValidatorFunc(clientFunc func() ([]byte, error)) *AgentValidatorFunc {

	var agentValidatorFunc = AgentValidatorFunc{
		Count: 0,
	}

	agentValidatorFunc.RetryFunc = func() error {
		agentValidatorFunc.Count++

		guid, err := executeAgentValidationRequest(clientFunc)
		if err != nil {
			return err
		}

		agentValidatorFunc.GUID = guid
		return nil
	}

	return &agentValidatorFunc
}

func executeAgentValidationRequest(clientFunc func() ([]byte, error)) (string, error) {
	data, err := clientFunc()
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
