package install

type recipeFile struct {
	InputVars []variableConfig       `yaml:"inputVars"`
	Install   map[string]interface{} `yaml:"install"`
	MetaData  metaData               `yaml:"metadata"`
}

type variableConfig struct {
	Name    string `yaml:"name"`
	Prompt  string `yaml:"prompt"`
	Default string `yaml:"default"`
}

type metaData struct {
	Description  string    `yaml:"description"`
	Keywords     []string  `yaml:"keywords"`
	MELTMatch    meltMatch `yaml:"meltMatch"`
	Name         string    `yaml:"name"`
	ProcessMatch []string  `yaml:"processMatch"`
	Repository   string    `yaml:"repository"`
	Variant      variant   `yaml:"variant"`
}

type meltMatch struct {
	Events  patternMatcher `yaml:"events"`
	Metrics patternMatcher `yaml:"metrics"`
	Logging loggingMatcher `yaml:"logging"`
}

type variant struct {
	Arch              []string `yaml:"arch"`
	OS                []string `yaml:"os"`
	TargetEnvironment []string `yaml:"targetEnvironment"`
}

type patternMatcher struct {
	Pattern []string `yaml:"pattern"`
}

type loggingMatcher struct {
	patternMatcher
	Files []string `yaml:"files"`
}

type recipeVariant struct {
	OS                []string `json:"os"`
	Arch              []string `json:"arch"`
	TargetEnvironment []string `json:"targetEnvironment"`
}
