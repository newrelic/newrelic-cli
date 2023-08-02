package synthetics

type SecureCredential struct {
	Key         string `yaml:"key"`
	OverrideKey string `yaml:"overrideKey"`
}

type MonitorConfig struct {
	Overrides  []SecureCredential `yaml:"secureCredential"`
	Location   string             `yaml:"location"`
	IsBlocking bool               `yaml:"isBlocking"`
}

type Monitor struct {
	GUID   string        `yaml:"guid"`
	Config MonitorConfig `yaml:"config"`
}

type TagSearch struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Configuration struct {
	Monitors  []Monitor   `yaml:"monitors"`
	TagSearch []TagSearch `yaml:"tagSearch"`
	Config    struct {
		Branch    string `yaml:"branch"`
		Commit    string `yaml:"commit"`
		Platform  string `yaml:"platform"`
		DeepLink  string `yaml:"deepLink"`
		BatchName string `yaml:"batchName"`
	} `yaml:"config"`
}

type Test struct {
	Id          string `json:"id"`
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
