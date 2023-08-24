package synthetics

import "github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"

// a wrapper structure for the input to be sent to syntheticsStartAutomatedTest
type SyntheticsStartAutomatedTestInput struct {
	Config synthetics.SyntheticsAutomatedTestConfigInput    `json:"config,omitempty"`
	Tests  []synthetics.SyntheticsAutomatedTestMonitorInput `json:"tests,omitempty"`
}


var TestResultExitCodes = map[synthetics.SyntheticsAutomatedTestStatus]*int{
	synthetics.SyntheticsAutomatedTestStatusTypes.FAILED:      intPtr(1),
	synthetics.SyntheticsAutomatedTestStatusTypes.PASSED:      intPtr(0),
	synthetics.SyntheticsAutomatedTestStatusTypes.TIMEOUT:     intPtr(3),

}

func intPtr(value int) *int {
	return &value
}

// Following are datatypes which are good to have to substitute
// the use of numbers as exit codes

// type AutomatedTestResultsExitStatus int

// const (
// 	AutomatedTestResultsExitStatusSuccess  AutomatedTestResultsExitStatus = 0
// 	AutomatedTestResultsExitStatusFailure  AutomatedTestResultsExitStatus = 1
// 	AutomatedTestResultsExitStatusTimedOut AutomatedTestResultsExitStatus = 3
// 	AutomatedTestResultsExitStatusUnknown  AutomatedTestResultsExitStatus = 2
// 	AutomatedTestResultsExitStatusInProgress AutomatedTestResultsExitStatus = -1
// )

// type AutomatedTestResultsStatus string

// const (
// 	AutomatedTestResultsStatusPassed     AutomatedTestResultsStatus = "PASSED"
// 	AutomatedTestResultsStatusFailure    AutomatedTestResultsStatus = "FAILED"
// 	AutomatedTestResultsStatusTimedOut   AutomatedTestResultsStatus = "TIMED_OUT"
// 	AutomatedTestResultsStatusInProgress AutomatedTestResultsStatus = "IN_PROGRESS"
// )

// var TestResultExitCodes = map[AutomatedTestResultsStatus]AutomatedTestResultsExitStatus{
// 	AutomatedTestResultsStatusPassed:     AutomatedTestResultsExitStatusSuccess,
// 	AutomatedTestResultsStatusFailure:    AutomatedTestResultsExitStatusFailure,
// 	AutomatedTestResultsStatusTimedOut:   AutomatedTestResultsExitStatusTimedOut,
// 	AutomatedTestResultsStatusInProgress: AutomatedTestResultsExitStatusUnknown,
// }
