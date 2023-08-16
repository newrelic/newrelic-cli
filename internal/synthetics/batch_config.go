package synthetics

// ----------------------------------------------------------------------------------
// After the API is out and is drawn into newrelic-client-go, all of the below
// datatypes will exist there (as fetched by Tutone), which functions in the CLI
// would call.
// ----------------------------------------------------------------------------------

// Following are the datatypes needed to unmarshal the
// the input YAML file so as to send a request to the syntheticsStartAutomatedTest mutation

type SecureCredential struct {
	Key         string `yaml:"key"`
	OverrideKey string `yaml:"overrideKey"`
}

type MonitorConfig struct {
	Overrides  []SecureCredential `yaml:"secureCredential,omitempty"`
	Location   string             `yaml:"location,omitempty"`
	IsBlocking bool               `yaml:"isBlocking"`
}

type Monitor struct {
	GUID   string        `yaml:"guid"`
	Config MonitorConfig `yaml:"config,omitempty"`
}

type Configuration struct {
	Monitors []Monitor `yaml:"monitors"`
	Config   struct {
		Branch    string `yaml:"branch"`
		Commit    string `yaml:"commit"`
		Platform  string `yaml:"platform"`
		DeepLink  string `yaml:"deepLink"`
		BatchName string `yaml:"batchName"`
	} `yaml:"config"`
}

// Following are the datatypes needed to unmarshal the
// mock JSON response of the query automatedTestResults

type Test struct {
	ID          string `json:"id"`
	BatchID     string `json:"batchId"`
	MonitorID   string `json:"monitorId"`
	MonitorGUID string `json:"monitorGuid"`
	MonitorName string `json:"monitorName"`
	Result      string `json:"result"`
}

type Root struct {
	Tests  []Test                 `json:"tests"`
	Config map[string]interface{} `json:"config"`
	Status string                 `json:"status"`
}

// Following are datatypes which are good to have to substitute
// the use of numbers as exit codes

type AutomatedTestResultsExitStatus int

const (
	AutomatedTestResultsExitStatusSuccess  AutomatedTestResultsExitStatus = 0
	AutomatedTestResultsExitStatusFailure  AutomatedTestResultsExitStatus = 1
	AutomatedTestResultsExitStatusTimedOut AutomatedTestResultsExitStatus = 3
	AutomatedTestResultsExitStatusUnknown  AutomatedTestResultsExitStatus = 2
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
