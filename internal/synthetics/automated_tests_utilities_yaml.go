package synthetics

import "github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"

// a wrapper structure for the input to be sent to syntheticsStartAutomatedTest
type StartAutomatedTestInput struct {
	Config synthetics.SyntheticsAutomatedTestConfigInput    `json:"config,omitempty"`
	Tests  []synthetics.SyntheticsAutomatedTestMonitorInput `json:"tests,omitempty"`
}

var globalResultExitCodes = map[synthetics.SyntheticsAutomatedTestStatus]*int{
	synthetics.SyntheticsAutomatedTestStatusTypes.FAILED:  intPtr(1),
	synthetics.SyntheticsAutomatedTestStatusTypes.PASSED:  intPtr(0),
	synthetics.SyntheticsAutomatedTestStatusTypes.TIMEOUT: intPtr(3),
}

func intPtr(value int) *int {
	return &value
}
