package synthetics

// ----------------------------------------------------------------------------------
// After the API is out and is drawn into newrelic-client-go, all of the below
// datatypes will exist there (as fetched by Tutone), which functions in the CLI
// would call.
// ----------------------------------------------------------------------------------

// Following are the datatypes needed to unmarshal the
// input YAML file to send a request to the syntheticsStartAutomatedTest mutation
// (updated schema based on changes made to the Product API)

type Config struct {
	BatchName  string `yaml:"batchName"`
	Branch     string `yaml:"branch"`
	Commit     string `yaml:"commit"`
	DeepLink   string `yaml:"deepLink"`
	Platform   string `yaml:"platform"`
	Repository string `yaml:"repository"`
}

type Domain struct {
	Domain   string `yaml:"domain"`
	Override string `yaml:"override"`
}

type SecureCredential struct {
	Key         string `yaml:"key"`
	OverrideKey string `yaml:"overrideKey"`
}

type Overrides struct {
	Domain           Domain           `yaml:"domain"`
	Location         string           `yaml:"location"`
	SecureCredential SecureCredential `yaml:"secureCredential"`
	StartingUrl      string           `yaml:"startingUrl"`
}

type TestConfig struct {
	IsBlocking bool      `yaml:"isBlocking"`
	Overrides  Overrides `yaml:"overrides"`
}

// SchemaTest (renamed to "SchemaTest" as "Test" already exists in the Mock JSON schema)
type SchemaTest struct {
	MonitorGuid string     `yaml:"monitorGuid"`
	Config      TestConfig `yaml:"config"`
}

type StartAutomatedTestInput struct {
	Config Config       `yaml:"config"`
	Tests  []SchemaTest `yaml:"tests"`
}
