package install

import (
	"gopkg.in/yaml.v2"
)

type recipe struct {
	ID             string     `json:"id"`
	File           string     `json:"file"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Repository     string     `json:"repository"`
	Keywords       []string   `json:"keywords"`
	ProcessMatch   []string   `json:"processMatch"`
	LogMatch       []logMatch `json:"logMatch"`
	ValidationNRQL string     `json:"validationNrql"`
	Vars           map[string]interface{}
}

// AddVar is responsible for including a new variable on the recipe Vars
// struct, which is used by go-task executor.
func (r *recipe) AddVar(key string, value interface{}) {
	if len(r.Vars) == 0 {
		r.Vars = make(map[string]interface{})
	}

	r.Vars[key] = value
}

func (r *recipe) ToRecipeFile() (*recipeFile, error) {
	var f recipeFile
	err := yaml.Unmarshal([]byte(r.File), &f)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
