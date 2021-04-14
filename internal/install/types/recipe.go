package types

import "strings"

type Recipe struct {
	ID                string                                   `json:"id"`
	Description       string                                   `json:"description"`
	DisplayName       string                                   `json:"displayName"`
	File              string                                   `json:"file"`
	InstallTargets    []OpenInstallationRecipeInstallTarget    `json:"installTargets"`
	Keywords          []string                                 `json:"keywords"`
	LogMatch          []LogMatch                               `json:"logMatch"`
	Name              string                                   `json:"name"`
	PreInstall        OpenInstallationPreInstallConfiguration  `json:"preInstall"`
	PostInstall       OpenInstallationPostInstallConfiguration `json:"postInstall"`
	ProcessMatch      []string                                 `json:"processMatch"`
	Repository        string                                   `json:"repository"`
	SuccessLinkConfig SuccessLinkConfig                        `json:"successLinkConfig"`
	ValidationNRQL    string                                   `json:"validationNrql"`
	Vars              map[string]interface{}
}

func (r *Recipe) PostInstallMessage() string {
	if r.PostInstall.Info != "" {
		return r.PostInstall.Info
	}

	return ""
}

func (r *Recipe) PreInstallMessage() string {
	if r.PreInstall.Info != "" {
		return r.PreInstall.Info
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

func (r *Recipe) IsApm() bool {
	return r.HasKeyword("apm")
}

func (r *Recipe) HasHostTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.HOST)
}

func (r *Recipe) HasApplicationTargetType() bool {
	return r.HasTargetType(OpenInstallationTargetTypeTypes.APPLICATION)
}

func (r *Recipe) HasKeyword(keyword string) bool {
	if len(r.Keywords) == 0 {
		return false
	}

	for _, single := range r.Keywords {
		return strings.EqualFold(single, keyword)
	}

	return false
}

func (r *Recipe) HasTargetType(t OpenInstallationTargetType) bool {
	if len(r.InstallTargets) == 0 {
		return false
	}

	for _, target := range r.InstallTargets {
		if target.Type == t {
			return true
		}
	}

	return false
}
