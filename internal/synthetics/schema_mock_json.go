package synthetics

// ----------------------------------------------------------------------------------
// After the API is out and is drawn into newrelic-client-go, all of the below
// datatypes will exist there (as fetched by Tutone), which functions in the CLI
// would call.
// ----------------------------------------------------------------------------------

// Following are the datatypes needed to unmarshal the
// mock JSON response of the query automatedTestResults

//type ResultantDomain struct {
//	Domain   string `json:"domain"`
//	Override string `json:"override"`
//}
//
//type ResultantSecureCredential struct {
//	Key         string `json:"key"`
//	OverrideKey string `json:"overrideKey"`
//}
//
//type ResultantOverrides struct {
//	Domain           ResultantDomain           `json:"domain"`
//	Location         string                    `json:"location"`
//	SecureCredential ResultantSecureCredential `json:"secureCredential"`
//	StartingUrl      string                    `json:"startingUrl"`
//}
//
//type AutomatedTestMonitorConfig struct {
//	IsBlocking bool               `json:"isBlocking"`
//	Overrides  ResultantOverrides `json:"overrides"`
//}
//
//type Test struct {
//	AutomatedTestMonitorConfig AutomatedTestMonitorConfig `json:"automatedTestMonitorConfig"`
//	BatchId                    string                     `json:"batchId"`
//	Duration                   int                        `json:"duration"`
//	Error                      string                     `json:"error"`
//	Id                         string                     `json:"id"`
//	Location                   string                     `json:"location"`
//	LocationLabel              string                     `json:"locationLabel"`
//	MonitorGuid                string                     `json:"monitorGuid"`
//	MonitorId                  string                     `json:"monitorId"`
//	MonitorName                string                     `json:"monitorName"`
//	Result                     string                     `json:"result"`
//	ResultsUrl                 string                     `json:"resultsUrl"`
//	Type                       string                     `json:"type"`
//	TypeLabel                  string                     `json:"typeLabel"`
//}
//
//type AutomatedTestResult struct {
//	Config struct {
//		BatchName  string `json:"batchName"`
//		Branch     string `json:"branch"`
//		Commit     string `json:"commit"`
//		DeepLink   string `json:"deepLink"`
//		Platform   string `json:"platform"`
//		Repository string `json:"repository"`
//	} `json:"config"`
//	Status string `json:"status"`
//	Tests  []Test `json:"tests"`
//}

// Following are datatypes which are good to have to substitute
// the use of numbers as exit codes

type AutomatedTestResultsExitStatus int

const (
	AutomatedTestResultsExitStatusSuccess  AutomatedTestResultsExitStatus = 0
	AutomatedTestResultsExitStatusFailure  AutomatedTestResultsExitStatus = 1
	AutomatedTestResultsExitStatusTimedOut AutomatedTestResultsExitStatus = 3
	AutomatedTestResultsExitStatusUnknown  AutomatedTestResultsExitStatus = 2
	AutomatedTestResultsExitStatusInProgress AutomatedTestResultsExitStatus = -1
)

type AutomatedTestResultsStatus string

const (
	AutomatedTestResultsStatusPassed     AutomatedTestResultsStatus = "PASSED"
	AutomatedTestResultsStatusFailure    AutomatedTestResultsStatus = "FAILED"
	AutomatedTestResultsStatusTimedOut   AutomatedTestResultsStatus = "TIMED_OUT"
	AutomatedTestResultsStatusInProgress AutomatedTestResultsStatus = "IN_PROGRESS"
)

var TestResultExitCodes = map[AutomatedTestResultsStatus]AutomatedTestResultsExitStatus{
	AutomatedTestResultsStatusPassed:     AutomatedTestResultsExitStatusSuccess,
	AutomatedTestResultsStatusFailure:    AutomatedTestResultsExitStatusFailure,
	AutomatedTestResultsStatusTimedOut:   AutomatedTestResultsExitStatusTimedOut,
	AutomatedTestResultsStatusInProgress: AutomatedTestResultsExitStatusUnknown,
}
