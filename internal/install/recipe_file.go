package install

type recipeFile struct {
	Description    string                 `yaml:"description"`
	InputVars      []variableConfig       `yaml:"inputVars"`
	Install        map[string]interface{} `yaml:"install"`
	InstallTargets recipeInstallTarget    `yaml:"installTargets"`
	Keywords       []string               `yaml:"keywords"`
	MELTMatch      meltMatch              `yaml:"meltMatch"`
	Name           string                 `yaml:"name"`
	ProcessMatch   []string               `yaml:"processMatch"`
	Repository     string                 `yaml:"repository"`
	Variant        variant                `yaml:"variant"`
}

type variableConfig struct {
	Name    string `yaml:"name"`
	Prompt  string `yaml:"prompt"`
	Default string `yaml:"default"`
}

type recipeInstallTarget struct {
	Type            string `yaml:"type"`
	OS              string `yaml:"os"`
	Platform        string `yaml:"platform"`
	PlatformFamily  string `yaml:"platformFamily"`
	PlatformVersion string `yaml:"platformVersion"`
	KernelVersion   string `yaml:"kernelVersion"`
	KernelArch      string `yaml:"kernelArch"`
}

type meltMatch struct {
	Events  patternMatcher `yaml:"events"`
	Metrics patternMatcher `yaml:"metrics"`
	Logging []logMatcher   `yaml:"logging"`
}

type variant struct {
	Arch              []string `yaml:"arch"`
	OS                []string `yaml:"os"`
	TargetEnvironment []string `yaml:"targetEnvironment"`
}

type patternMatcher struct {
	Pattern []string `yaml:"pattern"`
}

type logMatcher struct {
	// The name of the log file, as used by the Infra agent.
	Name string `yaml:"name"`

	// The path to the files, wildcards accepted.
	File string `yaml:"file"`

	// The pattern of the log messages themselves to make parsing and querying easier.
	Pattern string `yaml:"pattern"`
}

type recipeVariant struct {
	OS                []string `json:"os"`
	Arch              []string `json:"arch"`
	TargetEnvironment []string `json:"targetEnvironment"`
}
