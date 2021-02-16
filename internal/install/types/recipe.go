package types

type Recipe struct {
	ID             string                                `json:"id"`
	Description    string                                `json:"description"`
	DisplayName    string                                `json:"displayName"`
	File           string                                `json:"file"`
	InstallTargets []OpenInstallationRecipeInstallTarget `json:"installTargets"`
	Keywords       []string                              `json:"keywords"`
	LogMatch       []LogMatch                            `json:"logMatch"`
	Name           string                                `json:"name"`
	PreInstall     RecipePreInstall                      `json:"preInstall"`
	PostInstall    RecipePostInstall                     `json:"postInstall"`
	ProcessMatch   []string                              `json:"processMatch"`
	Repository     string                                `json:"repository"`
	ValidationNRQL string                                `json:"validationNrql"`
	Vars           map[string]interface{}
}

// RecipePreInstall represents the information used prior to recipe execution.
type RecipePreInstall struct {
	Info   string `yaml:"info"`
	Prompt string `yaml:"prompt"`
}

// RecipePostInstall represents the information used after recipe execution has completed.
type RecipePostInstall struct {
	Info   string `yaml:"info"`
	Prompt string `yaml:"prompt"`
}

func (r *Recipe) PostInstallMessage() string {
	if r.PostInstall.Info != "" {
		return r.PostInstall.Info
	}

	if r.PostInstall.Prompt != "" {
		return r.PostInstall.Prompt
	}

	return ""
}

func (r *Recipe) PreInstallMessage() string {
	if r.PreInstall.Info != "" {
		return r.PreInstall.Info
	}

	if r.PreInstall.Prompt != "" {
		return r.PreInstall.Prompt
	}

	return ""
}

// LogMatch represents a pattern that may match one or more logs on the underlying host.
type LogMatch struct {
	Name       string             `yaml:"name"`
	File       string             `yaml:"file"`
	Attributes LogMatchAttributes `yaml:"attributes,omitempty"`
	Pattern    string             `yaml:"pattern,omitempty"`
	Systemd    string             `yaml:"systemd,omitempty"`
}

// LogMatchAttributes contains metadata about its parent LogMatch.
type LogMatchAttributes struct {
	LogType string `yaml:"logtype"`
}

type RecipeVars map[string]string

// AddVar is responsible for including a new variable on the recipe Vars
// struct, which is used by go-task executor.
func (r *Recipe) AddVar(key string, value interface{}) {
	if len(r.Vars) == 0 {
		r.Vars = make(map[string]interface{})
	}

	r.Vars[key] = value
}
