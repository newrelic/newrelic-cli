package types

type Recipe struct {
	ID             string     `json:"id"`
	File           string     `json:"file"`
	Name           string     `json:"name"`
	DisplayName    string     `json:"displayName"`
	Description    string     `json:"description"`
	Repository     string     `json:"repository"`
	Keywords       []string   `json:"keywords"`
	ProcessMatch   []string   `json:"processMatch"`
	LogMatch       []LogMatch `json:"logMatch"`
	ValidationNRQL string     `json:"validationNrql"`
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
