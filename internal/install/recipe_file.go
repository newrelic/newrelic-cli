package install

type recipeFile struct {
	Description    string                 `yaml:"description"`
	InputVars      []variableConfig       `yaml:"inputVars"`
	Install        map[string]interface{} `yaml:"install"`
	InstallTargets recipeInstallTarget    `yaml:"installTargets"`
	Keywords       []string               `yaml:"keywords"`
	Name           string                 `yaml:"name"`
	ProcessMatch   []string               `yaml:"processMatch"`
	Repository     string                 `yaml:"repository"`
	Variant        variant                `yaml:"variant"`
	ValidationNRQL string                 `yaml:"validationNrql"`
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

type variant struct {
	Arch              []string `yaml:"arch"`
	OS                []string `yaml:"os"`
	TargetEnvironment []string `yaml:"targetEnvironment"`
}

type recipeVariant struct {
	OS                []string `json:"os"`
	Arch              []string `json:"arch"`
	TargetEnvironment []string `json:"targetEnvironment"`
}
