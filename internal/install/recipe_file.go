package install

type recipeFile struct {
	Description    string                 `yaml:"description"`
	InputVars      []variableConfig       `yaml:"inputVars"`
	Install        map[string]interface{} `yaml:"install"`
	InstallTargets []recipeInstallTarget  `yaml:"installTargets"`
	Keywords       []string               `yaml:"keywords"`
	LogMatch       logMatch               `yaml:"logMatch"`
	Name           string                 `yaml:"name"`
	ProcessMatch   []string               `yaml:"processMatch"`
	Repository     string                 `yaml:"repository"`
	ValidationNRQL string                 `yaml:"validationNrql"`
}

type variableConfig struct {
	Name    string `yaml:"name"`
	Prompt  string `yaml:"prompt"`
	Secret  bool   `secret:"prompt"`
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

type logMatch struct {
	Name       string             `yaml:"name"`
	File       string             `yaml:"file"`
	Attributes logMatchAttributes `yaml:"attributes"`
	Pattern    string             `yaml:"pattern"`
	Systemd    string             `yaml:"systemd"`
}

type logMatchAttributes struct {
	LogType string `yaml:"logType"`
}
