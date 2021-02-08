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
	ProcessMatch   []string                              `json:"processMatch"`
	Repository     string                                `json:"repository"`
	ValidationNRQL string                                `json:"validationNrql"`
	Vars           map[string]interface{}
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
